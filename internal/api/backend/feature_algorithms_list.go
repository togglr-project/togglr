package apibackend

import (
	"context"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListFeatureAlgorithms(
	ctx context.Context,
	params generatedapi.ListFeatureAlgorithmsParams,
) (generatedapi.ListFeatureAlgorithmsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)
	envKey := params.EnvironmentKey

	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		slog.Error("get environment failed", "error", err, "env", envKey, "project_id", projectID)

		return &generatedapi.ErrorNotFound{
			Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("environment not found"),
			},
		}, nil
	}

	list, err := r.featureAlgorithmsUseCase.ListByProjectIDWithEnvID(ctx, projectID, env.ID)
	if err != nil {
		slog.Error("list feature algorithms failed", "error", err, "project_id", projectID, "env", envKey)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	result := dto.DomainFeatureAlgorithmsToAPI(list)
	for i := range result {
		feature, err := r.featuresUseCase.GetByIDWithEnv(ctx, domain.FeatureID(result[i].FeatureID), envKey)
		if err != nil {
			slog.Error("get feature failed", "error", err, "id", result[i].FeatureID)

			return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		result[i].Feature = dto.DomainFeatureToAPI(feature)
	}

	return &generatedapi.ListFeatureAlgorithmsResponse{
		FeatureAlgorithms: result,
	}, nil
}
