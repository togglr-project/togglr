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

func (r *RestAPI) AddFeatureTag(
	ctx context.Context,
	req *generatedapi.AddFeatureTagRequest,
	params generatedapi.AddFeatureTagParams,
) (generatedapi.AddFeatureTagRes, error) {
	userID := appcontext.UserID(ctx)
	featureID := domain.FeatureID(params.FeatureID.String())
	tagID := domain.TagID(req.TagID.String())

	// Resolve feature and environment scope for pending changes.
	// Tags are environment-agnostic; use 'prod' environment for guarded workflow context per guidelines.
	const envKey = "prod"
	feature, err := r.featuresUseCase.GetByIDWithEnv(ctx, featureID, envKey)
	if err != nil {
		slog.Error("get feature for tag add failed", "error", err, "feature_id", featureID)
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString("feature not found"),
				},
			}, nil
		}

		return nil, err
	}
	// Get environment by key within the project's scope
	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, feature.ProjectID, envKey)
	if err != nil {
		slog.Error("get environment for tag add failed",
			"error", err, "project_id", feature.ProjectID, "environment_key", envKey)

		return nil, err
	}

	// Guarded flow: if a feature is guarded, create a pending change and return 202
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
			Reason:        "Add tag to feature via API",
			Origin:        "feature-tag-add",
			Action:        domain.EntityActionInsert,
			OldEntity:     nil, // No old entity for insert
			NewEntity:     featureTagData,
		},
	)
	if err != nil {
		slog.Error("guard check for tag add failed", "error", err)

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

	// Add tag to feature
	err = r.featureTagsUseCase.AddFeatureTag(ctx, featureID, tagID)
	if err != nil {
		slog.Error("add feature tag failed",
			"error", err, "user_id", userID, "feature_id", featureID, "tag_id", tagID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		// Check for the "already associated" error -> 409
		if err.Error() == "tag already associated with feature" {
			return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
				Message: generatedapi.NewOptString("tag already associated with feature"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.AddFeatureTagCreated{}, nil
}
