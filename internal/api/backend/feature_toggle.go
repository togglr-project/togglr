package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ToggleFeature handles PUT /api/v1/features/{feature_id}/toggle.
func (r *RestAPI) ToggleFeature(
	ctx context.Context,
	req *generatedapi.ToggleFeatureRequest,
	params generatedapi.ToggleFeatureParams,
) (generatedapi.ToggleFeatureRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Get environment key from query parameters
	environmentKey := params.EnvironmentKey

	// Ensure a feature exists and get project ID
	feature, err := r.featuresUseCase.GetByIDWithEnvironment(ctx, featureID, environmentKey)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature for toggle failed", "error", err)

		return nil, err
	}

	// Check permission to toggle feature within the project's scope
	ok, perr := r.permissionsService.HasProjectPermission(ctx, feature.ProjectID, domain.PermFeatureToggle)
	if perr != nil {
		slog.Error("permission check failed", "error", perr, "project_id", feature.ProjectID)

		if errors.Is(perr, domain.ErrUserNotFound) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		}

		return nil, perr
	}

	if !ok {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	updated, guardResult, err := r.featuresUseCase.Toggle(ctx, featureID, req.Enabled, environmentKey)
	if err != nil {
		slog.Error("toggle feature failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		return nil, err
	}

	// Handle guard result
	if guardResult.Error != nil {
		slog.Error("guard check failed", "error", guardResult.Error)

		return nil, guardResult.Error
	}

	if guardResult.ChangeConflict {
		return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
			Message: generatedapi.NewOptString("Feature is already locked by another pending change"),
		}}, nil
	}

	if guardResult.Pending {
		// Convert pending change to response
		pendingChangeResp := convertPendingChangeToResponse(guardResult.PendingChange)

		return &pendingChangeResp, nil
	}

	resp := &generatedapi.FeatureResponse{Feature: generatedapi.Feature{
		ID:          updated.ID.String(),
		ProjectID:   updated.ProjectID.String(),
		Key:         updated.Key,
		Name:        updated.Name,
		Description: generatedapi.NewOptNilString(updated.Description),
		Kind:        generatedapi.FeatureKind(updated.Kind),
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}}

	return resp, nil
}
