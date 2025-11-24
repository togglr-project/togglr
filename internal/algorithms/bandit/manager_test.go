package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/togglr-project/togglr/internal/domain"
)

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
