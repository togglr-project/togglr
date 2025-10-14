package bandit

import (
	"github.com/shopspring/decimal"
)

func (m *BanditManager) evalEpsilonGreedy(state *AlgorithmState) (string, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	epsilon := decimal.NewFromFloat(0.1)
	if value, ok := state.Settings["epsilon"]; ok {
		epsilon = value
	}

	// exploration?
	if decimal.NewFromFloat(m.randSrc.Float64()).LessThan(epsilon) {
		// choose uniformly random variant
		return state.VariantsArr[m.randSrc.Intn(len(state.VariantsArr))], nil
	}

	// exploitation: choose the best success rate (successes / impressions) or metric_avg if you prefer
	var bestKey string
	bestScore := -1.0
	for key, value := range state.Variants {
		var score float64
		if value.Evaluations > 0 {
			score = float64(value.Successes) / float64(value.Evaluations)
		} else {
			score = 0.0
		}
		if score > bestScore {
			bestScore = score
			bestKey = key
		}
	}

	// If all zero (cold start), pick random
	if bestKey == "" {
		return state.VariantsArr[m.randSrc.Intn(len(state.VariantsArr))], nil
	}

	// increment impression (evaluation) immediately for real-time accounting
	state.Variants[bestKey].Evaluations++

	return bestKey, nil
}
