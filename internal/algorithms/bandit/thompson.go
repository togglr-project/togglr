package bandit

import (
	"github.com/shopspring/decimal"
	"gonum.org/v1/gonum/stat/distuv"
)

func (m *BanditManager) evalThompson(state *AlgorithmState) string {
	state.mu.Lock()
	defer state.mu.Unlock()

	// priors
	priorAlpha := decimal.NewFromFloat(1.0)
	priorBeta := decimal.NewFromFloat(1.0)
	if value, ok := state.Settings["prior_alpha"]; ok {
		priorAlpha = value
	}
	if value, ok := state.Settings["prior_beta"]; ok {
		priorBeta = value
	}

	bestKey := ""
	bestSample := -1.0
	for key, value := range state.Variants {
		alpha := priorAlpha.Add(decimal.NewFromUint64(value.Successes))
		beta := priorBeta.Add(decimal.NewFromUint64(value.Failures))
		// Beta sampling using gonum distuv.Beta
		distuvBeta := distuv.Beta{Alpha: alpha.InexactFloat64(), Beta: beta.InexactFloat64(), Src: m.randSrc}
		sample := distuvBeta.Rand()
		if sample > bestSample {
			bestSample = sample
			bestKey = key
		}
	}
	if bestKey == "" {
		// fallback
		return state.VariantsArr[m.randSrc.Intn(len(state.VariantsArr))]
	}
	state.Variants[bestKey].Evaluations++

	return bestKey
}
