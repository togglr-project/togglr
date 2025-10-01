package apisdk

import (
	"context"

	appcontext "github.com/togglr-project/togglr/internal/context"
	generatedapi "github.com/togglr-project/togglr/internal/generated/sdkserver"
)

func (s *SDKRestAPI) GetFeatureHealth(
	ctx context.Context,
	params generatedapi.GetFeatureHealthParams,
) (generatedapi.GetFeatureHealthRes, error) {
	projectID := appcontext.ProjectID(ctx)
	envKey := appcontext.EnvKey(ctx)
	featureKey := params.FeatureKey

	health, err := s.errorReportsUseCase.GetFeatureHealth(ctx, projectID, featureKey, envKey)
	if err != nil {
		return nil, err
	}

	var lastAt generatedapi.OptDateTime
	if !health.LastErrorAt.IsZero() {
		lastAt.SetTo(health.LastErrorAt)
	}

	// Get threshold from project settings
	threshold := 20 // default threshold
	// TODO: implement threshold retrieval from project settings

	return &generatedapi.FeatureHealth{
		FeatureKey:     featureKey,
		EnvironmentKey: envKey,
		Enabled:        health.Enabled,
		AutoDisabled:   !health.Enabled,
		ErrorRate:      generatedapi.NewOptFloat32(float32(health.ErrorRate)),
		Threshold:      generatedapi.NewOptFloat32(float32(threshold)),
		LastErrorAt:    lastAt,
	}, nil
}
