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
