package bandit

import (
	"math"

	"github.com/shopspring/decimal"
)

func (m *BanditManager) evalSimulatedAnnealing(state *AlgorithmState, reward decimal.Decimal) decimal.Decimal {
	if state.Temperature.IsZero() {
		state.Temperature = getSettingAsDecimal(state.Settings, "temp", 1.0)
	}
	if state.StepSize.IsZero() {
		state.StepSize = getSettingAsDecimal(state.Settings, "step_scale", 0.1)
	}

	cooling := getSettingAsDecimal(state.Settings, "cooling", 0.95)
	candidate := state.CurrentValue.Add(decimal.NewFromFloat((m.randSrc.Float64()*2 - 1) * state.StepSize.InexactFloat64()))

	delta := reward.Sub(state.BestReward)
	if delta.GreaterThan(decimal.Zero) ||
		m.randSrc.Float64() < math.Exp(delta.InexactFloat64()/state.Temperature.InexactFloat64()) {
		state.CurrentValue = candidate
		if reward.GreaterThan(state.BestReward) {
			state.BestReward = reward
			state.BestValue = candidate
		}
	}

	state.Temperature = state.Temperature.Mul(cooling)
	state.MetricSum = state.MetricSum.Add(reward)
	state.Iteration++

	return state.CurrentValue
}
