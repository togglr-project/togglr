package bandit

import (
	"github.com/shopspring/decimal"
)

func (m *BanditManager) evalPID(state *AlgorithmState, measured, target decimal.Decimal) decimal.Decimal {
	kp := getSettingAsDecimal(state.Settings, "kp", 0.2)
	ki := getSettingAsDecimal(state.Settings, "ki", 0.05)
	kd := getSettingAsDecimal(state.Settings, "kd", 0.01)

	errorValue := target.Sub(measured)
	state.Integral = state.Integral.Add(errorValue)
	derivative := errorValue.Sub(state.LastError)
	state.LastError = errorValue

	output := kp.Mul(errorValue).
		Add(ki.Mul(state.Integral)).
		Add(kd.Mul(derivative))

	state.CurrentValue = state.CurrentValue.Add(output)
	state.MetricSum = state.MetricSum.Add(measured)
	state.Iteration++

	return state.CurrentValue
}
