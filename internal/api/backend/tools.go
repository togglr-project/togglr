package apibackend

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func optString2FlagVariantIDRef(optString generatedapi.OptUUID) *domain.FlagVariantID {
	if !optString.IsSet() {
		return nil
	}

	id := domain.FlagVariantID(optString.Value.String())

	return &id
}

func flagVariantRef2OptString(flagVariantID *domain.FlagVariantID) generatedapi.OptString {
	if flagVariantID == nil {
		return generatedapi.OptString{}
	}

	return generatedapi.NewOptString(flagVariantID.String())
}

func ruleAttribute2OptString(attr domain.RuleAttribute) generatedapi.OptString {
	if attr == "" {
		return generatedapi.OptString{}
	}

	return generatedapi.NewOptString(attr.String())
}

func optNilDurationToPtr(v generatedapi.OptNilDuration) *time.Duration {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	t, _ := v.Get()

	return &t
}

func optNilDateTimeToPtr(v generatedapi.OptNilDateTime) *time.Time {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	t, _ := v.Get()

	return &t
}

func ptrToOptNilDateTime(p *time.Time) generatedapi.OptNilDateTime {
	if p == nil {
		var o generatedapi.OptNilDateTime
		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilDateTime(*p)
}

func ptrToOptNilDuration(p *time.Duration) generatedapi.OptNilDuration {
	if p == nil {
		var o generatedapi.OptNilDuration
		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilDuration(*p)
}

func optNilStringToPtr(v generatedapi.OptNilString) *string {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	s, _ := v.Get()

	return &s
}

func ptrToOptNilString(p *string) generatedapi.OptNilString {
	if p == nil {
		var o generatedapi.OptNilString
		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilString(*p)
}
