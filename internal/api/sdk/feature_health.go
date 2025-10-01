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

	_ = projectID

	// TODO: implement

	return &generatedapi.FeatureHealth{
		FeatureKey:     featureKey,
		EnvironmentKey: envKey,
		Enabled:        false,                      // TODO: implement
		AutoDisabled:   false,                      // TODO: implement
		ErrorRate:      generatedapi.OptFloat32{},  // TODO: implement
		Threshold:      generatedapi.OptFloat32{},  // TODO: implement
		LastErrorAt:    generatedapi.OptDateTime{}, // TODO: implement
	}, nil
}
