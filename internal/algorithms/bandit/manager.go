package bandit

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/algorithms/wasm"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

const (
	defaultSyncStatsInterval = time.Second * 5
	defaultPollInterval      = time.Second * 5
)

type VariantStats struct {
	Evaluations uint64
	Successes   uint64
	Failures    uint64
	MetricSum   decimal.Decimal
}

// AlgorithmState holds stats for one algorithm (per feature).
type AlgorithmState struct {
	AlgorithmType domain.AlgorithmType
	Enabled       bool

	// multivariant state (for regular bandits)
	Variants    map[string]*VariantStats // variant_key -> stats
	VariantsArr []string
	Settings    map[string]decimal.Decimal

	// optimizer state
	IsOptimizer  bool
	Iteration    uint64
	CurrentValue decimal.Decimal
	MetricSum    decimal.Decimal
	BestValue    decimal.Decimal
	BestReward   decimal.Decimal
	LastError    decimal.Decimal
	Integral     decimal.Decimal
	StepSize     decimal.Decimal
	Temperature  decimal.Decimal

	// contextual bandit state
	IsContextual    bool
	ContextualState *ContextualAlgorithmState

	mu sync.RWMutex
}

type StateKey struct {
	FeatureKey string
	EnvKey     string
}

type BanditManager struct {
	mu                sync.RWMutex
	state             map[StateKey]*AlgorithmState
	randSrc           *rand.Rand
	syncStatsInterval time.Duration
	pollInterval      time.Duration
	lastSeen          time.Time
	stopCh            chan struct{}

	featureAlgorithmsRepo contract.FeatureAlgorithmsRepository
	featureVariantsRepo   contract.FlagVariantsRepository
	statsRepo             contract.FeatureAlgorithmStatsRepository
	optimizerStatsRepo    contract.FeatureOptimizerStatsRepository
	contextualStatsRepo   contract.FeatureContextualStatsRepository
	auditRepo             contract.AuditLogRepository
	featuresUseCase       contract.FeaturesUseCase
	envsUseCase           contract.EnvironmentsUseCase

	// Custom WASM algorithms support
	wasmManager *wasm.Manager
	customState map[StateKey]*customState
	customMu    sync.RWMutex

	nowFunc func() time.Time
}

func New(
	featureAlgorithmsRepo contract.FeatureAlgorithmsRepository,
	featureVariantsRepo contract.FlagVariantsRepository,
	statsRepo contract.FeatureAlgorithmStatsRepository,
	optimizerStatsRepo contract.FeatureOptimizerStatsRepository,
	contextualStatsRepo contract.FeatureContextualStatsRepository,
	auditRepo contract.AuditLogRepository,
	featuresUseCase contract.FeaturesUseCase,
	envsUseCase contract.EnvironmentsUseCase,
) (*BanditManager, error) {
	mngr := &BanditManager{
		state:                 make(map[StateKey]*AlgorithmState),
		customState:           make(map[StateKey]*customState),
		randSrc:               rand.New(rand.NewSource(time.Now().UnixNano())),
		syncStatsInterval:     defaultSyncStatsInterval,
		pollInterval:          defaultPollInterval,
		featureAlgorithmsRepo: featureAlgorithmsRepo,
		featureVariantsRepo:   featureVariantsRepo,
		statsRepo:             statsRepo,
		optimizerStatsRepo:    optimizerStatsRepo,
		contextualStatsRepo:   contextualStatsRepo,
		auditRepo:             auditRepo,
		featuresUseCase:       featuresUseCase,
		envsUseCase:           envsUseCase,
		stopCh:                make(chan struct{}),
		nowFunc:               time.Now,
	}

	if err := mngr.loadState(); err != nil {
		return nil, fmt.Errorf("load state: %w", err)
	}

	return mngr, nil
}

func (m *BanditManager) Start(context.Context) error {
	go m.syncToDBLoop() //nolint:contextcheck // false positive

	go func() { //nolint:contextcheck // false positive
		if err := m.Watch(context.Background()); err != nil {
			slog.Error("Bandits: failed to watch features updates", "error", err)
		}
	}()

	return nil
}

