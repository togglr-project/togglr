package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/tetratelabs/wazero"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// StateKey identifies a custom algorithm instance for a specific feature/environment.
type StateKey struct {
	FeatureKey string
	EnvKey     string
}

// CustomAlgorithmState holds runtime state for a custom algorithm instance.
type CustomAlgorithmState struct {
	AlgorithmID  domain.CustomAlgorithmID
	Kind         domain.AlgorithmKind
	Enabled      bool
	Settings     map[string]decimal.Decimal
	WASMHash     string
	Variants     []string
	VariantStats map[string]*VariantStats
	State        json.RawMessage // persisted algorithm state

	// Optimizer-specific state
	CurrentValue float64
	Iteration    uint64
	MetricSum    float64
	BestValue    float64
	BestReward   float64

	mu sync.RWMutex
}

// VariantStats tracks per-variant statistics.
type VariantStats struct {
	Evaluations uint64
	Successes   uint64
	Failures    uint64
	MetricSum   float64
}

// Manager handles custom WASM algorithm execution and state management.
type Manager struct {
	runtime  *Runtime
	executor *Executor

	// State for each feature/env combination
	state map[StateKey]*CustomAlgorithmState
	mu    sync.RWMutex

	// Compiled modules cache (by algorithm ID)
	modules  map[domain.CustomAlgorithmID]wazero.CompiledModule
	moduleMu sync.RWMutex

	// Module instance pool (by algorithm ID)
	instances   map[domain.CustomAlgorithmID]*ModuleInstance
	instancesMu sync.RWMutex

	customAlgRepo   contract.CustomAlgorithmsRepository
	customStatsRepo contract.CustomAlgorithmStatsRepository
}

// NewManager creates a new WASM algorithm manager.
func NewManager(
	ctx context.Context,
	customAlgRepo contract.CustomAlgorithmsRepository,
	customStatsRepo contract.CustomAlgorithmStatsRepository,
) (*Manager, error) {
	runtime, err := NewRuntime(ctx)
	if err != nil {
		return nil, fmt.Errorf("create WASM runtime: %w", err)
	}

	return &Manager{
		runtime:         runtime,
		executor:        NewExecutor(runtime),
		state:           make(map[StateKey]*CustomAlgorithmState),
		modules:         make(map[domain.CustomAlgorithmID]wazero.CompiledModule),
		instances:       make(map[domain.CustomAlgorithmID]*ModuleInstance),
		customAlgRepo:   customAlgRepo,
		customStatsRepo: customStatsRepo,
	}, nil
}

// LoadState loads all custom algorithm states from database.
func (m *Manager) LoadState(ctx context.Context) error {
	// Load all custom algorithms
	algorithms, err := m.customAlgRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("list custom algorithms: %w", err)
	}

	// Pre-compile all modules
	for _, alg := range algorithms {
		if _, err := m.getOrCompileModule(ctx, alg); err != nil {
			slog.Error("failed to compile custom algorithm",
				"algorithm", alg.Slug, "error", err)
		}
	}

	// Load stats
	stats, err := m.customStatsRepo.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("load custom algorithm stats: %w", err)
	}

	// Build state map
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, stat := range stats {
		key := StateKey{
			FeatureKey: stat.FeatureKey,
			EnvKey:     stat.EnvironmentKey,
		}

		state, ok := m.state[key]
		if !ok {
			state = &CustomAlgorithmState{
				AlgorithmID:  stat.AlgorithmID,
				Variants:     []string{},
				VariantStats: make(map[string]*VariantStats),
				State:        stat.State,
			}
			m.state[key] = state
		}

		if stat.VariantKey != "" {
			state.Variants = append(state.Variants, stat.VariantKey)
			state.VariantStats[stat.VariantKey] = &VariantStats{
				Evaluations: stat.Evaluations,
				Successes:   stat.Successes,
				Failures:    stat.Failures,
				MetricSum:   stat.MetricSum.InexactFloat64(),
			}
		}
	}

	return nil
}

