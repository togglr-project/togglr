package dto

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainFeatureToAPI converts domain Feature to generated API Feature.
func DomainFeatureToAPI(feature domain.Feature) generatedapi.Feature {
	return generatedapi.Feature{
		ID:           feature.ID.String(),
		ProjectID:    feature.ProjectID.String(),
		Key:          feature.Key,
		Name:         feature.Name,
		Description:  ptrToOptNilString(&feature.Description),
		Kind:         generatedapi.FeatureKind(feature.Kind),
		RolloutKey:   ruleAttribute2OptString(feature.RolloutKey),
		Enabled:      feature.Enabled,
		DefaultValue: feature.DefaultValue,
		CreatedAt:    feature.CreatedAt,
		UpdatedAt:    feature.UpdatedAt,
	}
}

// DomainFeatureExtendedToAPI converts domain FeatureExtended to generated API FeatureExtended.
func DomainFeatureExtendedToAPI(
	feature domain.FeatureExtended,
	isActive bool,
	nextState *bool,
	nextStateTime *time.Time,
) generatedapi.FeatureExtended {
	item := generatedapi.FeatureExtended{
		ID:           feature.ID.String(),
		ProjectID:    feature.ProjectID.String(),
		Key:          feature.Key,
		Name:         feature.Name,
		Description:  ptrToOptNilString(&feature.Description),
		Kind:         generatedapi.FeatureKind(feature.Kind),
		RolloutKey:   ruleAttribute2OptString(feature.RolloutKey),
		Enabled:      feature.Enabled,
		DefaultValue: feature.DefaultValue,
		CreatedAt:    feature.CreatedAt,
		UpdatedAt:    feature.UpdatedAt,
		IsActive:     isActive,
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
