package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) RemoveFeatureTag(
	ctx context.Context,
	params generatedapi.RemoveFeatureTagParams,
) (generatedapi.RemoveFeatureTagRes, error) {
	userID := appcontext.UserID(ctx)
	featureID := domain.FeatureID(params.FeatureID.String())
	tagID := domain.TagID(params.TagID.String())

	// Resolve feature and environment scope for pending changes.
	// Tags are environment-agnostic; use 'prod' environment for guarded workflow context.
	const envKey = "prod"
	feature, err := r.featuresUseCase.GetByIDWithEnv(ctx, featureID, envKey)
	if err != nil {
		slog.Error("get feature for tag remove failed", "error", err, "feature_id", featureID)
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		return nil, err
	}
	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, feature.ProjectID, envKey)
	if err != nil {
		slog.Error("get environment for tag remove failed",
			"error", err, "project_id", feature.ProjectID, "environment_key", envKey)

		return nil, err
	}

	// Guarded flow: if feature is guarded, create a pending change and return 202
	// Create a simple struct for feature tag relationship
	featureTagData := domain.FeatureTags{
		FeatureID: featureID,
		TagID:     tagID,
	}

	pendingChange, conflict, _, err := r.guardEngine.CheckGuardedOperation(
		ctx,
		contract.GuardRequest{
			ProjectID:     feature.ProjectID,
			EnvironmentID: env.ID,
			FeatureID:     featureID,
			Reason:        "Remove tag from feature via API",
			Origin:        "feature-tag-remove",
			Action:        domain.EntityActionDelete,
			OldEntity:     featureTagData, // For delete, we need the old entity
			NewEntity:     nil,            // No new entity for delete
		},
	)
	if err != nil {
		slog.Error("guard check for tag remove failed", "error", err)

		return nil, err
	}
	if conflict {
		return &generatedapi.ErrorConflict{
			Error: generatedapi.ErrorConflictError{
				Message: generatedapi.NewOptString("Feature is already locked by another pending change"),
			},
		}, nil
	}
	if pendingChange != nil {
		resp := convertPendingChangeToResponse(pendingChange)

		return &resp, nil
	}

	// Remove tag from the feature
	err = r.featureTagsUseCase.RemoveFeatureTag(ctx, featureID, tagID)
	if err != nil {
		slog.Error("remove feature tag failed", "error", err, "user_id", userID, "feature_id", featureID, "tag_id", tagID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.RemoveFeatureTagNoContent{}, nil
}
