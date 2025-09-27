package dto

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// Common helper functions for optional type conversions

// Time conversion helpers.
func ptrToOptNilDateTime(p *time.Time) generatedapi.OptNilDateTime {
	if p == nil {
		var o generatedapi.OptNilDateTime

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilDateTime(*p)
}

func optNilDateTimeToPtr(v generatedapi.OptNilDateTime) *time.Time {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	t, _ := v.Get()

	return &t
}

func ptrToOptNilDuration(p *time.Duration) generatedapi.OptNilDuration {
	if p == nil {
		var o generatedapi.OptNilDuration

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilDuration(*p)
}

func optNilDurationToPtr(v generatedapi.OptNilDuration) *time.Duration {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	t, _ := v.Get()

	return &t
}

// String conversion helpers.
func ptrToOptNilString(p *string) generatedapi.OptNilString {
	if p == nil {
		var o generatedapi.OptNilString

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilString(*p)
}

func optNilStringToPtr(v generatedapi.OptNilString) *string {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	s, _ := v.Get()

	return &s
}

func optNilStringToString(v generatedapi.OptNilString) string {
	if !v.IsSet() || v.IsNull() {
		return ""
	}

	s, _ := v.Get()

	return s
}

func optStringToString(v generatedapi.OptString) string {
	if !v.IsSet() {
		return ""
	}

	return v.Value
}

func optStringToPtr(v generatedapi.OptString) *string {
	if !v.IsSet() {
		return nil
	}

	return &v.Value
}

// Boolean conversion helpers.
func ptrToOptNilBool(p *bool) generatedapi.OptNilBool {
	if p == nil {
		var o generatedapi.OptNilBool

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilBool(*p)
}

func optNilBoolToPtr(v generatedapi.OptNilBool) *bool {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	b, _ := v.Get()

	return &b
}

// Rule attribute conversion helpers.
func ruleAttribute2OptString(attr domain.RuleAttribute) generatedapi.OptString {
	if attr == "" {
		return generatedapi.OptString{}
	}

	return generatedapi.NewOptString(attr.String())
}

func optStringToRuleAttribute(opt generatedapi.OptString) domain.RuleAttribute {
	if !opt.IsSet() {
		return ""
	}

	return domain.RuleAttribute(opt.Value)
}

// Flag variant ID conversion helpers.
func flagVariantRef2OptString(flagVariantID *domain.FlagVariantID) generatedapi.OptString {
	if flagVariantID == nil {
		return generatedapi.OptString{}
	}

	return generatedapi.NewOptString(flagVariantID.String())
}

func optString2FlagVariantIDRef(optString generatedapi.OptString) *domain.FlagVariantID {
	if !optString.IsSet() {
		return nil
	}

	id := domain.FlagVariantID(optString.Value)

	return &id
}

// UUID conversion helpers.
func optString2FlagVariantIDRefFromUUID(optUUID generatedapi.OptUUID) *domain.FlagVariantID {
	if !optUUID.IsSet() {
		return nil
	}

	id := domain.FlagVariantID(optUUID.Value.String())

	return &id
}

func flagVariantIDRef2OptUUID(flagVariantID *domain.FlagVariantID) generatedapi.OptUUID {
	if flagVariantID == nil {
		return generatedapi.OptUUID{}
	}
	// Note: This would need uuid.Parse to convert string to uuid.UUID
	// For now, returning empty OptUUID
	return generatedapi.OptUUID{}
}
