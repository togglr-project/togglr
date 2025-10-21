package bandit

import (
	"github.com/shopspring/decimal"
)

func (m *BanditManager) evalHillClimb(state *AlgorithmState, lastReward decimal.Decimal) decimal.Decimal {
	step := getSettingAsFloat64(state.Settings, "step", 0.05)
	direction := getSettingAsFloat64(state.Settings, "direction", 1)
	prev := state.CurrentValue
	candidate := prev.Add(decimal.NewFromFloat(step).Mul(decimal.NewFromFloat(direction)))

	if lastReward.GreaterThan(state.MetricSum) {
		state.CurrentValue = candidate
	} else {
		state.Settings["direction"] = state.Settings["direction"].Mul(decimal.NewFromFloat(-1))
	}

	state.Iteration++

	return state.CurrentValue
}
