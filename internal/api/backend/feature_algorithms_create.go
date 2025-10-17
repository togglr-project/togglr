package apibackend

import (
	"context"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateFeatureAlgorithm(
	ctx context.Context,
	req *generatedapi.CreateFeatureAlgorithmRequest,
	params generatedapi.CreateFeatureAlgorithmParams,
) (generatedapi.CreateFeatureAlgorithmRes, error) {
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

	err = r.featureAlgorithmsUseCase.Create(ctx, domain.FeatureAlgorithmDTO{
		ProjectID:     feature.ProjectID,
		EnvironmentID: envID,
		FeatureID:     featureID,
		AlgorithmSlug: req.AlgorithmSlug,
		Enabled:       req.Enabled,
		Settings:      nil,
	})
	if err != nil {
		slog.Error("create feature algorithm failed", "error", err, "feature_id", featureID)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.CreateFeatureAlgorithmCreated{}, nil
}
