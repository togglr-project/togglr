package bandit

import (
	"github.com/shopspring/decimal"
)

func (m *BanditManager) evalPID(state *AlgorithmState, measured, target decimal.Decimal) decimal.Decimal {
	kp := getSettingAsDecimal(state.Settings, "kp", 0.2)
	ki := getSettingAsDecimal(state.Settings, "ki", 0.05)
	kd := getSettingAsDecimal(state.Settings, "kd", 0.01)

	errorValue := target.Sub(measured)
	state.Settings["integral"] = state.Settings["integral"].Add(errorValue)
	derivative := errorValue.Sub(state.Settings["prev_error"])
	state.Settings["prev_error"] = errorValue

	output := kp.Mul(errorValue).
		Add(ki.Mul(state.Settings["integral"])).
		Add(kd.Mul(derivative))

	state.CurrentValue = state.CurrentValue.Add(output)

	return state.CurrentValue
}
