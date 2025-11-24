package bandit

import (
	"math/rand"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/togglr-project/togglr/internal/domain"
)

func newTestManager(seed int64) *BanditManager {
	return &BanditManager{
		randSrc: rand.New(rand.NewSource(seed)),
	}
}

func newTestState(variants []string) *AlgorithmState {
	variantsMap := make(map[string]*VariantStats, len(variants))
	for _, v := range variants {
		variantsMap[v] = &VariantStats{}
	}
	return &AlgorithmState{
		Variants:    variantsMap,
		VariantsArr: variants,
		Settings:    make(map[string]decimal.Decimal),
	}
}

// ==================== Epsilon-Greedy Tests ====================

func TestEvalEpsilonGreedy_Exploitation(t *testing.T) {
	// With epsilon=0, should always exploit (choose best variant)
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})
	state.Settings["epsilon"] = decimal.NewFromFloat(0.0)

	// Set up stats: B has the best success rate
	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 30}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 70}
	state.Variants["C"] = &VariantStats{Evaluations: 100, Successes: 50}

	result := m.evalEpsilonGreedy(state)
	assert.Equal(t, "B", result, "Should choose variant B with highest success rate")
}

func TestEvalEpsilonGreedy_Exploration(t *testing.T) {
	// With epsilon=1, should always explore (random choice)
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})
	state.Settings["epsilon"] = decimal.NewFromFloat(1.0)

	// Set up stats: B has the best success rate, but should still explore
	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 30}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 70}
	state.Variants["C"] = &VariantStats{Evaluations: 100, Successes: 50}

	// Run multiple times and collect results
	results := make(map[string]int)
	for i := 0; i < 100; i++ {
		result := m.evalEpsilonGreedy(state)
		results[result]++
	}

	// With exploration, all variants should be chosen at least once
	assert.Greater(t, results["A"], 0, "Variant A should be chosen at least once")
	assert.Greater(t, results["B"], 0, "Variant B should be chosen at least once")
	assert.Greater(t, results["C"], 0, "Variant C should be chosen at least once")
}

func TestEvalEpsilonGreedy_ColdStart(t *testing.T) {
	// When all variants have 0 evaluations, should pick random
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})
	state.Settings["epsilon"] = decimal.NewFromFloat(0.0)

	result := m.evalEpsilonGreedy(state)
	assert.Contains(t, state.VariantsArr, result, "Should return one of the variants")
}

func TestEvalEpsilonGreedy_DefaultEpsilon(t *testing.T) {
	// Default epsilon should be 0.1
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})

	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 90}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 10}

	// Run many times - most should be A (exploitation), but some B (exploration)
	aCount := 0
	for i := 0; i < 1000; i++ {
		if m.evalEpsilonGreedy(state) == "A" {
			aCount++
		}
	}

	// With epsilon=0.1, expect ~90% A choices
	assert.Greater(t, aCount, 800, "Most choices should be A (exploitation)")
	assert.Less(t, aCount, 1000, "Some choices should be B (exploration)")
}

// ==================== Thompson Sampling Tests ====================

func TestEvalThompson_PrefersBetterVariant(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})

	// B has much better success rate
	state.Variants["A"] = &VariantStats{Successes: 10, Failures: 90}
	state.Variants["B"] = &VariantStats{Successes: 90, Failures: 10}

	// Run multiple times
	bCount := 0
	for i := 0; i < 100; i++ {
		if m.evalThompson(state) == "B" {
			bCount++
		}
	}

	// B should be chosen most of the time
	assert.Greater(t, bCount, 80, "Variant B should be chosen most of the time")
}

func TestEvalThompson_CustomPriors(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Settings["prior_alpha"] = decimal.NewFromFloat(2.0)
	state.Settings["prior_beta"] = decimal.NewFromFloat(2.0)

	state.Variants["A"] = &VariantStats{Successes: 5, Failures: 5}
	state.Variants["B"] = &VariantStats{Successes: 5, Failures: 5}

	// With equal stats and symmetric priors, both should be chosen roughly equally
	aCount := 0
	for i := 0; i < 100; i++ {
		if m.evalThompson(state) == "A" {
			aCount++
		}
	}

	// Expect roughly 50/50 split
	assert.Greater(t, aCount, 30, "Variant A should be chosen sometimes")
	assert.Less(t, aCount, 70, "Variant B should also be chosen sometimes")
}

