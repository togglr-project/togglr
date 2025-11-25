package bandit

import (
	"github.com/shopspring/decimal"
)

func (m *BanditManager) evalHillClimb(state *AlgorithmState, lastReward decimal.Decimal) decimal.Decimal {
	if state.StepSize.IsZero() {
		state.StepSize = getSettingAsDecimal(state.Settings, "step", 0.05)
	}

	step := state.StepSize
	candidate := state.CurrentValue.Add(step)

	if lastReward.GreaterThan(state.BestReward) {
		state.BestReward = lastReward
		state.BestValue = state.CurrentValue
		state.CurrentValue = candidate
	} else {
		state.StepSize = step.Neg()
	}

	state.MetricSum = state.MetricSum.Add(lastReward)
	state.Iteration++

	return state.CurrentValue
}
