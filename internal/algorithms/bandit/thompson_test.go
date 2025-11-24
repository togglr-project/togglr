package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEvalThompson_PrefersBetterVariant(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})

	// B has much better success rate
	state.Variants["A"] = &VariantStats{Successes: 10, Failures: 90}
	state.Variants["B"] = &VariantStats{Successes: 90, Failures: 10}

	// Run multiple times
	bCount := 0
	for i := 0; i < 100; i++ {
		if m.evalThompson(state) == "B" {
			bCount++
		}
	}

	// B should be chosen most of the time
	assert.Greater(t, bCount, 80, "Variant B should be chosen most of the time")
}

func TestEvalThompson_CustomPriors(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Settings["prior_alpha"] = decimal.NewFromFloat(2.0)
	state.Settings["prior_beta"] = decimal.NewFromFloat(2.0)

	state.Variants["A"] = &VariantStats{Successes: 5, Failures: 5}
	state.Variants["B"] = &VariantStats{Successes: 5, Failures: 5}

	// With equal stats and symmetric priors, both should be chosen roughly equally
	aCount := 0
	for i := 0; i < 100; i++ {
		if m.evalThompson(state) == "A" {
			aCount++
		}
	}

	// Expect roughly 50/50 split
	assert.Greater(t, aCount, 30, "Variant A should be chosen sometimes")
	assert.Less(t, aCount, 70, "Variant B should also be chosen sometimes")
}

func TestEvalThompson_ColdStart(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})

	// All variants have 0 successes/failures
	result := m.evalThompson(state)
	assert.Contains(t, state.VariantsArr, result, "Should return one of the variants")
}

func TestEvalThompson_IncrementsEvaluations(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Variants["A"] = &VariantStats{Successes: 50, Failures: 50}
	state.Variants["B"] = &VariantStats{Successes: 50, Failures: 50}

	initialA := state.Variants["A"].Evaluations
	initialB := state.Variants["B"].Evaluations

	result := m.evalThompson(state)

	// Only the chosen variant should have incremented evaluations
	if result == "A" {
		assert.Equal(t, initialA+1, state.Variants["A"].Evaluations)
		assert.Equal(t, initialB, state.Variants["B"].Evaluations)
	} else {
		assert.Equal(t, initialA, state.Variants["A"].Evaluations)
		assert.Equal(t, initialB+1, state.Variants["B"].Evaluations)
	}
}
