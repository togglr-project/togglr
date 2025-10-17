package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteFeatureAlgorithm(
	ctx context.Context,
	params generatedapi.DeleteFeatureAlgorithmParams,
) (generatedapi.DeleteFeatureAlgorithmRes, error) {
	featureID := domain.FeatureID(params.FeatureID)
	envID := domain.EnvironmentID(params.EnvironmentID)

	err := r.featureAlgorithmsUseCase.DeleteByFeatureIDWithEnvID(ctx, featureID, envID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature algorithm not found"),
			}}, nil
		}

		slog.Error("failed to delete feature algorithm", "featureID", featureID, "envID", envID)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.DeleteFeatureAlgorithmNoContent{}, nil
}
