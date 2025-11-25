package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEvalPID_ConvergesToTarget(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		Integral:     decimal.Zero,
		LastError:    decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"kp": decimal.NewFromFloat(0.5),
			"ki": decimal.NewFromFloat(0.1),
			"kd": decimal.NewFromFloat(0.05),
		},
	}

	target := decimal.NewFromFloat(10.0)

	for i := 0; i < 50; i++ {
		measured := state.CurrentValue
		m.evalPID(state, measured, target)
	}

	diff := target.Sub(state.CurrentValue).Abs()
	assert.True(t, diff.LessThan(decimal.NewFromFloat(1.0)),
		"PID should converge towards target value")
}

func TestEvalPID_RespondsToError(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(5.0),
		Integral:     decimal.Zero,
		LastError:    decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"kp": decimal.NewFromFloat(0.2),
			"ki": decimal.NewFromFloat(0.0),
			"kd": decimal.NewFromFloat(0.0),
		},
	}

	target := decimal.NewFromFloat(10.0)
	measured := decimal.NewFromFloat(5.0)

	result := m.evalPID(state, measured, target)

	assert.True(t, result.GreaterThan(decimal.NewFromFloat(5.0)),
		"PID output should increase when below target")
}

func TestEvalPID_IntegralAccumulates(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		Integral:     decimal.Zero,
		LastError:    decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"kp": decimal.NewFromFloat(0.0),
			"ki": decimal.NewFromFloat(1.0),
			"kd": decimal.NewFromFloat(0.0),
		},
	}

	target := decimal.NewFromFloat(1.0)
	measured := decimal.NewFromFloat(0.0)

	for i := 0; i < 5; i++ {
		m.evalPID(state, measured, target)
	}

	assert.True(t, state.Integral.GreaterThan(decimal.NewFromFloat(4.0)),
		"Integral term should accumulate over iterations")
}

func TestEvalPID_DerivativeRespondsToChange(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		Integral:     decimal.Zero,
		LastError:    decimal.NewFromFloat(1.0),
		Settings: map[string]decimal.Decimal{
			"kp": decimal.NewFromFloat(0.0),
			"ki": decimal.NewFromFloat(0.0),
			"kd": decimal.NewFromFloat(1.0),
		},
	}

	target := decimal.NewFromFloat(5.0)
	measured := decimal.NewFromFloat(3.0)

	result := m.evalPID(state, measured, target)

	assert.True(t, result.Equal(decimal.NewFromFloat(1.0)),
		"PID derivative term should respond to error change")
}

func TestEvalPID_IterationIncrement(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.NewFromFloat(0.0),
		Settings: map[string]decimal.Decimal{
			"kp": decimal.NewFromFloat(0.2),
			"ki": decimal.NewFromFloat(0.1),
			"kd": decimal.NewFromFloat(0.05),
		},
	}

	m.evalPID(state, decimal.NewFromFloat(0.5), decimal.NewFromFloat(1.0))

	assert.Equal(t, uint64(1), state.Iteration, "Iteration should increment")
}