// getOrCompileModule gets a compiled module from cache or compiles it.
func (m *Manager) getOrCompileModule(ctx context.Context, alg domain.CustomAlgorithm) (wazero.CompiledModule, error) {
	m.moduleMu.RLock()
	if compiled, ok := m.modules[alg.ID]; ok {
		m.moduleMu.RUnlock()
		return compiled, nil
	}
	m.moduleMu.RUnlock()

	m.moduleMu.Lock()
	defer m.moduleMu.Unlock()

	// Double-check after acquiring write lock
	if compiled, ok := m.modules[alg.ID]; ok {
		return compiled, nil
	}

	compiled, err := m.runtime.CompileModule(ctx, alg.WASMHash, alg.WASMBinary)
	if err != nil {
		return nil, err
	}

	m.modules[alg.ID] = compiled
	return compiled, nil
}

// getOrCreateInstance gets or creates a module instance for an algorithm.
func (m *Manager) getOrCreateInstance(ctx context.Context, algID domain.CustomAlgorithmID) (*ModuleInstance, error) {
	m.instancesMu.RLock()
	if instance, ok := m.instances[algID]; ok {
		m.instancesMu.RUnlock()
		return instance, nil
	}
	m.instancesMu.RUnlock()

	m.instancesMu.Lock()
	defer m.instancesMu.Unlock()

	// Double-check
	if instance, ok := m.instances[algID]; ok {
		return instance, nil
	}

	// Get algorithm
	alg, err := m.customAlgRepo.GetByID(ctx, algID)
	if err != nil {
		return nil, fmt.Errorf("get algorithm: %w", err)
	}

	// Get compiled module
	compiled, err := m.getOrCompileModule(ctx, alg)
	if err != nil {
		return nil, err
	}

	// Create instance
	instance, err := m.executor.InstantiateModule(ctx, compiled, string(algID))
	if err != nil {
		return nil, err
	}

	m.instances[algID] = instance
	return instance, nil
}

// GetState returns the state for a feature/environment.
func (m *Manager) GetState(featureKey, envKey string) (*CustomAlgorithmState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, ok := m.state[StateKey{FeatureKey: featureKey, EnvKey: envKey}]
	return state, ok
}

// SetState sets the state for a feature/environment.
func (m *Manager) SetState(featureKey, envKey string, state *CustomAlgorithmState) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state[StateKey{FeatureKey: featureKey, EnvKey: envKey}] = state
}

// RemoveState removes state for a feature/environment.
func (m *Manager) RemoveState(featureKey, envKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.state, StateKey{FeatureKey: featureKey, EnvKey: envKey})
}

// EvaluateBandit evaluates a custom bandit algorithm.
func (m *Manager) EvaluateBandit(ctx context.Context, featureKey, envKey string) (string, bool) {
	state, ok := m.GetState(featureKey, envKey)
	if !ok || !state.Enabled {
		return "", false
	}

	instance, err := m.getOrCreateInstance(ctx, state.AlgorithmID)
	if err != nil {
		slog.Error("failed to get WASM instance", "error", err)
		return "", false
	}

	// Build input
	input := m.buildBanditInput(state, nil)

	// Execute
	output, err := m.executor.Evaluate(ctx, instance, input)
	if err != nil {
		slog.Error("WASM evaluate failed", "error", err)
		return "", false
	}

	// Update state
	state.mu.Lock()
	if output.NewState != nil {
		state.State = output.NewState
	}
	state.mu.Unlock()

	return output.SelectedVariant, output.SelectedVariant != ""
}

// EvaluateContextual evaluates a custom contextual bandit algorithm.
func (m *Manager) EvaluateContextual(ctx context.Context, featureKey, envKey string, userCtx map[string]any) (string, bool) {
	state, ok := m.GetState(featureKey, envKey)
	if !ok || !state.Enabled {
		return "", false
	}

	instance, err := m.getOrCreateInstance(ctx, state.AlgorithmID)
	if err != nil {
		slog.Error("failed to get WASM instance", "error", err)
		return "", false
	}

	// Build input
	input := m.buildBanditInput(state, userCtx)

	// Execute
	output, err := m.executor.Evaluate(ctx, instance, input)
	if err != nil {
		slog.Error("WASM evaluate failed", "error", err)
		return "", false
	}

	// Update state
	state.mu.Lock()
	if output.NewState != nil {
		state.State = output.NewState
	}
	state.mu.Unlock()

	return output.SelectedVariant, output.SelectedVariant != ""
}

