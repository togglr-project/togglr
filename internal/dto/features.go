package dto

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainFeatureToAPI converts domain Feature to generated API Feature
func DomainFeatureToAPI(feature domain.Feature) generatedapi.Feature {
	return generatedapi.Feature{
		ID:             feature.ID.String(),
		ProjectID:      feature.ProjectID.String(),
		Key:            feature.Key,
		Name:           feature.Name,
		Description:    ptrToOptNilString(&feature.Description),
		Kind:           generatedapi.FeatureKind(feature.Kind),
		DefaultVariant: feature.DefaultVariant,
		Enabled:        feature.Enabled,
		RolloutKey:     ruleAttribute2OptString(feature.RolloutKey),
		CreatedAt:      feature.CreatedAt,
		UpdatedAt:      feature.UpdatedAt,
	}
}

// DomainFeatureExtendedToAPI converts domain FeatureExtended to generated API FeatureExtended
func DomainFeatureExtendedToAPI(feature domain.FeatureExtended, isActive bool, nextState *bool, nextStateTime *time.Time) generatedapi.FeatureExtended {
	item := generatedapi.FeatureExtended{
		ID:             feature.ID.String(),
		ProjectID:      feature.ProjectID.String(),
		Key:            feature.Key,
		Name:           feature.Name,
		Description:    ptrToOptNilString(&feature.Description),
		Kind:           generatedapi.FeatureKind(feature.Kind),
		DefaultVariant: feature.DefaultVariant,
		Enabled:        feature.Enabled,
		RolloutKey:     ruleAttribute2OptString(feature.RolloutKey),
		CreatedAt:      feature.CreatedAt,
		UpdatedAt:      feature.UpdatedAt,
		IsActive:       isActive,
	}

	// Handle next state
	if nextState != nil && nextStateTime != nil && !nextStateTime.IsZero() {
		item.NextState = generatedapi.NewOptNilBool(*nextState)
		item.NextStateTime = generatedapi.NewOptNilDateTime(*nextStateTime)
	}

	// Convert tags
	item.Tags = DomainTagsToAPI(feature.Tags)

	return item
}

// DomainFeaturesToAPI converts slice of domain Features to slice of generated API Features
func DomainFeaturesToAPI(features []domain.Feature) []generatedapi.Feature {
	resp := make([]generatedapi.Feature, 0, len(features))
	for _, feature := range features {
		resp = append(resp, DomainFeatureToAPI(feature))
	}
	return resp
}

// APIFeatureToDomain converts generated API Feature to domain Feature
func APIFeatureToDomain(feature generatedapi.Feature) domain.Feature {
	return domain.Feature{
		ID:             domain.FeatureID(feature.ID),
		ProjectID:      domain.ProjectID(feature.ProjectID),
		Key:            feature.Key,
		Name:           feature.Name,
		Description:    *optNilStringToPtr(feature.Description),
		Kind:           domain.FeatureKind(feature.Kind),
		DefaultVariant: feature.DefaultVariant,
		Enabled:        feature.Enabled,
		RolloutKey:     domain.RuleAttribute(optStringToRuleAttribute(feature.RolloutKey)),
		CreatedAt:      feature.CreatedAt,
		UpdatedAt:      feature.UpdatedAt,
	}
}
