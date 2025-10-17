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