// EvaluateOptimizer evaluates a custom optimizer algorithm.
func (m *Manager) EvaluateOptimizer(ctx context.Context, featureKey, envKey string) (decimal.Decimal, bool) {
	state, ok := m.GetState(featureKey, envKey)
	if !ok || !state.Enabled {
		return decimal.Zero, false
	}

	instance, err := m.getOrCreateInstance(ctx, state.AlgorithmID)
	if err != nil {
		slog.Error("failed to get WASM instance", "error", err)
		return decimal.Zero, false
	}

	// Build input
	input := m.buildOptimizerInput(state)

	// Execute
	output, err := m.executor.Evaluate(ctx, instance, input)
	if err != nil {
		slog.Error("WASM evaluate failed", "error", err)
		return decimal.Zero, false
	}

	// Update state
	state.mu.Lock()
	if output.NewState != nil {
		state.State = output.NewState
	}
	state.CurrentValue = output.OptimizedValue
	state.mu.Unlock()

	return decimal.NewFromFloat(output.OptimizedValue), true
}

// HandleFeedback processes a feedback event for a custom algorithm.
func (m *Manager) HandleFeedback(
	ctx context.Context,
	featureKey, envKey, variantKey string,
	eventType domain.FeedbackEventType,
	reward decimal.Decimal,
	userCtx map[string]any,
) {
	state, ok := m.GetState(featureKey, envKey)
	if !ok {
		return
	}

	instance, err := m.getOrCreateInstance(ctx, state.AlgorithmID)
	if err != nil {
		slog.Error("failed to get WASM instance for feedback", "error", err)
		return
	}

	if !instance.HasFeedbackHandler() {
		// No feedback handler, just update local stats
		m.updateLocalStats(state, variantKey, eventType, reward)
		return
	}

	// Build feedback input
	settings := make(map[string]float64, len(state.Settings))
	for k, v := range state.Settings {
		settings[k] = v.InexactFloat64()
	}

	input := &domain.WASMFeedbackInput{
		Settings:   settings,
		State:      state.State,
		VariantKey: variantKey,
		EventType:  string(eventType),
		Reward:     reward.InexactFloat64(),
		Context:    userCtx,
	}

	// Execute feedback handler
	output, err := m.executor.HandleFeedback(ctx, instance, input)
	if err != nil {
		slog.Error("WASM feedback handler failed", "error", err)
		m.updateLocalStats(state, variantKey, eventType, reward)
		return
	}

	// Update state
	state.mu.Lock()
	if output.NewState != nil {
		state.State = output.NewState
	}
	state.mu.Unlock()

	m.updateLocalStats(state, variantKey, eventType, reward)
}

// buildBanditInput creates input for bandit algorithms.
func (m *Manager) buildBanditInput(state *CustomAlgorithmState, userCtx map[string]any) *domain.WASMInput {
	state.mu.RLock()
	defer state.mu.RUnlock()

	settings := make(map[string]float64, len(state.Settings))
	for k, v := range state.Settings {
		settings[k] = v.InexactFloat64()
	}

	variantStats := make(map[string]domain.WASMVariantStats, len(state.VariantStats))
	for k, v := range state.VariantStats {
		variantStats[k] = domain.WASMVariantStats{
			Evaluations: v.Evaluations,
			Successes:   v.Successes,
			Failures:    v.Failures,
			MetricSum:   v.MetricSum,
		}
	}

	return &domain.WASMInput{
		Settings:     settings,
		State:        state.State,
		Variants:     state.Variants,
		VariantStats: variantStats,
		Context:      userCtx,
	}
}