func (m *BanditManager) Stop(context.Context) error {
	close(m.stopCh)

	if err := m.flushAllToDB(); err != nil { //nolint:contextcheck // false positive
		slog.Error("bandit: failed flush", "error", err)
	}

	return nil
}

func (m *BanditManager) GetAlgorithmState(featureKey, envKey string) (*AlgorithmState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := StateKey{FeatureKey: featureKey, EnvKey: envKey}
	state, ok := m.state[key]

	return state, ok
}

func (m *BanditManager) HasAlgorithm(featureKey, envKey string) bool {
	state, ok := m.GetAlgorithmState(featureKey, envKey)
	if !ok {
		return false
	}

	return state.Enabled
}

func (m *BanditManager) GetAlgorithmKind(featureKey, envKey string) (domain.AlgorithmKind, bool) {
	state, ok := m.GetAlgorithmState(featureKey, envKey)
	if !ok {
		return domain.AlgorithmKindUnknown, false
	}

	if !state.Enabled {
		return domain.AlgorithmKindUnknown, false
	}

	return state.AlgorithmType.Kind(), true
}

// EvaluateFeature chooses a variant according to the bandit algorithm (multi-variant).
//
//nolint:exhaustive // false positive
func (m *BanditManager) EvaluateFeature(featureKey, envKey string) (string, bool) {
	state, ok := m.GetAlgorithmState(featureKey, envKey)
	if !ok {
		return "", false
	}

	if !state.Enabled {
		return "", false
	}

	if state.IsOptimizer {
		return "", false
	}

	var variant string

	// choose bandit algorithm
	switch state.AlgorithmType {
	case domain.AlgorithmTypeEpsilonGreedy:
		variant = m.evalEpsilonGreedy(state)
	case domain.AlgorithmTypeThompsonSampling:
		variant = m.evalThompson(state)
	case domain.AlgorithmTypeUCB:
		variant = m.evalUCB(state)
	case domain.AlgorithmTypeUnknown:
		return "", false
	default:
		return "", false
	}

	return variant, true
}

