package apisdk

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/shopspring/decimal"

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

	featAlg, err := s.featureAlgorithmsUC.GetByFeatureIDWithEnvID(ctx, feature.ID, feature.EnvironmentID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		slog.Error("get feature algorithm failed", "error", err)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	// For optimizer algorithms, variant_key is optional
	variantKey := ""
	if req.VariantKey.IsSet() {
		variantKey = req.VariantKey.Value
	}

	var algSlug string
	if featAlg.AlgorithmSlug != nil {
		algSlug = *featAlg.AlgorithmSlug
	}

	event := domain.FeedbackEventDTO{
		ProjectID:     feature.ProjectID,
		EnvironmentID: feature.EnvironmentID,
		FeatureID:     feature.ID,
		FeatureKey:    featureKey,
		EnvKey:        envKey,
		VariantKey:    variantKey,
		EventType:     domain.FeedbackEventType(req.EventType),
		AlgorithmSlug: algSlug,
		Reward:        decimal.NewFromFloat32(req.Reward.Or(0.0)),
		Context:       reqCtx,
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
