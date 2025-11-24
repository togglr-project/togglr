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
