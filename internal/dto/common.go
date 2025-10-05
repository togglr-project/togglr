package dto

import (
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// Common helper functions for optional type conversions

// String conversion helpers.
func ptrToOptNilString(p *string) generatedapi.OptNilString {
	if p == nil {
		var o generatedapi.OptNilString

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilString(*p)
}

// Rule attribute conversion helpers.
func ruleAttribute2OptString(attr domain.RuleAttribute) generatedapi.OptString {
	if attr == "" {
		return generatedapi.OptString{}
	}

	return generatedapi.NewOptString(attr.String())
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