// EvaluateOptimizer returns optimized value for single-variant feature (optimizer algorithms).
func (m *BanditManager) EvaluateOptimizer(featureKey, envKey string) (decimal.Decimal, bool) {
	state, ok := m.GetAlgorithmState(featureKey, envKey)
	if !ok {
		return decimal.Zero, false
	}

	if !state.Enabled {
		return decimal.Zero, false
	}

	if !state.IsOptimizer {
		return decimal.Zero, false
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	// For optimizers, we return the current value
	// The value will be updated on next track event with reward
	return state.CurrentValue, true
}

// EvaluateContextual chooses a variant according to contextual bandit algorithm.
func (m *BanditManager) EvaluateContextual(featureKey, envKey string, ctx map[string]any) (string, bool) {
	state, ok := m.GetAlgorithmState(featureKey, envKey)
	if !ok {
		return "", false
	}

	if !state.Enabled {
		return "", false
	}

	if !state.IsContextual {
		return "", false
	}

	var variant string

	//nolint:exhaustive // only contextual algorithms
	switch state.AlgorithmType {
	case domain.AlgorithmTypeLinUCB:
		variant = m.evalLinUCB(state, ctx)
	case domain.AlgorithmTypeContextualThompson:
		variant = m.evalContextualThompson(state, ctx)
	case domain.AlgorithmTypeContextualEpsilon:
		variant = m.evalContextualEpsilon(state, ctx)
	default:
		return "", false
	}

	return variant, variant != ""
}

// HandleContextualTrackEvent updates contextual bandit model with observed reward.
func (m *BanditManager) HandleContextualTrackEvent(
	featureKey string,
	envKey string,
	variantKey string,
	eventType domain.FeedbackEventType,
	metric decimal.Decimal,
	ctx map[string]any,
) {
	state, ok := m.GetAlgorithmState(featureKey, envKey)
	if !ok {
		return
	}

	if !state.IsContextual {
		return
	}

	// Only process success/failure events
	if eventType != domain.FeedbackEventTypeSuccess && eventType != domain.FeedbackEventTypeFailure {
		return
	}

	reward := metric.InexactFloat64()
	if eventType == domain.FeedbackEventTypeFailure {
		reward = -reward
	}

	//nolint:exhaustive // only contextual algorithms
	switch state.AlgorithmType {
	case domain.AlgorithmTypeLinUCB:
		m.updateLinUCB(state, variantKey, reward, ctx)
	case domain.AlgorithmTypeContextualThompson:
		m.updateContextualThompson(state, variantKey, reward, ctx)
	case domain.AlgorithmTypeContextualEpsilon:
		m.updateContextualEpsilon(state, variantKey, reward, ctx)
	}
}

// HandleTrackEvent called by track consumer to update in-memory counters.
//
//nolint:exhaustive // false positive
func (m *BanditManager) HandleTrackEvent(
	featureKey string,
	envKey string,
	variantKey string,
	eventType domain.FeedbackEventType,
	metric decimal.Decimal,
) {
	state, ok := m.GetAlgorithmState(featureKey, envKey)
	if !ok {
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	// Handle optimizer algorithms
	if state.IsOptimizer {
		if eventType == domain.FeedbackEventTypeSuccess || eventType == domain.FeedbackEventTypeFailure {
			reward := metric
			if eventType == domain.FeedbackEventTypeFailure {
				reward = metric.Neg() // negative reward for failures
			}

			switch state.AlgorithmType {
			case domain.AlgorithmTypeHillClimb:
				m.evalHillClimb(state, reward)
			case domain.AlgorithmTypeSimAnnealing:
				m.evalSimulatedAnnealing(state, reward)
			case domain.AlgorithmTypePIDController:
				// PID needs target value from settings
				// TODO: get target from feature-flag settings
				target := getSettingAsDecimal(state.Settings, "target", 1.0)
				m.evalPID(state, reward, target)
			// BayesOpt and CEM need batch processing, skip for now
			case domain.AlgorithmTypeBayesOpt, domain.AlgorithmTypeCEM:
				// TODO: implement BayesOpt and CEM
				// These require batch evaluation, not suitable for online learning
				// Store metric for potential batch update later
				state.MetricSum = state.MetricSum.Add(reward)
			default:
			}
		}

		return
	}

	// Handle bandit algorithms (multi-variant)
	vs, ok := state.Variants[variantKey]
	if !ok {
		vs = &VariantStats{}
		state.Variants[variantKey] = vs
		state.VariantsArr = append(state.VariantsArr, variantKey)
	}
	switch eventType {
	case domain.FeedbackEventTypeEvaluation:
		vs.Evaluations++
	case domain.FeedbackEventTypeSuccess:
		vs.Successes++
		vs.MetricSum = vs.MetricSum.Add(metric)
	case domain.FeedbackEventTypeFailure:
		vs.Failures++
	case domain.FeedbackEventTypeError:
		vs.Failures++
	case domain.FeedbackEventTypeUnknown:
	default:
		// custom handling
	}
}

//nolint:gocognit // This is a complex function.
func (m *BanditManager) Watch(ctx context.Context) error {
	slog.Info("Bandit Manager: start watching features")

	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	last := m.lastSeen

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-m.stopCh:
			return nil
		case <-ticker.C:
			windowEnd := m.nowFunc().UTC()

			logs, err := m.auditRepo.ListSince(ctx, last)
			if err != nil {
				slog.Error("Watcher: audit log list since failed", "err", err)

				last = windowEnd

				continue
			}

			if len(logs) == 0 {
				continue
			}

			slog.Info("Bandits Manager: changes on features detected", "count", len(logs))

			// Deduplicate by project+feature; keep delete if any delete for a feature entity appears.
			type changeKey struct {
				featureID domain.FeatureID
				envID     domain.EnvironmentID
				envKey    string
			}

			changes := make(map[changeKey]domain.AuditAction)

			for _, row := range logs {
				if row.FeatureID == "" {
					continue
				}

				key := changeKey{
					featureID: row.FeatureID,
					envID:     row.EnvironmentID,
					envKey:    row.EnvKey,
				}

				if row.Entity == domain.EntityFeatureAlgorithm && row.Action == domain.AuditActionDelete {
					changes[key] = domain.AuditActionDelete

					continue
				}

				if row.Entity == domain.EntityFeatureAlgorithm || row.Entity == domain.EntityFlagVariant {
					if _, ok := changes[key]; !ok {
						changes[key] = domain.AuditActionUpdate
					}
				}
			}

			for key, action := range changes {
				if action == domain.AuditActionDelete {
					if err := m.removeFeatureFromState(
						ctx,
						key.featureID,
						key.envKey,
					); err != nil {
						slog.Error("Bandit Manager: delete feature from state failed",
							"featureID", key.featureID, "env", key.envKey, "err", err)
					}

					continue
				}

				// refresh feature by loading from repo
				if err := m.refreshFeature(
					ctx,
					key.featureID,
					key.envID,
					key.envKey,
				); err != nil {
					slog.Error("Bandit Manager: refresh feature failed",
						"feature", key.featureID, "error", err)
				}
			}

			last = windowEnd
		}
	}
}

