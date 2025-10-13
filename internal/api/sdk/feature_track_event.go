package apisdk

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/sdkserver"
)

//nolint:nilerr // it's ok here
func (s *SDKRestAPI) TrackFeatureEvent(
	ctx context.Context,
	req *generatedapi.TrackRequest,
	params generatedapi.TrackFeatureEventParams,
) (generatedapi.TrackFeatureEventRes, error) {
	envKey := appcontext.EnvKey(ctx)
	featureKey := params.FeatureKey

	var algorithmIDRef *domain.AlgorithmID
	if req.AlgorithmID.IsSet() {
		algorithmID := domain.AlgorithmID(req.AlgorithmID.Value.String())
		algorithmIDRef = &algorithmID
	}

	var reqCtx map[string]any
	if req.Context.IsSet() {
		reqCtx = make(map[string]any, len(req.Context.Value))

		for key, raw := range req.Context.Value {
			var value any
			if err := json.Unmarshal(raw, &value); err != nil {
				return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("invalid request context"),
				}}, nil
			}

			reqCtx[key] = value
		}
	}

	feature, err := s.featureUseCase.GetByKeyWithEnvCached(ctx, featureKey, envKey)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		slog.Error("get feature failed", "error", err)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	event := domain.FeedbackEventDTO{
		FeatureID:   feature.ID,
		AlgorithmID: algorithmIDRef,
		VariantKey:  req.VariantKey,
		EventType:   req.EventType,
		Reward:      float64(req.Reward.Or(0.0)),
		Context:     reqCtx,
	}

	err = s.bus.PublishFeedbackEvent(ctx, event)
	if err != nil {
		slog.Error("publish feedback event failed", "error", err)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.TrackFeatureEventAccepted{}, nil
}
