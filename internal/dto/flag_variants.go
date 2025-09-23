package dto

import (
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainFlagVariantToAPI converts domain FlagVariant to generated API FlagVariant
func DomainFlagVariantToAPI(variant domain.FlagVariant) generatedapi.FlagVariant {
	return generatedapi.FlagVariant{
		ID:             variant.ID.String(),
		FeatureID:      variant.FeatureID.String(),
		Name:           variant.Name,
		RolloutPercent: int(variant.RolloutPercent),
	}
}

// DomainFlagVariantsToAPI converts slice of domain FlagVariants to slice of generated API FlagVariants
func DomainFlagVariantsToAPI(variants []domain.FlagVariant) []generatedapi.FlagVariant {
	resp := make([]generatedapi.FlagVariant, 0, len(variants))
	for _, variant := range variants {
		resp = append(resp, DomainFlagVariantToAPI(variant))
	}
	return resp
}

// APIFlagVariantToDomain converts generated API FlagVariant to domain FlagVariant
func APIFlagVariantToDomain(variant generatedapi.FlagVariant) domain.FlagVariant {
	return domain.FlagVariant{
		ID:             domain.FlagVariantID(variant.ID),
		FeatureID:      domain.FeatureID(variant.FeatureID),
		Name:           variant.Name,
		RolloutPercent: uint8(variant.RolloutPercent),
	}
}
