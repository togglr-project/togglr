package apisdk

import (
	"context"
	"encoding/json"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/sdkserver"
)

func (s *SDKRestAPI) SdkV1FeaturesFeatureKeyEvaluatePost(
	ctx context.Context,
	req generatedapi.EvaluateRequest,
	params generatedapi.SdkV1FeaturesFeatureKeyEvaluatePostParams,
) (generatedapi.SdkV1FeaturesFeatureKeyEvaluatePostRes, error) {
	projectID := etogglcontext.ProjectID(ctx)
	reqCtx := make(map[domain.RuleAttribute]any, len(req))
	for key, valueRaw := range req {
		attr := domain.RuleAttribute(key)
		var value any
		if err := json.Unmarshal(valueRaw, &value); err != nil {
			return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString("invalid request context"),
			}}, nil
		}

		reqCtx[attr] = value
	}

	variant, enabled, found := s.featureProcessor.Evaluate(projectID, params.FeatureKey, reqCtx)
	if !found {
		return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
			Message: generatedapi.NewOptString("feature not found"),
		}}, nil
	}

	return &generatedapi.EvaluateResponse{
		FeatureKey: params.FeatureKey,
		Enabled:    enabled,
		Value:      variant,
	}, nil
}
