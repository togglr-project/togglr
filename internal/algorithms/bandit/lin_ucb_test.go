package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/togglr-project/togglr/internal/domain"
)

func TestEvalLinUCB_ChoosesVariant(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B", "C"}
	state := &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeLinUCB,
		Enabled:       true,
		IsContextual:  true,
		VariantsArr:   variants,
		Settings: map[string]decimal.Decimal{
			"alpha":       decimal.NewFromFloat(1.0),
			"feature_dim": decimal.NewFromFloat(8),
		},
		ContextualState: NewContextualAlgorithmState(8, variants, map[string]decimal.Decimal{
			"alpha": decimal.NewFromFloat(1.0),
		}),
	}

	ctx := map[string]any{
		"user.id":      "123",
		"country_code": "US",
		"device_type":  "mobile",
	}

	result := m.evalLinUCB(state, ctx)
	assert.Contains(t, variants, result, "Should return one of the variants")
}

func TestEvalLinUCB_LearnsFromRewards(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B"}
	state := &AlgorithmState{
		AlgorithmType: domain.AlgorithmTypeLinUCB,
		Enabled:       true,
		IsContextual:  true,
		VariantsArr:   variants,
		Settings: map[string]decimal.Decimal{
			"alpha":       decimal.NewFromFloat(0.1), // Low exploration
			"feature_dim": decimal.NewFromFloat(8),
		},
		ContextualState: NewContextualAlgorithmState(8, variants, map[string]decimal.Decimal{
			"alpha": decimal.NewFromFloat(0.1),
		}),
	}

	// Context for US mobile users
	ctx := map[string]any{
		"country_code": "US",
		"device_type":  "mobile",
	}

	// Train: variant B is better for this context
	for i := 0; i < 50; i++ {
		m.updateLinUCB(state, "A", 0.2, ctx)
		m.updateLinUCB(state, "B", 0.8, ctx)
	}

	// After training, B should be chosen more often for similar context
	bCount := 0
	for i := 0; i < 100; i++ {
		if m.evalLinUCB(state, ctx) == "B" {
			bCount++
		}
	}

	assert.Greater(t, bCount, 60, "Better variant B should be chosen more often after training")
}

func TestUpdateLinUCB_UpdatesModel(t *testing.T) {
	m := newTestManager(42)

	variants := []string{"A", "B"}
	state := &AlgorithmState{
		AlgorithmType:   domain.AlgorithmTypeLinUCB,
		Enabled:         true,
		IsContextual:    true,
		VariantsArr:     variants,
		Settings:        map[string]decimal.Decimal{"feature_dim": decimal.NewFromFloat(8)},
		ContextualState: NewContextualAlgorithmState(8, variants, nil),
	}

	ctx := map[string]any{"user.id": "test"}

	initialPulls := state.ContextualState.Variants["A"].Pulls
	initialTotalRew := state.ContextualState.Variants["A"].TotalRew

	m.updateLinUCB(state, "A", 1.0, ctx)

	// Check that stats were updated
	// Note: Pulls is updated in evalLinUCB, not updateLinUCB
	assert.Greater(t, state.ContextualState.Variants["A"].TotalRew, initialTotalRew,
		"TotalRew should increase")
	assert.Equal(t, initialPulls, state.ContextualState.Variants["A"].Pulls,
		"Pulls not changed in update")
}
