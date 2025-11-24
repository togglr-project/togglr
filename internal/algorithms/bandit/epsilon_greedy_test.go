package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEvalEpsilonGreedy_Exploitation(t *testing.T) {
	// With epsilon=0, should always exploit (choose best variant)
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})
	state.Settings["epsilon"] = decimal.NewFromFloat(0.0)

	// Set up stats: B has the best success rate
	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 30}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 70}
	state.Variants["C"] = &VariantStats{Evaluations: 100, Successes: 50}

	result := m.evalEpsilonGreedy(state)
	assert.Equal(t, "B", result, "Should choose variant B with highest success rate")
}

func TestEvalEpsilonGreedy_Exploration(t *testing.T) {
	// With epsilon=1, should always explore (random choice)
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})
	state.Settings["epsilon"] = decimal.NewFromFloat(1.0)

	// Set up stats: B has the best success rate, but should still explore
	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 30}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 70}
	state.Variants["C"] = &VariantStats{Evaluations: 100, Successes: 50}

	// Run multiple times and collect results
	results := make(map[string]int)
	for i := 0; i < 100; i++ {
		result := m.evalEpsilonGreedy(state)
		results[result]++
	}

	// With exploration, all variants should be chosen at least once
	assert.Greater(t, results["A"], 0, "Variant A should be chosen at least once")
	assert.Greater(t, results["B"], 0, "Variant B should be chosen at least once")
	assert.Greater(t, results["C"], 0, "Variant C should be chosen at least once")
}

func TestEvalEpsilonGreedy_ColdStart(t *testing.T) {
	// When all variants have 0 evaluations, should pick random
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})
	state.Settings["epsilon"] = decimal.NewFromFloat(0.0)

	result := m.evalEpsilonGreedy(state)
	assert.Contains(t, state.VariantsArr, result, "Should return one of the variants")
}

func TestEvalEpsilonGreedy_DefaultEpsilon(t *testing.T) {
	// Default epsilon should be 0.1
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})

	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 90}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 10}

	// Run many times - most should be A (exploitation), but some B (exploration)
	aCount := 0
	for i := 0; i < 1000; i++ {
		if m.evalEpsilonGreedy(state) == "A" {
			aCount++
		}
	}

	// With epsilon=0.1, expect ~90% A choices
	assert.Greater(t, aCount, 800, "Most choices should be A (exploitation)")
	assert.Less(t, aCount, 1000, "Some choices should be B (exploration)")
}
