package bandit

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/algorithms/wasm"
	"github.com/togglr-project/togglr/internal/domain"
)

// customState tracks which features use custom algorithms.
type customState struct {
	AlgorithmID domain.CustomAlgorithmID
	Kind        domain.AlgorithmKind
	Enabled     bool
}

// SetWASMManager sets the WASM manager for custom algorithm support.
func (m *BanditManager) SetWASMManager(wasmMgr *wasm.Manager) {
	m.wasmManager = wasmMgr
}

// IsCustomAlgorithm checks if the feature uses a custom WASM algorithm.
func (m *BanditManager) IsCustomAlgorithm(featureKey, envKey string) bool {
	m.customMu.RLock()
	defer m.customMu.RUnlock()

	key := StateKey{FeatureKey: featureKey, EnvKey: envKey}
	state, ok := m.customState[key]
	return ok && state.Enabled
}

// GetCustomAlgorithmKind returns the kind of custom algorithm for a feature.
func (m *BanditManager) GetCustomAlgorithmKind(featureKey, envKey string) (domain.AlgorithmKind, bool) {
	m.customMu.RLock()
	defer m.customMu.RUnlock()

	key := StateKey{FeatureKey: featureKey, EnvKey: envKey}
	state, ok := m.customState[key]
	if !ok || !state.Enabled {
		return domain.AlgorithmKindUnknown, false
	}
	return state.Kind, true
}

// EvaluateCustom evaluates a custom bandit algorithm.
func (m *BanditManager) EvaluateCustom(featureKey, envKey string, ctx map[string]any) (string, bool) {
	if m.wasmManager == nil {
		return "", false
	}

	kind, ok := m.GetCustomAlgorithmKind(featureKey, envKey)
	if !ok {
		return "", false
	}

	switch kind {
	case domain.AlgorithmKindBandit:
		return m.wasmManager.EvaluateBandit(context.Background(), featureKey, envKey)
	case domain.AlgorithmKindContextualBandit:
		return m.wasmManager.EvaluateContextual(context.Background(), featureKey, envKey, ctx)
	default:
		return "", false
	}
}

// EvaluateCustomOptimizer evaluates a custom optimizer algorithm.
func (m *BanditManager) EvaluateCustomOptimizer(featureKey, envKey string) (decimal.Decimal, bool) {
	if m.wasmManager == nil {
		return decimal.Zero, false
	}

	kind, ok := m.GetCustomAlgorithmKind(featureKey, envKey)
	if !ok || kind != domain.AlgorithmKindOptimizer {
		return decimal.Zero, false
	}

	return m.wasmManager.EvaluateOptimizer(context.Background(), featureKey, envKey)
}

// HandleCustomTrackEvent handles feedback for custom algorithms.
func (m *BanditManager) HandleCustomTrackEvent(
	featureKey string,
	envKey string,
	variantKey string,
	eventType domain.FeedbackEventType,
	metric decimal.Decimal,
	ctx map[string]any,
) {
	if m.wasmManager == nil {
		return
	}

	m.wasmManager.HandleFeedback(
		context.Background(),
		featureKey,
		envKey,
		variantKey,
		eventType,
		metric,
		ctx,
	)
}

// RegisterCustomAlgorithm registers a custom algorithm for a feature.
func (m *BanditManager) RegisterCustomAlgorithm(
	featureKey, envKey string,
	algID domain.CustomAlgorithmID,
	kind domain.AlgorithmKind,
	enabled bool,
	variants []string,
	settings map[string]decimal.Decimal,
) {
	m.customMu.Lock()
	defer m.customMu.Unlock()

	key := StateKey{FeatureKey: featureKey, EnvKey: envKey}
	m.customState[key] = &customState{
		AlgorithmID: algID,
		Kind:        kind,
		Enabled:     enabled,
	}

	// Also register in WASM manager
	if m.wasmManager != nil {
		state := &wasm.CustomAlgorithmState{
			AlgorithmID:  algID,
			Kind:         kind,
			Enabled:      enabled,
			Settings:     settings,
			Variants:     variants,
			VariantStats: make(map[string]*wasm.VariantStats),
		}
		for _, v := range variants {
			state.VariantStats[v] = &wasm.VariantStats{}
		}
		m.wasmManager.SetState(featureKey, envKey, state)
	}
}

// UnregisterCustomAlgorithm removes a custom algorithm for a feature.
func (m *BanditManager) UnregisterCustomAlgorithm(featureKey, envKey string) {
	m.customMu.Lock()
	defer m.customMu.Unlock()

	key := StateKey{FeatureKey: featureKey, EnvKey: envKey}
	delete(m.customState, key)

	if m.wasmManager != nil {
		m.wasmManager.RemoveState(featureKey, envKey)
	}
}

// flushCustomStats saves custom algorithm stats to database.
func (m *BanditManager) flushCustomStats() error {
	if m.wasmManager == nil {
		return nil
	}
	return m.wasmManager.FlushStats(context.Background())
}
