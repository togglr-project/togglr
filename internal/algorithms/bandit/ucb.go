package bandit

import (
	"math"
)

func (m *BanditManager) evalUCB(state *AlgorithmState) string {
	state.mu.Lock()
	defer state.mu.Unlock()

	confidence := 2.0
	if value, ok := state.Settings["confidence"]; ok {
		confidence = value.InexactFloat64()
	}

	total := uint64(0)
	for _, stats := range state.Variants {
		total += stats.Evaluations
	}

	// cold start: any variant with 0 impressions is preferred
	for key, stats := range state.Variants {
		if stats.Evaluations == 0 {
			state.Variants[key].Evaluations++

			return key
		}
	}

	bestKey := ""
	bestScore := -1.0
	for key, stats := range state.Variants {
		mean := float64(stats.Successes) / float64(stats.Evaluations)
		score := mean + confidence*math.Sqrt(math.Log(float64(total))/float64(stats.Evaluations))
		if score > bestScore {
			bestScore = score
			bestKey = key
		}
	}

	if bestKey == "" {
		return state.VariantsArr[m.randSrc.Intn(len(state.VariantsArr))]
	}

	state.Variants[bestKey].Evaluations++

	return bestKey
}