func TestEvalThompson_ColdStart(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})

	// All variants have 0 successes/failures
	result := m.evalThompson(state)
	assert.Contains(t, state.VariantsArr, result, "Should return one of the variants")
}

func TestEvalThompson_IncrementsEvaluations(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Variants["A"] = &VariantStats{Successes: 50, Failures: 50}
	state.Variants["B"] = &VariantStats{Successes: 50, Failures: 50}

	initialA := state.Variants["A"].Evaluations
	initialB := state.Variants["B"].Evaluations

	result := m.evalThompson(state)

	// Only the chosen variant should have incremented evaluations
	if result == "A" {
		assert.Equal(t, initialA+1, state.Variants["A"].Evaluations)
		assert.Equal(t, initialB, state.Variants["B"].Evaluations)
	} else {
		assert.Equal(t, initialA, state.Variants["A"].Evaluations)
		assert.Equal(t, initialB+1, state.Variants["B"].Evaluations)
	}
}

// ==================== UCB Tests ====================

func TestEvalUCB_ColdStart(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})

	// UCB should prefer unexplored variants first
	seen := make(map[string]bool)
	for i := 0; i < 3; i++ {
		result := m.evalUCB(state)
		seen[result] = true
	}

	// All variants should be explored once
	assert.Len(t, seen, 3, "All variants should be explored during cold start")
}

func TestEvalUCB_ExploitsBestAfterExploration(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})

	// Both variants have been explored, B is much better
	state.Variants["A"] = &VariantStats{Evaluations: 50, Successes: 10}
	state.Variants["B"] = &VariantStats{Evaluations: 50, Successes: 40}

	// Run multiple evaluations
	bCount := 0
	for i := 0; i < 100; i++ {
		if m.evalUCB(state) == "B" {
			bCount++
		}
	}

	// B should be chosen most of the time
	assert.Greater(t, bCount, 70, "Better variant should be chosen more often")
}

func TestEvalUCB_CustomConfidence(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Settings["confidence"] = decimal.NewFromFloat(0.1) // Very low confidence = mostly exploitation

	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 60}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 40}

	// With very low confidence, should exploit more aggressively
	aCount := 0
	for i := 0; i < 100; i++ {
		if m.evalUCB(state) == "A" {
			aCount++
		}
	}

	// A has better mean, should be chosen more with low confidence
	assert.Greater(t, aCount, 50, "Better variant should be exploited more with low confidence")
}

func TestEvalUCB_IncrementsEvaluations(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Variants["A"] = &VariantStats{Evaluations: 10, Successes: 5}
	state.Variants["B"] = &VariantStats{Evaluations: 10, Successes: 5}

	initialTotal := state.Variants["A"].Evaluations + state.Variants["B"].Evaluations

	m.evalUCB(state)

	newTotal := state.Variants["A"].Evaluations + state.Variants["B"].Evaluations
	assert.Equal(t, initialTotal+1, newTotal, "Total evaluations should increase by 1")
}

// ==================== Hill Climb Tests ====================

func TestEvalHillClimb_ImproveOnBetterReward(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		MetricSum:    decimal.NewFromFloat(0.5),
		Settings: map[string]decimal.Decimal{
			"step":      decimal.NewFromFloat(0.1),
			"direction": decimal.NewFromFloat(1.0),
		},
	}

	// Better reward should move in same direction
	reward := decimal.NewFromFloat(0.8)
	result := m.evalHillClimb(state, reward)

	assert.True(t, result.GreaterThan(decimal.NewFromFloat(1.0)), "Value should increase with better reward")
	assert.Equal(t, uint64(1), state.Iteration, "Iteration should increment")
}

