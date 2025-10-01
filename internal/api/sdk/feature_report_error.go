package apisdk

import (
	"context"
	"encoding/json"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/sdkserver"
)

func (s *SDKRestAPI) ReportFeatureError(
	ctx context.Context,
	req *generatedapi.FeatureErrorReport,
	params generatedapi.ReportFeatureErrorParams,
) (generatedapi.ReportFeatureErrorRes, error) {
	projectID := appcontext.ProjectID(ctx)
	envKey := appcontext.EnvKey(ctx)
	featureKey := params.FeatureKey
	reqCtx := make(map[domain.RuleAttribute]any)

	if req.Context.IsSet() {
		for key, valueRaw := range req.Context.Value {
			attr := domain.RuleAttribute(key)

			var value any
			if err := json.Unmarshal(valueRaw, &value); err != nil {
				return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("invalid request context"),
				}}, nil
			}

			reqCtx[attr] = value
		}
	}

	_ = projectID
	// TODO: implement

	// return &generatedapi.ReportFeatureErrorAccepted{}, nil // for 202 status code (Error reported, evaluation pending)

	// for 200 status code (Error reported successfully)
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