//nolint:nestif,maintidx,gocognit
func (m *BanditManager) loadState() error {
	now := m.nowFunc()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	records, err := m.featureAlgorithmsRepo.ListAllExtended(ctx)
	if err != nil {
		return fmt.Errorf("list feature algorithms: %w", err)
	}

	banditStats, err := m.statsRepo.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("load bandit stats: %w", err)
	}

	optimizerStats, err := m.optimizerStatsRepo.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("load optimizer stats: %w", err)
	}

	contextualStats, err := m.contextualStatsRepo.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("load contextual stats: %w", err)
	}

	allVariants, err := m.featureVariantsRepo.ListExtended(ctx)
	if err != nil {
		return fmt.Errorf("load all variants: %w", err)
	}

	allVariantsMap := make(map[StateKey][]string)
	for _, variant := range allVariants {
		key := StateKey{
			FeatureKey: variant.FeatureKey,
			EnvKey:     variant.EnvKey,
		}
		value, ok := allVariantsMap[key]
		if !ok {
			value = []string{variant.Name}
		} else {
			value = append(value, variant.Name)
		}
		allVariantsMap[key] = value
	}

	type statsKey struct {
		FeatureID domain.FeatureID
		EnvID     domain.EnvironmentID
		AlgSlug   string
		Variant   string
	}

	banditStatsMap := make(map[statsKey]VariantStats)
	for _, stat := range banditStats {
		key := statsKey{
			FeatureID: stat.FeatureID,
			EnvID:     stat.EnvironmentID,
			AlgSlug:   stat.AlgorithmSlug,
			Variant:   stat.VariantKey,
		}
		banditStatsMap[key] = VariantStats{
			Evaluations: stat.Evaluations,
			Successes:   stat.Successes,
			Failures:    stat.Failures,
			MetricSum:   stat.MetricSum,
		}
	}

	type optimizerKey struct {
		FeatureID domain.FeatureID
		EnvID     domain.EnvironmentID
		AlgSlug   string
	}

	optimizerStatsMap := make(map[optimizerKey]domain.FeatureOptimizerStats)
	for _, stat := range optimizerStats {
		key := optimizerKey{
			FeatureID: stat.FeatureID,
			EnvID:     stat.EnvironmentID,
			AlgSlug:   stat.AlgorithmSlug,
		}
		optimizerStatsMap[key] = stat
	}

	contextualStatsMap := make(map[statsKey]domain.FeatureContextualStats)
	for _, stat := range contextualStats {
		key := statsKey{
			FeatureID: stat.FeatureID,
			EnvID:     stat.EnvironmentID,
			AlgSlug:   stat.AlgorithmSlug,
			Variant:   stat.VariantKey,
		}
		contextualStatsMap[key] = stat
	}

	state := make(map[StateKey]*AlgorithmState, len(records))
	for _, record := range records {
		// Skip custom algorithms (handled by WASM manager)
		if record.AlgorithmSlug == nil {
			continue
		}

		algSlug := *record.AlgorithmSlug

		key := StateKey{
			FeatureKey: record.FeatureKey,
			EnvKey:     record.EnvKey,
		}

		algType := domain.AlgorithmSlugToType(algSlug)
		algKind := algType.Kind()

		if algKind == domain.AlgorithmKindOptimizer {
			optStats, hasStats := optimizerStatsMap[optimizerKey{
				FeatureID: record.FeatureID,
				EnvID:     record.EnvironmentID,
				AlgSlug:   algSlug,
			}]

			algState := &AlgorithmState{
				AlgorithmType: algType,
				Enabled:       record.Enabled,
				Settings:      record.Settings,
				IsOptimizer:   true,
				Iteration:     0,
				CurrentValue:  getSettingAsDecimal(record.Settings, "initial_value", 0.0),
				MetricSum:     decimal.Zero,
				Variants:      make(map[string]*VariantStats),
				VariantsArr:   []string{},
			}

			if hasStats {
				algState.Iteration = optStats.Iteration
				algState.CurrentValue = optStats.CurrentValue
				algState.MetricSum = optStats.MetricSum
				algState.BestValue = optStats.BestValue
				algState.BestReward = optStats.BestReward
				algState.LastError = optStats.LastError
				algState.Integral = optStats.Integral
				algState.StepSize = optStats.StepSize
				algState.Temperature = optStats.Temperature
			}

			state[key] = algState

			continue
		}

		variantsArr, ok := allVariantsMap[key]
		if !ok {
			slog.Warn("variants not found for feature",
				"feature", record.FeatureID, "env", record.EnvironmentID)

			continue
		}

		variants := make(map[string]*VariantStats, len(variantsArr))
		for _, variantKey := range variantsArr {
			stat := banditStatsMap[statsKey{
				FeatureID: record.FeatureID,
				EnvID:     record.EnvironmentID,
				AlgSlug:   algSlug,
				Variant:   variantKey,
			}]
			variants[variantKey] = &stat
		}

		if algKind == domain.AlgorithmKindContextualBandit {
			featureDim := int(getSettingAsFloat64(record.Settings, "feature_dim", DefaultFeatureDim))
			ctxState := NewContextualAlgorithmState(featureDim, variantsArr, record.Settings)

			for _, variantKey := range variantsArr {
				ctxStat, hasCtxStats := contextualStatsMap[statsKey{
					FeatureID: record.FeatureID,
					EnvID:     record.EnvironmentID,
					AlgSlug:   algSlug,
					Variant:   variantKey,
				}]
				if hasCtxStats && ctxState.Variants[variantKey] != nil {
					vs := ctxState.Variants[variantKey]
					if len(ctxStat.MatrixA) == featureDim*featureDim {
						vs.A = ctxStat.MatrixA
					}
					if len(ctxStat.VectorB) == featureDim {
						vs.B = ctxStat.VectorB
					}
					vs.Pulls = ctxStat.Pulls
					vs.TotalRew = ctxStat.TotalReward.InexactFloat64()
					vs.Successes = ctxStat.Successes
					vs.Failures = ctxStat.Failures
				}
			}

			state[key] = &AlgorithmState{
				AlgorithmType:   algType,
				Enabled:         record.Enabled,
				Variants:        variants,
				VariantsArr:     slices.Clone(variantsArr),
				Settings:        record.Settings,
				IsOptimizer:     false,
				IsContextual:    true,
				ContextualState: ctxState,
			}

			continue
		}

		state[key] = &AlgorithmState{
			AlgorithmType: algType,
			Enabled:       record.Enabled,
			Variants:      variants,
			VariantsArr:   slices.Clone(variantsArr),
			Settings:      record.Settings,
			IsOptimizer:   false,
			IsContextual:  false,
		}
	}

	m.mu.Lock()
	m.state = state
	m.lastSeen = now
	m.mu.Unlock()

	return nil
}

