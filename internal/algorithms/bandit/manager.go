package bandit

import (
	"context"
	"errors"
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

	featureAlgorithmsRepo contract.FeatureAlgorithmsRepository
	featureVariantsRepo   contract.FlagVariantsRepository
	statsRepo             contract.FeatureAlgorithmStatsRepository
}

func New(
	featureAlgorithmsRepo contract.FeatureAlgorithmsRepository,
	featureVariantsRepo contract.FlagVariantsRepository,
	statsRepo contract.FeatureAlgorithmStatsRepository,
	syncInterval time.Duration,
) (*BanditManager, error) {
	mngr := &BanditManager{
		state:                 make(map[StateKey]*AlgorithmState),
		randSrc:               rand.New(rand.NewSource(time.Now().UnixNano())),
		syncInterval:          syncInterval,
		featureAlgorithmsRepo: featureAlgorithmsRepo,
		featureVariantsRepo:   featureVariantsRepo,
		statsRepo:             statsRepo,
	}

	if err := mngr.loadState(); err != nil {
		return nil, fmt.Errorf("load state: %w", err)
	}

	return mngr, nil
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

// EvaluateFeature chooses a variant according to algorithm.
func (m *BanditManager) EvaluateFeature(
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) (string, bool, error) {
	state, ok := m.GetAlgorithmState(featureID, envID)
	if !ok {
		return "", false, errors.New("no such feature")
	}

	if !state.Enabled {
		return "", false, nil
	}

	// choose algorithm
	var (
		variant string
		err     error
	)
	switch state.AlgorithmType {
	case domain.AlgorithmTypeEpsilonGreedy:
		variant, err = m.evalEpsilonGreedy(state)
	case domain.AlgorithmTypeThompsonSampling:
		variant, err = m.evalThompson(state)
	case domain.AlgorithmTypeUCB:
		variant, err = m.evalUCB(state)
	case domain.AlgorithmTypeUnknown:
		return "", false, errors.New("unknown algorithm")
	default:
		return "", false, fmt.Errorf("unknown algorithm: %v", state.AlgorithmType)
	}

	return variant, true, err
}

// HandleTrackEvent called by track consumer to update in-memory counters.
func (m *BanditManager) HandleTrackEvent(
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	variantKey string,
	eventType string,
	metric decimal.Decimal,
) error {
	state, ok := m.GetAlgorithmState(featureID, envID)
	if !ok {
		return errors.New("unknown state")
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	vs, ok := state.Variants[variantKey]
	if !ok {
		vs = &VariantStats{}
		state.Variants[variantKey] = vs
	}
	switch eventType {
	case "evaluation":
		vs.Evaluations++
	case "success":
		vs.Successes++
		vs.MetricSum = vs.MetricSum.Add(metric)
	case "failure":
		vs.Failures++
	case "error":
		vs.Failures++
	default:
		// custom handling
	}

	return nil
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
	defer m.mu.Unlock()
	m.state = state

	return nil
}

//// SyncToDBLoop periodically flushes in-memory counters to DB (UPSERT)
// func (m *BanditManager) SyncToDBLoop(ctx context.Context) {
//	ticker := time.NewTicker(m.syncInterval)
//	defer ticker.Stop()
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case <-ticker.C:
//			if err := m.flushAllToDB(ctx); err != nil {
//				fmt.Printf("bandit: failed flush: %v\n", err)
//			}
//		}
//	}
//}
//
// func (m *BanditManager) flushAllToDB(ctx context.Context) error {
//	type dumpKey struct {
//		FeatureID domain.FeatureID
//		EnvID     domain.EnvironmentID
//		AlgSlug   string
//	}
//
//	m.mu.RLock()
//	snapshot := make(map[dumpKey]map[string]VariantStats)
//	for feat, algState := range m.state {
//		snapKey := dumpKey{
//			FeatureID: feat.FeatureID,
//			EnvID:     feat.EnvID,
//			AlgSlug:   algState.AlgorithmType.Slug(),
//		}
//		snapshot[snapKey] = make(map[string]VariantStats)
//
//		algState.mu.RLock()
//		for variant, variantStats := range algState.Variants {
//			snapshot[snapKey][variant] = *variantStats
//		}
//		algState.mu.RUnlock()
//	}
//	m.mu.RUnlock()
//
//	batch := &pgx.Batch{}
//	query := `INSERT INTO monitoring.feature_algorithm_stats
//      (feature_id, environment_id, algorithm_slug, variant_key,
//     evaluations, successes, failures, metric_sum, updated_at)
//      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())
//      ON CONFLICT (feature_id, environment_id, algorithm_slug, variant_key) DO UPDATE SET
//        evaluations = monitoring.feature_algorithm_stats.evaluations + EXCLUDED.evaluations,
//        successes   = monitoring.feature_algorithm_stats.successes + EXCLUDED.successes,
//        failures    = monitoring.feature_algorithm_stats.failures + EXCLUDED.failures,
//        metric_sum  = monitoring.feature_algorithm_stats.metric_sum + EXCLUDED.metric_sum,
//        updated_at  = NOW();`
//
//	for feat, vsmap := range snapshot {
//		for vk, vs := range vsmap {
//			// skip zeros
//			if vs.Evaluations == 0 && vs.Successes == 0 && vs.Failures == 0 && vs.MetricSum.IsZero() {
//				continue
//			}
//			batch.Queue(query, feat.FeatureID, feat.EnvID, feat.AlgSlug, vk,
//			vs.Evaluations, vs.Successes, vs.Failures, vs.MetricSum)
//		}
//	}
//
//	if batch.Len() == 0 {
//		return nil
//	}
//
//	// send batch using m.db (pgx.Conn)
//	batchResults := m.db.SendBatch(ctx, batch)
//	defer batchResults.Close()
//	for i := 0; i < batch.Len(); i++ {
//		if _, err := batchResults.Exec(); err != nil {
//			return err
//		}
//	}
//	return nil
//}