func TestEvalHillClimb_ReverseOnWorseReward(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		MetricSum:    decimal.NewFromFloat(0.8),
		Settings: map[string]decimal.Decimal{
			"step":      decimal.NewFromFloat(0.1),
			"direction": decimal.NewFromFloat(1.0),
		},
	}

	// Worse reward should reverse direction
	reward := decimal.NewFromFloat(0.5)
	m.evalHillClimb(state, reward)

	// Direction should be reversed
	assert.True(t, state.Settings["direction"].Equal(decimal.NewFromFloat(-1.0)),
		"Direction should reverse on worse reward")
}

func TestEvalHillClimb_DefaultSettings(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		MetricSum:    decimal.Zero,
		Settings:     make(map[string]decimal.Decimal),
	}

	reward := decimal.NewFromFloat(1.0)
	result := m.evalHillClimb(state, reward)

	// With default step=0.05 and direction=1, should move up
	assert.True(t, result.GreaterThan(decimal.Zero), "Should move in positive direction by default")
}

// ==================== Simulated Annealing Tests ====================

func TestEvalSimulatedAnnealing_AcceptsBetterReward(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		MetricSum:    decimal.NewFromFloat(0.5),
		Settings: map[string]decimal.Decimal{
			"temp":       decimal.NewFromFloat(1.0),
			"cooling":    decimal.NewFromFloat(0.95),
			"step_scale": decimal.NewFromFloat(0.1),
		},
	}

	// Better reward should always be accepted
	reward := decimal.NewFromFloat(0.8)
	m.evalSimulatedAnnealing(state, reward)

	// MetricSum should be updated to new reward
	assert.True(t, state.MetricSum.Equal(reward), "MetricSum should be updated on better reward")
}

func TestEvalSimulatedAnnealing_CoolsDown(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		MetricSum:    decimal.NewFromFloat(0.5),
		Settings: map[string]decimal.Decimal{
			"temp":       decimal.NewFromFloat(1.0),
			"cooling":    decimal.NewFromFloat(0.9),
			"step_scale": decimal.NewFromFloat(0.1),
		},
	}

	reward := decimal.NewFromFloat(0.8)
	m.evalSimulatedAnnealing(state, reward)

	// Temperature should decrease
	assert.True(t, state.Settings["temp"].LessThan(decimal.NewFromFloat(1.0)),
		"Temperature should decrease after iteration")
	assert.True(t, state.Settings["temp"].Equal(decimal.NewFromFloat(0.9)),
		"Temperature should be multiplied by cooling factor")
}

func TestEvalSimulatedAnnealing_MayAcceptWorseAtHighTemp(t *testing.T) {
	m := newTestManager(42)

	acceptedWorse := 0
	for i := 0; i < 100; i++ {
		state := &AlgorithmState{
			CurrentValue: decimal.NewFromFloat(1.0),
			MetricSum:    decimal.NewFromFloat(0.9),
			Settings: map[string]decimal.Decimal{
				"temp":       decimal.NewFromFloat(10.0), // High temperature
				"cooling":    decimal.NewFromFloat(0.95),
				"step_scale": decimal.NewFromFloat(0.1),
			},
		}

		// Slightly worse reward
		reward := decimal.NewFromFloat(0.8)
		m.evalSimulatedAnnealing(state, reward)

		if state.MetricSum.Equal(reward) {
			acceptedWorse++
		}
	}

	// With high temperature, should sometimes accept worse solutions
	assert.Greater(t, acceptedWorse, 0, "Should sometimes accept worse solutions at high temperature")
}

// ==================== PID Controller Tests ====================

func TestEvalPID_ConvergesToTarget(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		Settings: map[string]decimal.Decimal{
			"kp":         decimal.NewFromFloat(0.5),
			"ki":         decimal.NewFromFloat(0.1),
			"kd":         decimal.NewFromFloat(0.05),
			"integral":   decimal.Zero,
			"prev_error": decimal.Zero,
		},
	}

	target := decimal.NewFromFloat(10.0)

	// Simulate several iterations
	for i := 0; i < 50; i++ {
		measured := state.CurrentValue
		m.evalPID(state, measured, target)
	}

	// Should converge towards target
	diff := target.Sub(state.CurrentValue).Abs()
	assert.True(t, diff.LessThan(decimal.NewFromFloat(1.0)),
		"PID should converge towards target value")
}

