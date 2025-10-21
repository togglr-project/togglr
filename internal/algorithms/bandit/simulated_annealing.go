package bandit

import (
	"math"

	"github.com/shopspring/decimal"
)

func (m *BanditManager) evalSimulatedAnnealing(state *AlgorithmState, reward decimal.Decimal) decimal.Decimal {
	temp := getSettingAsDecimal(state.Settings, "temp", 1.0)
	cooling := getSettingAsDecimal(state.Settings, "cooling", 0.95)
	step := getSettingAsDecimal(state.Settings, "step_scale", 0.1)
	candidate := state.CurrentValue.Add(decimal.NewFromFloat((m.randSrc.Float64()*2 - 1) * step.InexactFloat64()))

	delta := reward.Sub(state.MetricSum)
	if delta.GreaterThan(decimal.Zero) ||
		m.randSrc.Float64() < math.Exp(delta.InexactFloat64()/temp.InexactFloat64()) {
		state.CurrentValue = candidate
		state.MetricSum = reward
	}
	state.Settings["temp"] = temp.Mul(cooling)

	return state.CurrentValue
}
