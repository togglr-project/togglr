package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

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