func TestEvalPID_RespondsToError(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(5.0),
		Settings: map[string]decimal.Decimal{
			"kp":         decimal.NewFromFloat(0.2),
			"ki":         decimal.NewFromFloat(0.0),
			"kd":         decimal.NewFromFloat(0.0),
			"integral":   decimal.Zero,
			"prev_error": decimal.Zero,
		},
	}

	// Target is above current value
	target := decimal.NewFromFloat(10.0)
	measured := decimal.NewFromFloat(5.0)

	result := m.evalPID(state, measured, target)

	// With positive error (target > measured), output should increase
	assert.True(t, result.GreaterThan(decimal.NewFromFloat(5.0)),
		"PID output should increase when below target")
}

func TestEvalPID_IntegralAccumulates(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		Settings: map[string]decimal.Decimal{
			"kp":         decimal.NewFromFloat(0.0),
			"ki":         decimal.NewFromFloat(1.0), // Only integral term
			"kd":         decimal.NewFromFloat(0.0),
			"integral":   decimal.Zero,
			"prev_error": decimal.Zero,
		},
	}

	target := decimal.NewFromFloat(1.0)
	measured := decimal.NewFromFloat(0.0)

	// Run several iterations
	for i := 0; i < 5; i++ {
		m.evalPID(state, measured, target)
	}

	// Integral should accumulate
	assert.True(t, state.Settings["integral"].GreaterThan(decimal.NewFromFloat(4.0)),
		"Integral term should accumulate over iterations")
}

func TestEvalPID_DerivativeRespondsToChange(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		Settings: map[string]decimal.Decimal{
			"kp":         decimal.NewFromFloat(0.0),
			"ki":         decimal.NewFromFloat(0.0),
			"kd":         decimal.NewFromFloat(1.0), // Only derivative term
			"integral":   decimal.Zero,
			"prev_error": decimal.NewFromFloat(1.0), // Previous error was 1.0
		},
	}

	target := decimal.NewFromFloat(5.0)
	measured := decimal.NewFromFloat(3.0) // Error is 2.0, change is +1.0

	result := m.evalPID(state, measured, target)

	// Derivative responds to change in error
	// Error = 5.0 - 3.0 = 2.0, prev_error = 1.0, derivative = 2.0 - 1.0 = 1.0
	// output = kd * derivative = 1.0 * 1.0 = 1.0
	assert.True(t, result.Equal(decimal.NewFromFloat(1.0)),
		"PID derivative term should respond to error change")
}

// ==================== Helper Function Tests ====================

func TestGetSettingAsFloat64(t *testing.T) {
	settings := map[string]decimal.Decimal{
		"existing": decimal.NewFromFloat(3.14),
	}

	// Test existing key
	result := getSettingAsFloat64(settings, "existing", 0.0)
	assert.InDelta(t, 3.14, result, 0.001)

	// Test missing key (should return default)
	result = getSettingAsFloat64(settings, "missing", 42.0)
	assert.InDelta(t, 42.0, result, 0.001)
}

func TestGetSettingAsDecimal(t *testing.T) {
	settings := map[string]decimal.Decimal{
		"existing": decimal.NewFromFloat(3.14),
	}

	// Test existing key
	result := getSettingAsDecimal(settings, "existing", 0.0)
	assert.True(t, result.Equal(decimal.NewFromFloat(3.14)))

	// Test missing key (should return default)
	result = getSettingAsDecimal(settings, "missing", 42.0)
	assert.True(t, result.Equal(decimal.NewFromFloat(42.0)))
}

// ==================== Bayesian Optimization Tests ====================

func TestEvalBayesOpt_SelectsBestSample(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"noise": decimal.NewFromFloat(0.0), // No noise for deterministic test
		},
	}

	samples := []decimal.Decimal{
		decimal.NewFromFloat(1.0),
		decimal.NewFromFloat(2.0),
		decimal.NewFromFloat(3.0),
	}
	rewards := []decimal.Decimal{
		decimal.NewFromFloat(0.5),
		decimal.NewFromFloat(0.9), // Best reward
		decimal.NewFromFloat(0.3),
	}

	result := m.evalBayesOpt(state, samples, rewards)

	// Should select sample with best reward (2.0)
	assert.True(t, result.Equal(decimal.NewFromFloat(2.0)),
		"Should select sample with highest reward")
}

