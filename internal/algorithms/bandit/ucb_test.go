package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEvalUCB_ColdStart(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B", "C"})

	// UCB should prefer unexplored variants first
	seen := make(map[string]bool)
	for i := 0; i < 3; i++ {
		result := m.evalUCB(state)
		seen[result] = true
	}

	// All variants should be explored once
	assert.Len(t, seen, 3, "All variants should be explored during cold start")
}

func TestEvalUCB_ExploitsBestAfterExploration(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})

	// Both variants have been explored, B is much better
	state.Variants["A"] = &VariantStats{Evaluations: 50, Successes: 10}
	state.Variants["B"] = &VariantStats{Evaluations: 50, Successes: 40}

	// Run multiple evaluations
	bCount := 0
	for i := 0; i < 100; i++ {
		if m.evalUCB(state) == "B" {
			bCount++
		}
	}

	// B should be chosen most of the time
	assert.Greater(t, bCount, 70, "Better variant should be chosen more often")
}

func TestEvalUCB_CustomConfidence(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Settings["confidence"] = decimal.NewFromFloat(0.1) // Very low confidence = mostly exploitation

	state.Variants["A"] = &VariantStats{Evaluations: 100, Successes: 60}
	state.Variants["B"] = &VariantStats{Evaluations: 100, Successes: 40}

	// With very low confidence, should exploit more aggressively
	aCount := 0
	for i := 0; i < 100; i++ {
		if m.evalUCB(state) == "A" {
			aCount++
		}
	}

	// A has better mean, should be chosen more with low confidence
	assert.Greater(t, aCount, 50, "Better variant should be exploited more with low confidence")
}

func TestEvalUCB_IncrementsEvaluations(t *testing.T) {
	m := newTestManager(42)
	state := newTestState([]string{"A", "B"})
	state.Variants["A"] = &VariantStats{Evaluations: 10, Successes: 5}
	state.Variants["B"] = &VariantStats{Evaluations: 10, Successes: 5}

	initialTotal := state.Variants["A"].Evaluations + state.Variants["B"].Evaluations

	m.evalUCB(state)

	newTotal := state.Variants["A"].Evaluations + state.Variants["B"].Evaluations
	assert.Equal(t, initialTotal+1, newTotal, "Total evaluations should increase by 1")
}
