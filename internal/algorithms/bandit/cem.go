package bandit

import (
	"sort"

	"github.com/shopspring/decimal"
)

//nolint:unused // will be implemented later
func (m *BanditManager) evalCEM(
	state *AlgorithmState,
	samples, rewards []decimal.Decimal,
) decimal.Decimal {
	eliteFrac := getSettingAsFloat64(state.Settings, "elite_fraction", 0.2)
	popSize := int(getSettingAsFloat64(state.Settings, "population_size", 20))

	type sampleReward struct{ val, rew float64 }
	sr := make([]sampleReward, len(samples))
	for i := range samples {
		sr[i] = sampleReward{samples[i].InexactFloat64(), rewards[i].InexactFloat64()}
	}
	sort.Slice(sr, func(i, j int) bool { return sr[i].rew > sr[j].rew })
	elites := sr[:int(float64(popSize)*eliteFrac)]
	mean := 0.0
	for _, e := range elites {
		mean += e.val
	}
	mean /= float64(len(elites))
	state.CurrentValue = decimal.NewFromFloat(mean)

	return state.CurrentValue
}