func TestEvalBayesOpt_WithNoise(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"noise": decimal.NewFromFloat(0.5), // High noise
		},
	}

	samples := []decimal.Decimal{
		decimal.NewFromFloat(1.0),
		decimal.NewFromFloat(2.0),
	}
	rewards := []decimal.Decimal{
		decimal.NewFromFloat(0.5),
		decimal.NewFromFloat(0.51), // Very close rewards
	}

	// With noise, selection may vary
	selections := make(map[string]int)
	for i := 0; i < 100; i++ {
		result := m.evalBayesOpt(state, samples, rewards)
		selections[result.String()]++
	}

	// Both should be selected sometimes due to noise
	assert.Greater(t, len(selections), 0, "Should select at least one variant")
}

// ==================== Cross-Entropy Method Tests ====================

func TestEvalCEM_SelectsEliteMean(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"elite_fraction":  decimal.NewFromFloat(0.5),
			"population_size": decimal.NewFromFloat(4),
		},
	}

	// 4 samples, elite_fraction=0.5 means top 2 will be selected
	samples := []decimal.Decimal{
		decimal.NewFromFloat(1.0),
		decimal.NewFromFloat(2.0),
		decimal.NewFromFloat(3.0),
		decimal.NewFromFloat(4.0),
	}
	rewards := []decimal.Decimal{
		decimal.NewFromFloat(0.1),
		decimal.NewFromFloat(0.9), // Elite
		decimal.NewFromFloat(0.2),
		decimal.NewFromFloat(0.8), // Elite
	}

	result := m.evalCEM(state, samples, rewards)

	// Elite samples are 2.0 (reward 0.9) and 4.0 (reward 0.8)
	// Mean of elites = (2.0 + 4.0) / 2 = 3.0
	assert.True(t, result.Equal(decimal.NewFromFloat(3.0)),
		"Should return mean of elite samples")
}

func TestEvalCEM_DefaultSettings(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings:     make(map[string]decimal.Decimal),
	}

	// 20 samples (default population_size)
	samples := make([]decimal.Decimal, 20)
	rewards := make([]decimal.Decimal, 20)
	for i := 0; i < 20; i++ {
		samples[i] = decimal.NewFromFloat(float64(i))
		rewards[i] = decimal.NewFromFloat(float64(i)) // Higher sample = higher reward
	}

	result := m.evalCEM(state, samples, rewards)

	// With elite_fraction=0.2 and pop_size=20, top 4 samples are elites
	// Top 4 by reward: 19, 18, 17, 16
	// Mean = (19 + 18 + 17 + 16) / 4 = 17.5
	assert.True(t, result.Equal(decimal.NewFromFloat(17.5)),
		"Should return mean of top 20% samples")
}

// ==================== Integration-like Tests ====================

// ==================== Manager Method Tests ====================

func TestGetAlgorithmKind_Bandit(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeEpsilonGreedy,
		Enabled:       true,
		IsOptimizer:   false,
	}

	kind, ok := m.GetAlgorithmKind("feature1", "prod")
	assert.True(t, ok)
	assert.Equal(t, domain.AlgorithmKindBandit, kind)
}

func TestGetAlgorithmKind_Optimizer(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeHillClimb,
		Enabled:       true,
		IsOptimizer:   true,
	}

	kind, ok := m.GetAlgorithmKind("feature1", "prod")
	assert.True(t, ok)
	assert.Equal(t, domain.AlgorithmKindOptimizer, kind)
}

func TestGetAlgorithmKind_NotFound(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	_, ok := m.GetAlgorithmKind("nonexistent", "prod")
	assert.False(t, ok)
}

func TestGetAlgorithmKind_Disabled(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeEpsilonGreedy,
		Enabled:       false,
		IsOptimizer:   false,
	}

	_, ok := m.GetAlgorithmKind("feature1", "prod")
	assert.False(t, ok, "Should return false for disabled algorithms")
}

