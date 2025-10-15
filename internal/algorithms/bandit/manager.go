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

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
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
	Variants      map[string]*VariantStats // variant_key -> stats
	VariantsArr   []string
	Settings      map[string]decimal.Decimal

	mu sync.RWMutex
}

type StateKey struct {
	FeatureID domain.FeatureID
	EnvID     domain.EnvironmentID
}

type BanditManager struct {
	mu           sync.RWMutex
	state        map[StateKey]*AlgorithmState
	randSrc      *rand.Rand
	syncInterval time.Duration
	stopCh       chan struct{}

	featureAlgorithmsRepo contract.FeatureAlgorithmsRepository
	featureVariantsRepo   contract.FlagVariantsRepository
	statsRepo             contract.FeatureAlgorithmStatsRepository
}

func New(
	featureAlgorithmsRepo contract.FeatureAlgorithmsRepository,
	featureVariantsRepo contract.FlagVariantsRepository,
	statsRepo contract.FeatureAlgorithmStatsRepository,
) (*BanditManager, error) {
	mngr := &BanditManager{
		state:                 make(map[StateKey]*AlgorithmState),
		randSrc:               rand.New(rand.NewSource(time.Now().UnixNano())),
		syncInterval:          time.Second * 5,
		featureAlgorithmsRepo: featureAlgorithmsRepo,
		featureVariantsRepo:   featureVariantsRepo,
		statsRepo:             statsRepo,
		stopCh:                make(chan struct{}),
	}

	if err := mngr.loadState(); err != nil {
		return nil, fmt.Errorf("load state: %w", err)
	}

	return mngr, nil
}

func (m *BanditManager) Start(context.Context) error {
	go m.syncToDBLoop() //nolint:contextcheck // false positive

	return nil
}

func (m *BanditManager) Stop(context.Context) error {
	close(m.stopCh)

	return nil
}

func (m *BanditManager) GetAlgorithmState(
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) (*AlgorithmState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := StateKey{FeatureID: featureID, EnvID: envID}
	state, ok := m.state[key]

	return state, ok
}

// EvaluateFeature chooses a variant according to the algorithm.
func (m *BanditManager) EvaluateFeature(
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) (string, bool) {
	state, ok := m.GetAlgorithmState(featureID, envID)
	if !ok {
		return "", false
	}

	if !state.Enabled {
		return "", false
	}

	var variant string

	// choose algorithm
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

// HandleTrackEvent called by track consumer to update in-memory counters.
func (m *BanditManager) HandleTrackEvent(
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	variantKey string,
	eventType domain.FeedbackEventType,
	metric decimal.Decimal,
) {
	state, ok := m.GetAlgorithmState(featureID, envID)
	if !ok {
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

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

func (m *BanditManager) loadState() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	records, err := m.featureAlgorithmsRepo.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("list feature algorithms: %w", err)
	}

	stats, err := m.statsRepo.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("load all stats: %w", err)
	}

	allVariants, err := m.featureVariantsRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("load all variants: %w", err)
	}

	allVariantsMap := make(map[StateKey][]string)
	for _, variant := range allVariants {
		key := StateKey{
			FeatureID: variant.FeatureID,
			EnvID:     variant.EnvironmentID,
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

	statsMap := make(map[statsKey]VariantStats)
	for _, stat := range stats {
		key := statsKey{
			FeatureID: stat.FeatureID,
			EnvID:     stat.EnvironmentID,
			AlgSlug:   stat.AlgorithmSlug,
			Variant:   stat.VariantKey,
		}
		statsMap[key] = VariantStats{
			Evaluations: stat.Evaluations,
			Successes:   stat.Successes,
			Failures:    stat.Failures,
			MetricSum:   stat.MetricSum,
		}
	}

	state := make(map[StateKey]*AlgorithmState, len(records))
	for _, record := range records {
		key := StateKey{
			FeatureID: record.FeatureID,
			EnvID:     record.EnvironmentID,
		}

		variantsArr, ok := allVariantsMap[key]
		if !ok {
			slog.Warn("variants not found for feature",
				"feature", record.FeatureID, "env", record.EnvironmentID)

			continue
		}

		variants := make(map[string]*VariantStats, len(variantsArr))
		for _, variantKey := range variantsArr {
			stat := statsMap[statsKey{
				FeatureID: record.FeatureID,
				EnvID:     record.EnvironmentID,
				AlgSlug:   record.AlgorithmSlug,
				Variant:   variantKey,
			}]
			variants[variantKey] = &stat
		}

		state[key] = &AlgorithmState{
			AlgorithmType: domain.AlgorithmSlugToType(record.AlgorithmSlug),
			Enabled:       record.Enabled,
			Variants:      variants,
			VariantsArr:   slices.Clone(variantsArr),
			Settings:      record.Settings,
		}
	}

	m.mu.Lock()
	m.state = state
	m.mu.Unlock()

	return nil
}

func (m *BanditManager) syncToDBLoop() {
	ticker := time.NewTicker(m.syncInterval)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopCh:
			if err := m.flushAllToDB(); err != nil {
				slog.Error("bandit: failed flush", "error", err)
			}

			return
		case <-ticker.C:
			if err := m.flushAllToDB(); err != nil {
				slog.Error("bandit: failed flush", "error", err)
			}
		}
	}
}

func (m *BanditManager) flushAllToDB() error {
	var records []domain.FeatureAlgorithmStats

	m.mu.RLock()
	for feat, algState := range m.state {
		algState.mu.RLock()
		for variant, variantStats := range algState.Variants {
			records = append(records, domain.FeatureAlgorithmStats{
				FeatureID:     feat.FeatureID,
				EnvironmentID: feat.EnvID,
				AlgorithmSlug: algState.AlgorithmType.Slug(),
				VariantKey:    variant,
				Evaluations:   variantStats.Evaluations,
				Successes:     variantStats.Successes,
				Failures:      variantStats.Failures,
				MetricSum:     variantStats.MetricSum,
			})
		}
		algState.mu.RUnlock()
	}
	m.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return m.statsRepo.InsertBatch(ctx, records)
}
