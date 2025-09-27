package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateFeatureFlagVariant(
	ctx context.Context,
	req *generatedapi.CreateFlagVariantRequest,
	params generatedapi.CreateFeatureFlagVariantParams,
) (generatedapi.CreateFeatureFlagVariantRes, error) {
	featureID := domain.FeatureID(params.FeatureID)
	environmentKey := params.EnvironmentKey

	feature, err := r.featuresUseCase.GetByIDWithEnv(ctx, featureID, environmentKey)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature for variant create failed", "error", err)

		return nil, err
	}

	// Check permissions on the owning project
	if err := r.permissionsService.CanManageProject(ctx, feature.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", feature.ProjectID)

		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			}}, nil
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		}

		return nil, err
	}

	variant := domain.FlagVariant{
		ProjectID:      feature.ProjectID,
		FeatureID:      featureID,
		Name:           req.Name,
		RolloutPercent: uint8(req.RolloutPercent),
	}

	// Resolve environment
	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, feature.ProjectID, environmentKey)
	if err != nil {
		slog.Error("get environment for flag variant create failed", "error", err)

		return nil, err
	}

	// Guarded flow: if feature is guarded, create a pending change and return 202
	pc, conflict, _, err := r.guardEngine.CheckAndMaybeCreatePending(
		ctx,
		contract.GuardEngineInput{
			ProjectID:       feature.ProjectID,
			EnvironmentID:   env.ID,
			FeatureID:       featureID,
			Reason:          "Create flag variant via API",
			Origin:          "flag-variant-create",
			PrimaryEntity:   string(domain.EntityFlagVariant),
			PrimaryEntityID: "",
			Action:          domain.EntityActionInsert,
			ExtraChanges:    nil,
		},
	)
	if err != nil {
		slog.Error("guard check for flag variant create failed", "error", err)

		return nil, err
	}
	if conflict {
		return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
			Message: generatedapi.NewOptString("Feature is already locked by another pending change"),
		}}, nil
	}
	if pc != nil {
		resp := convertPendingChangeToResponse(pc)

		return &resp, nil
	}

	created, err := r.flagVariantsUseCase.Create(ctx, variant)
	if err != nil {
		slog.Error("create flag variant failed", "error", err)

		return nil, err
	}

	resp := &generatedapi.FlagVariantResponse{FlagVariant: generatedapi.FlagVariant{
		ID:             created.ID.String(),
		FeatureID:      created.FeatureID.String(),
		Name:           created.Name,
		RolloutPercent: int(created.RolloutPercent),
	}}

	return resp, nil
}
