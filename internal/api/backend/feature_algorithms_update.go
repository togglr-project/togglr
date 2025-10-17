//nolint:nilerr // false positive
package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/contract"
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

	// Guarded flow: if a feature is guarded, create a pending change and return 202
	// The guard engine will automatically compute changes by comparing old and new entities
	settings := make(map[string]decimal.Decimal, len(req.Settings))
	for key, value := range req.Settings {
		settings[key] = decimal.NewFromFloat(value)
	}

	featAlgNew := domain.FeatureAlgorithm{
		ID:            featAlg.ID,
		ProjectID:     feature.ProjectID,
		EnvironmentID: feature.EnvironmentID,
		FeatureID:     featureID,
		AlgorithmSlug: featAlg.AlgorithmSlug,
		Enabled:       req.Enabled,
		Settings:      settings,
	}
	pendingChange, conflict, _, err := r.guardEngine.CheckGuardedOperation(
		ctx,
		contract.GuardRequest{
			ProjectID:     feature.ProjectID,
			EnvironmentID: feature.EnvironmentID,
			FeatureID:     featureID,
			Reason:        "Update feature algorithm via API",
			Origin:        "feature-algorithms-update",
			Action:        domain.EntityActionUpdate,
			OldEntity:     featAlg,
			NewEntity:     featAlgNew,
		},
	)
	if err != nil {
		slog.Error("guard check for schedule update failed", "error", err)

		return nil, err
	}
	if conflict {
		return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
			Message: generatedapi.NewOptString("Feature is already locked by another pending change"),
		}}, nil
	}
	if pendingChange != nil {
		resp := convertPendingChangeToResponse(pendingChange)

		return &resp, nil
	}

	err = r.featureAlgorithmsUseCase.Update(ctx, featAlgNew)
	if err != nil {
		slog.Error("failed to update feature algorithm",
			"featureID", featureID, "envID", envID, "err", err)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	resp := dto.DomainFeatureAlgorithmToAPI(featAlgNew)

	return &resp, nil
}