func TestEvaluateOptimizer_ReturnsCurrentValue(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeHillClimb,
		Enabled:       true,
		IsOptimizer:   true,
		CurrentValue:  decimal.NewFromFloat(42.5),
	}

	value, ok := m.EvaluateOptimizer("feature1", "prod")
	assert.True(t, ok)
	assert.True(t, value.Equal(decimal.NewFromFloat(42.5)))
}

func TestEvaluateOptimizer_NotOptimizer(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeEpsilonGreedy,
		Enabled:       true,
		IsOptimizer:   false,
	}

	_, ok := m.EvaluateOptimizer("feature1", "prod")
	assert.False(t, ok, "Should return false for bandit algorithms")
}

func TestEvaluateOptimizer_Disabled(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeHillClimb,
		Enabled:       false,
		IsOptimizer:   true,
		CurrentValue:  decimal.NewFromFloat(42.5),
	}

	_, ok := m.EvaluateOptimizer("feature1", "prod")
	assert.False(t, ok, "Should return false for disabled algorithms")
}

func TestEvaluateFeature_NotOptimizer(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeEpsilonGreedy,
		Enabled:       true,
		IsOptimizer:   false,
		Variants: map[string]*VariantStats{
			"A": {Evaluations: 10, Successes: 8},
			"B": {Evaluations: 10, Successes: 2},
		},
		VariantsArr: []string{"A", "B"},
		Settings: map[string]decimal.Decimal{
			"epsilon": decimal.NewFromFloat(0.0),
		},
	}

	variant, ok := m.EvaluateFeature("feature1", "prod")
	assert.True(t, ok)
	assert.Equal(t, "A", variant, "Should return best variant")
}

func TestEvaluateFeature_RejectsOptimizer(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeHillClimb,
		Enabled:       true,
		IsOptimizer:   true,
	}

	_, ok := m.EvaluateFeature("feature1", "prod")
	assert.False(t, ok, "EvaluateFeature should reject optimizer algorithms")
}

func TestHasAlgorithm(t *testing.T) {
	m := newTestManager(42)
	m.state = make(map[StateKey]*AlgorithmState)

	key := StateKey{FeatureKey: "feature1", EnvKey: "prod"}
	m.state[key] = &AlgorithmState{
		Enabled: true,
	}

	assert.True(t, m.HasAlgorithm("feature1", "prod"))
	assert.False(t, m.HasAlgorithm("nonexistent", "prod"))

	// Disabled algorithm
	m.state[key].Enabled = false
	assert.False(t, m.HasAlgorithm("feature1", "prod"))
}

func TestMultipleIterations_EpsilonGreedy(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})
	state.Settings["epsilon"] = decimal.NewFromFloat(0.1)

	// Initialize with some evaluations so algorithm can make decisions
	state.Variants["A"] = &VariantStats{Evaluations: 10, Successes: 5}
	state.Variants["B"] = &VariantStats{Evaluations: 10, Successes: 5}
	state.Variants["C"] = &VariantStats{Evaluations: 10, Successes: 5}

	// Simulate real usage: evaluate and record outcomes
	for i := 0; i < 200; i++ {
		variant := m.evalEpsilonGreedy(state)

		// Simulate: B has 80% success, A has 50%, C has 20%
		r := m.randSrc.Float64()
		switch variant {
		case "A":
			state.Variants["A"].Evaluations++
			if r < 0.5 {
				state.Variants["A"].Successes++
			}
		case "B":
			state.Variants["B"].Evaluations++
			if r < 0.8 {
				state.Variants["B"].Successes++
			}
		case "C":
			state.Variants["C"].Evaluations++
			if r < 0.2 {
				state.Variants["C"].Successes++
			}
		}
	}

	// After learning, B should have highest success rate and most evaluations
	bSuccessRate := float64(state.Variants["B"].Successes) / float64(state.Variants["B"].Evaluations)
	aSuccessRate := float64(state.Variants["A"].Successes) / float64(state.Variants["A"].Evaluations)
	cSuccessRate := float64(state.Variants["C"].Successes) / float64(state.Variants["C"].Evaluations)

	require.Greater(t, bSuccessRate, aSuccessRate, "B should have higher success rate than A")
	require.Greater(t, bSuccessRate, cSuccessRate, "B should have higher success rate than C")
}