func (m *BanditManager) removeFeatureFromState(
	ctx context.Context,
	featureID domain.FeatureID,
	envKey string,
) error {
	feature, err := m.featuresUseCase.GetByIDWithEnv(ctx, featureID, envKey)
	if err != nil {
		return fmt.Errorf("find feature algorithms: %w", err)
	}

	key := StateKey{
		FeatureKey: feature.Key,
		EnvKey:     envKey,
	}

	m.mu.Lock()
	delete(m.state, key)
	m.mu.Unlock()

	return nil
}

//nolint:nestif // This is a complex function.
func (m *BanditManager) refreshFeature(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	envKey string,
) error {
	records, err := m.featureAlgorithmsRepo.ListExtendedByFeatureIDWithEnvID(ctx, featureID, envID)
	if err != nil {
		return fmt.Errorf("find feature algorithms: %w", err)
	}

	variants, err := m.featureVariantsRepo.ListByFeatureIDWithEnvID(ctx, featureID, envID)
	if err != nil {
		return fmt.Errorf("find feature variants: %w", err)
	}

	actualVariantsMap := make(map[string]*VariantStats, len(variants))
	actualVariantsArr := make([]string, 0, len(variants))
	for _, variant := range variants {
		actualVariantsArr = append(actualVariantsArr, variant.Name)
		actualVariantsMap[variant.Name] = &VariantStats{}
	}

	for _, record := range records {
		// Skip custom algorithms (handled by WASM manager)
		if record.AlgorithmSlug == nil {
			continue
		}

		algSlug := *record.AlgorithmSlug

		key := StateKey{
			FeatureKey: record.FeatureKey,
			EnvKey:     envKey,
		}

		algType := domain.AlgorithmSlugToType(algSlug)
		algKind := algType.Kind()

		state, ok := m.GetAlgorithmState(record.FeatureKey, envKey)
		if !ok {
			// add new
			switch algKind {
			case domain.AlgorithmKindBandit:
				state = &AlgorithmState{
					AlgorithmType: algType,
					Enabled:       record.Enabled,
					Variants:      actualVariantsMap,
					VariantsArr:   slices.Clone(actualVariantsArr),
					Settings:      record.Settings,
					IsOptimizer:   false,
					IsContextual:  false,
				}
			case domain.AlgorithmKindContextualBandit:
				featureDim := int(getSettingAsFloat64(record.Settings, "feature_dim", DefaultFeatureDim))
				state = &AlgorithmState{
					AlgorithmType:   algType,
					Enabled:         record.Enabled,
					Variants:        actualVariantsMap,
					VariantsArr:     slices.Clone(actualVariantsArr),
					Settings:        record.Settings,
					IsOptimizer:     false,
					IsContextual:    true,
					ContextualState: NewContextualAlgorithmState(featureDim, actualVariantsArr, record.Settings),
				}
			case domain.AlgorithmKindOptimizer:
				state = &AlgorithmState{
					AlgorithmType: algType,
					Enabled:       record.Enabled,
					Settings:      record.Settings,
					IsOptimizer:   true,
					Iteration:     0,
					CurrentValue:  getSettingAsDecimal(record.Settings, "initial_value", 0.0),
					MetricSum:     decimal.Zero,
					Variants:      make(map[string]*VariantStats),
					VariantsArr:   []string{},
				}
			case domain.AlgorithmKindUnknown:
				continue
			default:
				continue
			}

			m.mu.Lock()
			m.state[key] = state
			m.mu.Unlock()
		} else {
			// update existing
			state.mu.Lock()
			state.Enabled = record.Enabled
			state.Settings = record.Settings

			if algKind == domain.AlgorithmKindBandit || algKind == domain.AlgorithmKindContextualBandit {
				// Update variants for bandits
				for _, variantKey := range state.VariantsArr {
					if _, ok := actualVariantsMap[variantKey]; !ok {
						delete(state.Variants, variantKey)
					}
				}
				for _, variantKey := range actualVariantsArr {
					if _, ok := state.Variants[variantKey]; !ok {
						state.Variants[variantKey] = &VariantStats{}
					}
				}
				state.VariantsArr = slices.Clone(actualVariantsArr)
			}
			state.mu.Unlock()
		}
	}

	return nil
}

