//nolint:nilerr // false positive
package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteFeatureAlgorithm(
	ctx context.Context,
	params generatedapi.DeleteFeatureAlgorithmParams,
) (generatedapi.DeleteFeatureAlgorithmRes, error) {
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
	// The guard engine will automatically handle the delete operation
	pendingChange, conflict, _, err := r.guardEngine.CheckGuardedOperation(
		ctx,
		contract.GuardRequest{
			ProjectID:     feature.ProjectID,
			EnvironmentID: feature.EnvironmentID,
			FeatureID:     feature.ID,
			Reason:        "Delete feature algorithm via API",
			Origin:        "feature-algorithm-delete",
			Action:        domain.EntityActionDelete,
			OldEntity:     featAlg, // For delete, we need the old entity
			NewEntity:     nil,     // No new entity for delete
		},
	)
	if err != nil {
		slog.Error("guard check for schedule delete failed", "error", err)

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

	err = r.featureAlgorithmsUseCase.DeleteByFeatureIDWithEnvID(ctx, featureID, envID)
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
