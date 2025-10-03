package apisdk

import (
	"context"
	"encoding/json"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
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

	err := s.bus.PublishErrorReport(ctx, contract.ErrorReportEvent{
		RequestID:    appcontext.RequestID(ctx),
		ProjectID:    projectID,
		EnvKey:       envKey,
		FeatureKey:   featureKey,
		Context:      reqCtx,
		ErrorType:    req.ErrorType,
		ErrorMessage: req.ErrorMessage,
	})
	if err != nil {
		return &generatedapi.ErrorInternalServerError{
			Error: generatedapi.ErrorInternalServerErrorError{
				Message: generatedapi.NewOptString(err.Error()),
			},
		}, nil
	}

	return &generatedapi.ReportFeatureErrorAccepted{}, nil
}