func (m *BanditManager) syncToDBLoop() {
	ticker := time.NewTicker(m.syncStatsInterval)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			if err := m.flushAllToDB(); err != nil {
				slog.Error("bandit: failed flush", "error", err)
			}
		}
	}
}

func (m *BanditManager) flushAllToDB() error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var banditRecords []domain.FeatureAlgorithmStats
	var optimizerRecords []domain.FeatureOptimizerStats
	var contextualRecords []domain.FeatureContextualStats

	m.mu.RLock()
	for feat, algState := range m.state {
		feature, err := m.featuresUseCase.GetByKeyWithEnvCached(ctx, feat.FeatureKey, feat.EnvKey)
		if err != nil {
			m.mu.RUnlock()

			return fmt.Errorf("get feature: %w", err)
		}

		env, err := m.envsUseCase.GetByIDCached(ctx, feature.EnvironmentID)
		if err != nil {
			m.mu.RUnlock()

			return fmt.Errorf("get env: %w", err)
		}

		algState.mu.RLock()

		switch {
		case algState.IsOptimizer:
			optimizerRecords = append(optimizerRecords, domain.FeatureOptimizerStats{
				ProjectID:      feature.ProjectID,
				EnvironmentID:  feature.EnvironmentID,
				FeatureID:      feature.ID,
				AlgorithmSlug:  algState.AlgorithmType.Slug(),
				FeatureKey:     feature.Key,
				EnvironmentKey: env.Key,
				Iteration:      algState.Iteration,
				CurrentValue:   algState.CurrentValue,
				BestValue:      algState.BestValue,
				BestReward:     algState.BestReward,
				MetricSum:      algState.MetricSum,
				LastError:      algState.LastError,
				Integral:       algState.Integral,
				StepSize:       algState.StepSize,
				Temperature:    algState.Temperature,
			})
		case algState.IsContextual && algState.ContextualState != nil:
			algState.ContextualState.mu.RLock()
			for variant, vs := range algState.ContextualState.Variants {
				contextualRecords = append(contextualRecords, domain.FeatureContextualStats{
					ProjectID:      feature.ProjectID,
					EnvironmentID:  feature.EnvironmentID,
					FeatureID:      feature.ID,
					AlgorithmSlug:  algState.AlgorithmType.Slug(),
					VariantKey:     variant,
					FeatureKey:     feature.Key,
					EnvironmentKey: env.Key,
					FeatureDim:     algState.ContextualState.FeatureDim,
					MatrixA:        vs.A,
					VectorB:        vs.B,
					Pulls:          vs.Pulls,
					TotalReward:    decimal.NewFromFloat(vs.TotalRew),
					Successes:      vs.Successes,
					Failures:       vs.Failures,
				})
			}
			algState.ContextualState.mu.RUnlock()
		default:
			for variant, variantStats := range algState.Variants {
				banditRecords = append(banditRecords, domain.FeatureAlgorithmStats{
					ProjectID:      feature.ProjectID,
					EnvironmentID:  feature.EnvironmentID,
					FeatureID:      feature.ID,
					AlgorithmSlug:  algState.AlgorithmType.Slug(),
					VariantKey:     variant,
					FeatureKey:     feature.Key,
					EnvironmentKey: env.Key,
					Evaluations:    variantStats.Evaluations,
					Successes:      variantStats.Successes,
					Failures:       variantStats.Failures,
					MetricSum:      variantStats.MetricSum,
				})
			}
		}

		algState.mu.RUnlock()
	}
	m.mu.RUnlock()

	if err := m.statsRepo.InsertBatch(ctx, banditRecords); err != nil {
		return fmt.Errorf("insert bandit batch: %w", err)
	}

	if err := m.optimizerStatsRepo.InsertBatch(ctx, optimizerRecords); err != nil {
		return fmt.Errorf("insert optimizer batch: %w", err)
	}

	if err := m.contextualStatsRepo.InsertBatch(ctx, contextualRecords); err != nil {
		return fmt.Errorf("insert contextual batch: %w", err)
	}

	// Flush custom WASM algorithm stats
	if err := m.flushCustomStats(); err != nil {
		return fmt.Errorf("insert custom algorithm batch: %w", err)
	}

	slog.Debug("bandit: flushed all features", "elapsed", time.Since(start))

	return nil
}

func getSettingAsFloat64(settings map[string]decimal.Decimal, key string, defaultValue float64) float64 {
	value, ok := settings[key]
	if !ok {
		return defaultValue
	}

	return value.InexactFloat64()
}

func getSettingAsDecimal(settings map[string]decimal.Decimal, key string, defaultValue float64) decimal.Decimal {
	value, ok := settings[key]
	if !ok {
		return decimal.NewFromFloat(defaultValue)
	}

	return value
}
