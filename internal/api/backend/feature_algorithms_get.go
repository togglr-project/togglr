package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetFeatureAlgorithm(
	ctx context.Context,
	params generatedapi.GetFeatureAlgorithmParams,
) (generatedapi.GetFeatureAlgorithmRes, error) {
	featureID := domain.FeatureID(params.FeatureID)
	envID := domain.EnvironmentID(params.EnvironmentID)

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

	resp := dto.DomainFeatureAlgorithmToAPI(featAlg)

	return &resp, nil
}
