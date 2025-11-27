package bandit

import (
	"github.com/shopspring/decimal"
)

//nolint:unused // will be implemented later
func (m *BanditManager) evalBayesOpt(
	state *AlgorithmState,
	samples, rewards []decimal.Decimal,
) decimal.Decimal {
	best := samples[0]
	bestScore := rewards[0].InexactFloat64()
	noise := getSettingAsFloat64(state.Settings, "noise", 0.01)
	for i := 1; i < len(samples); i++ {
		s := rewards[i].InexactFloat64() + m.randSrc.Float64()*noise
		if s > bestScore {
			best = samples[i]
			bestScore = s
		}
	}
	state.CurrentValue = best

	return state.CurrentValue
}
