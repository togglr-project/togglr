package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) UpdateFeatureAlgorithm(
	ctx context.Context,
	req *generatedapi.UpdateFeatureAlgorithmRequest,
	params generatedapi.UpdateFeatureAlgorithmParams,
) (generatedapi.UpdateFeatureAlgorithmRes, error) {
	featureID := domain.FeatureID(params.FeatureID)
	envID := domain.EnvironmentID(params.EnvironmentID)

	env, err := r.environmentsUseCase.GetByIDCached(ctx, envID)
	if err != nil {
		slog.Error("get environment failed", "error", err, "environment_id", params.EnvironmentID)

		return &generatedapi.ErrorNotFound{
			Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("environment not found"),
			},
		}, nil
	}

	feature, err := r.featuresUseCase.GetByIDWithEnv(ctx, featureID, env.Key)
	if err != nil {
		slog.Error("get feature failed", "error", err, "feature_id", featureID)

		return &generatedapi.ErrorNotFound{
			Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			},
		}, nil
	}

	// Check if user can manage the project
	if err := r.permissionsService.CanManageProject(ctx, feature.ProjectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{
			Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			},
		}, nil
	}

	featAlg, err := r.featureAlgorithmsUseCase.GetByFeatureIDWithEnvID(ctx, featureID, envID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature algorithm not found"),
			}}, nil
		}

		slog.Error("failed to get feature algorithm", "featureID", featureID, "envID", envID)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	featAlg.Enabled = req.Enabled
	settings := make(map[string]decimal.Decimal, len(req.Settings))
	for key, value := range req.Settings {
		settings[key] = decimal.NewFromFloat(value)
	}
	featAlg.Settings = settings

	err = r.featureAlgorithmsUseCase.Update(ctx, featAlg)
	if err != nil {
		slog.Error("failed to update feature algorithm",
			"featureID", featureID, "envID", envID, "err", err)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	resp := dto.DomainFeatureAlgorithmToAPI(featAlg)

	return &resp, nil
}
