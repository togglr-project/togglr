package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/togglr-project/togglr/internal/domain"
)

func TestEvalContextualThompson_ChoosesVariant(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B", "C"}
	state := &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeContextualThompson,
		Enabled:       true,
		IsContextual:  true,
		VariantsArr:   variants,
		Settings: map[string]decimal.Decimal{
			"prior_variance": decimal.NewFromFloat(1.0),
			"feature_dim":    decimal.NewFromFloat(8),
		},
		ContextualState: NewContextualAlgorithmState(8, variants, map[string]decimal.Decimal{
			"prior_variance": decimal.NewFromFloat(1.0),
		}),
	}

	ctx := map[string]any{
		"user.id":      "456",
		"country_code": "DE",
		"browser":      "Chrome",
	}

	result := m.evalContextualThompson(state, ctx)
	assert.Contains(t, variants, result, "Should return one of the variants")
}

func TestEvalContextualThompson_LearnsFromRewards(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B"}
	state := &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeContextualThompson,
		Enabled:       true,
		IsContextual:  true,
		VariantsArr:   variants,
		Settings: map[string]decimal.Decimal{
			"prior_variance": decimal.NewFromFloat(0.1), // Lower variance = less exploration
			"feature_dim":    decimal.NewFromFloat(8),
		},
		ContextualState: NewContextualAlgorithmState(8, variants, map[string]decimal.Decimal{
			"prior_variance": decimal.NewFromFloat(0.1),
		}),
	}

	// Context for European desktop users
	ctx := map[string]any{
		"country_code": "DE",
		"device_type":  "desktop",
	}

	// Train: variant A is better for this context
	for i := 0; i < 100; i++ {
		m.updateContextualThompson(state, "A", 0.9, ctx)
		m.updateContextualThompson(state, "B", 0.1, ctx)
	}

	// After training, A should be chosen more often
	aCount := 0
	for i := 0; i < 100; i++ {
		if m.evalContextualThompson(state, ctx) == "A" {
			aCount++
		}
	}

	assert.Greater(t, aCount, 50, "Better variant A should be chosen more often after training")
}

func TestUpdateContextualThompson_UpdatesStats(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B"}
	state := &AlgorithmState{
		AlgorithmType:   domain.AlgorithmTypeContextualThompson,
		Enabled:         true,
		IsContextual:    true,
		VariantsArr:     variants,
		Settings:        map[string]decimal.Decimal{"feature_dim": decimal.NewFromFloat(8)},
		ContextualState: NewContextualAlgorithmState(8, variants, nil),
	}

	ctx := map[string]any{"user.id": "test"}

	initialSuccesses := state.ContextualState.Variants["A"].Successes

	m.updateContextualThompson(state, "A", 1.0, ctx)

	assert.Equal(t, initialSuccesses+1, state.ContextualState.Variants["A"].Successes,
		"Successes should increment on positive reward")
}
