//nolint:nilerr // false positive
package apibackend

import (
	"context"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/contract"
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

	_, err = r.featureAlgorithmsUseCase.GetByFeatureIDWithEnvID(ctx, featureID, envID)
	if err == nil {
		return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
			Message: generatedapi.NewOptString("feature algorithm already exists")}}, nil
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

	// Guarded flow: if a feature is guarded, create a pending change and return 202
	// The guard engine will automatically compute changes by comparing old and new entities
	settings := make(map[string]decimal.Decimal, len(req.Settings))
	for key, value := range req.Settings {
		settings[key] = decimal.NewFromFloat(value)
	}

	featAlg := domain.FeatureAlgorithm{
		ProjectID:     feature.ProjectID,
		EnvironmentID: feature.EnvironmentID,
		FeatureID:     featureID,
		AlgorithmSlug: &req.AlgorithmSlug,
		Enabled:       req.Enabled,
		Settings:      settings,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
	pendingChange, conflict, _, err := r.guardEngine.CheckGuardedOperation(
		ctx,
		contract.GuardRequest{
			ProjectID:     feature.ProjectID,
			EnvironmentID: feature.EnvironmentID,
			FeatureID:     featureID,
			Reason:        "Create feature algorithm via API",
			Origin:        "feature-algorithms-create",
			Action:        domain.EntityActionInsert,
			NewEntity:     featAlg,
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

	err = r.featureAlgorithmsUseCase.Create(ctx, domain.FeatureAlgorithmDTO{
		ProjectID:     feature.ProjectID,
		EnvironmentID: envID,
		FeatureID:     featureID,
		AlgorithmSlug: &req.AlgorithmSlug,
		Enabled:       req.Enabled,
		Settings:      settings,
	})
	if err != nil {
		slog.Error("create feature algorithm failed", "error", err, "feature_id", featureID)

		return &generatedapi.ErrorInternalServerError{Error: generatedapi.ErrorInternalServerErrorError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.CreateFeatureAlgorithmCreated{}, nil
}
