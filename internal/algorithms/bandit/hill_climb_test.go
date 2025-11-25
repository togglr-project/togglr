package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEvalHillClimb_ImproveOnBetterReward(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		BestReward:   decimal.NewFromFloat(0.5),
		StepSize:     decimal.NewFromFloat(0.1),
		Settings:     make(map[string]decimal.Decimal),
	}

	reward := decimal.NewFromFloat(0.8)
	result := m.evalHillClimb(state, reward)

	assert.True(t, result.GreaterThan(decimal.NewFromFloat(1.0)), "Value should increase with better reward")
	assert.Equal(t, uint64(1), state.Iteration, "Iteration should increment")
	assert.True(t, state.BestReward.Equal(reward), "BestReward should be updated")
}

func TestEvalHillClimb_ReverseOnWorseReward(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		BestReward:   decimal.NewFromFloat(0.8),
		StepSize:     decimal.NewFromFloat(0.1),
		Settings:     make(map[string]decimal.Decimal),
	}

	reward := decimal.NewFromFloat(0.5)
	m.evalHillClimb(state, reward)

	assert.True(t, state.StepSize.LessThan(decimal.Zero),
		"Step should be negative after worse reward")
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

	assert.True(t, result.GreaterThan(decimal.Zero), "Should move in positive direction by default")
}
