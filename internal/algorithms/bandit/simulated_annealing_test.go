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
		BestReward:   decimal.NewFromFloat(0.5),
		Temperature:  decimal.NewFromFloat(1.0),
		StepSize:     decimal.NewFromFloat(0.1),
		Settings: map[string]decimal.Decimal{
			"cooling": decimal.NewFromFloat(0.95),
		},
	}

	reward := decimal.NewFromFloat(0.8)
	m.evalSimulatedAnnealing(state, reward)

	assert.True(t, state.BestReward.Equal(reward), "BestReward should be updated on better reward")
}

func TestEvalSimulatedAnnealing_CoolsDown(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		BestReward:   decimal.NewFromFloat(0.5),
		Temperature:  decimal.NewFromFloat(1.0),
		StepSize:     decimal.NewFromFloat(0.1),
		Settings: map[string]decimal.Decimal{
			"cooling": decimal.NewFromFloat(0.9),
		},
	}

	reward := decimal.NewFromFloat(0.8)
	m.evalSimulatedAnnealing(state, reward)

	assert.True(t, state.Temperature.LessThan(decimal.NewFromFloat(1.0)),
		"Temperature should decrease after iteration")
	assert.True(t, state.Temperature.Equal(decimal.NewFromFloat(0.9)),
		"Temperature should be multiplied by cooling factor")
}

func TestEvalSimulatedAnnealing_MayAcceptWorseAtHighTemp(t *testing.T) {
	m := newTestManager(42)

	acceptedWorse := 0
	for i := 0; i < 100; i++ {
		state := &AlgorithmState{
			CurrentValue: decimal.NewFromFloat(1.0),
			BestReward:   decimal.NewFromFloat(0.9),
			Temperature:  decimal.NewFromFloat(10.0),
			StepSize:     decimal.NewFromFloat(0.1),
			Settings: map[string]decimal.Decimal{
				"cooling": decimal.NewFromFloat(0.95),
			},
		}

		initialValue := state.CurrentValue
		reward := decimal.NewFromFloat(0.8)
		m.evalSimulatedAnnealing(state, reward)

		if !state.CurrentValue.Equal(initialValue) {
			acceptedWorse++
		}
	}

	assert.Greater(t, acceptedWorse, 0, "Should sometimes accept worse solutions at high temperature")
}

func TestEvalSimulatedAnnealing_IterationIncrement(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(1.0),
		Temperature:  decimal.NewFromFloat(1.0),
		StepSize:     decimal.NewFromFloat(0.1),
		Settings: map[string]decimal.Decimal{
			"cooling": decimal.NewFromFloat(0.95),
		},
	}

	reward := decimal.NewFromFloat(0.5)
	m.evalSimulatedAnnealing(state, reward)

	assert.Equal(t, uint64(1), state.Iteration, "Iteration should increment")
}
