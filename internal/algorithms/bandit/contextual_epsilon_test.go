package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/togglr-project/togglr/internal/domain"
)

func TestEvalContextualEpsilon_ChoosesVariant(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B", "C"}
	state := &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeContextualEpsilon,
		Enabled:       true,
		IsContextual:  true,
		VariantsArr:   variants,
		Settings: map[string]decimal.Decimal{
			"epsilon":     decimal.NewFromFloat(0.1),
			"feature_dim": decimal.NewFromFloat(8),
		},
		ContextualState: NewContextualAlgorithmState(8, variants, map[string]decimal.Decimal{
			"epsilon": decimal.NewFromFloat(0.1),
		}),
	}

	ctx := map[string]any{
		"user.id":  "789",
		"platform": "ios",
		"age":      25,
	}

	result := m.evalContextualEpsilon(state, ctx)
	assert.Contains(t, variants, result, "Should return one of the variants")
}

func TestEvalContextualEpsilon_Exploration(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B"}
	state := &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeContextualEpsilon,
		Enabled:       true,
		IsContextual:  true,
		VariantsArr:   variants,
		Settings: map[string]decimal.Decimal{
			"epsilon":     decimal.NewFromFloat(1.0), // Always explore
			"feature_dim": decimal.NewFromFloat(8),
		},
		ContextualState: NewContextualAlgorithmState(8, variants, map[string]decimal.Decimal{
			"epsilon": decimal.NewFromFloat(1.0),
		}),
	}

	ctx := map[string]any{"user.id": "test"}

	// With epsilon=1, should explore and choose both variants
	choices := make(map[string]int)
	for i := 0; i < 100; i++ {
		variant := m.evalContextualEpsilon(state, ctx)
		choices[variant]++
	}

	assert.Greater(t, choices["A"], 0, "A should be chosen sometimes during exploration")
	assert.Greater(t, choices["B"], 0, "B should be chosen sometimes during exploration")
}

func TestEvalContextualEpsilon_Exploitation(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B"}
	state := &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeContextualEpsilon,
		Enabled:       true,
		IsContextual:  true,
		VariantsArr:   variants,
		Settings: map[string]decimal.Decimal{
			"epsilon":     decimal.NewFromFloat(0.0), // Always exploit
			"feature_dim": decimal.NewFromFloat(8),
		},
		ContextualState: NewContextualAlgorithmState(8, variants, map[string]decimal.Decimal{
			"epsilon": decimal.NewFromFloat(0.0),
		}),
	}

	ctx := map[string]any{
		"country_code": "JP",
		"device_type":  "mobile",
	}

	// Train: B is much better for this context
	for i := 0; i < 50; i++ {
		m.updateContextualEpsilon(state, "A", 0.1, ctx)
		m.updateContextualEpsilon(state, "B", 0.9, ctx)
	}

	// With epsilon=0, should always exploit (choose B)
	bCount := 0
	for i := 0; i < 100; i++ {
		if m.evalContextualEpsilon(state, ctx) == "B" {
			bCount++
		}
	}

	assert.Greater(t, bCount, 90, "Should mostly choose better variant B with epsilon=0")
}

func TestUpdateContextualEpsilon_UpdatesStats(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B"}
	state := &AlgorithmState{
		AlgorithmType:   domain.AlgorithmTypeContextualEpsilon,
		Enabled:         true,
		IsContextual:    true,
		VariantsArr:     variants,
		Settings:        map[string]decimal.Decimal{"feature_dim": decimal.NewFromFloat(8)},
		ContextualState: NewContextualAlgorithmState(8, variants, nil),
	}

	ctx := map[string]any{"user.id": "test"}

	initialTotalRew := state.ContextualState.Variants["B"].TotalRew

	m.updateContextualEpsilon(state, "B", 0.5, ctx)

	assert.Greater(t, state.ContextualState.Variants["B"].TotalRew, initialTotalRew,
		"TotalRew should increase after update")
}