// buildOptimizerInput creates input for optimizer algorithms.
func (m *Manager) buildOptimizerInput(state *CustomAlgorithmState) *domain.WASMInput {
	state.mu.RLock()
	defer state.mu.RUnlock()

	settings := make(map[string]float64, len(state.Settings))
	for k, v := range state.Settings {
		settings[k] = v.InexactFloat64()
	}

	return &domain.WASMInput{
		Settings:     settings,
		State:        state.State,
		CurrentValue: state.CurrentValue,
		Iteration:    state.Iteration,
		MetricSum:    state.MetricSum,
		BestValue:    state.BestValue,
		BestReward:   state.BestReward,
	}
}

// updateLocalStats updates in-memory statistics.
func (m *Manager) updateLocalStats(state *CustomAlgorithmState, variantKey string, eventType domain.FeedbackEventType, reward decimal.Decimal) {
	state.mu.Lock()
	defer state.mu.Unlock()

	if state.Kind == domain.AlgorithmKindOptimizer {
		state.Iteration++
		rewardFloat := reward.InexactFloat64()
		state.MetricSum += rewardFloat
		if rewardFloat > state.BestReward {
			state.BestReward = rewardFloat
			state.BestValue = state.CurrentValue
		}
		return
	}

	// Bandit stats
	vs, ok := state.VariantStats[variantKey]
	if !ok {
		vs = &VariantStats{}
		state.VariantStats[variantKey] = vs
		state.Variants = append(state.Variants, variantKey)
	}

	switch eventType {
	case domain.FeedbackEventTypeEvaluation:
		vs.Evaluations++
	case domain.FeedbackEventTypeSuccess:
		vs.Successes++
		vs.MetricSum += reward.InexactFloat64()
	case domain.FeedbackEventTypeFailure:
		vs.Failures++
	}
}

// FlushStats saves all state to database.
func (m *Manager) FlushStats(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var records []domain.CustomAlgorithmStats

	for key, state := range m.state {
		state.mu.RLock()

		if state.Kind == domain.AlgorithmKindOptimizer {
			// Single record for optimizer
			records = append(records, domain.CustomAlgorithmStats{
				AlgorithmID:    state.AlgorithmID,
				FeatureKey:     key.FeatureKey,
				EnvironmentKey: key.EnvKey,
				VariantKey:     "",
				State:          state.State,
				Evaluations:    state.Iteration,
				MetricSum:      decimal.NewFromFloat(state.MetricSum),
			})
		} else {
			// Per-variant records for bandits
			for variant, vs := range state.VariantStats {
				records = append(records, domain.CustomAlgorithmStats{
					AlgorithmID:    state.AlgorithmID,
					FeatureKey:     key.FeatureKey,
					EnvironmentKey: key.EnvKey,
					VariantKey:     variant,
					State:          state.State,
					Evaluations:    vs.Evaluations,
					Successes:      vs.Successes,
					Failures:       vs.Failures,
					MetricSum:      decimal.NewFromFloat(vs.MetricSum),
				})
			}
		}

		state.mu.RUnlock()
	}

	if len(records) == 0 {
		return nil
	}

	return m.customStatsRepo.UpsertBatch(ctx, records)
}

// InvalidateModule removes a module from cache (e.g., when algorithm is updated).
func (m *Manager) InvalidateModule(ctx context.Context, algID domain.CustomAlgorithmID, wasmHash string) {
	m.moduleMu.Lock()
	if compiled, ok := m.modules[algID]; ok {
		_ = compiled.Close(ctx)
		delete(m.modules, algID)
	}
	m.moduleMu.Unlock()

	m.instancesMu.Lock()
	if instance, ok := m.instances[algID]; ok {
		_ = instance.Close(ctx)
		delete(m.instances, algID)
	}
	m.instancesMu.Unlock()

	m.runtime.RemoveFromCache(wasmHash)
}

// Close releases all resources.
func (m *Manager) Close(ctx context.Context) error {
	m.instancesMu.Lock()
	for id, instance := range m.instances {
		_ = instance.Close(ctx)
		delete(m.instances, id)
	}
	m.instancesMu.Unlock()

	return m.runtime.Close(ctx)
}
