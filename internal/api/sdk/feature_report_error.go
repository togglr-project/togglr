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

	health, accepted, threshold, err := s.errorReportsUseCase.ReportError(
		ctx,
		projectID,
		featureKey,
		envKey,
		reqCtx,
		req.ErrorType,
		req.ErrorMessage,
	)
	if err != nil {
		return nil, err
	}
	if accepted {
		return &generatedapi.ReportFeatureErrorAccepted{}, nil
	}

	var lastAt generatedapi.OptDateTime
	if !health.LastErrorAt.IsZero() {
		lastAt.SetTo(health.LastErrorAt)
	}

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
