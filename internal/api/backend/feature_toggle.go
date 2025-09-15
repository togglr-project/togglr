package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// ToggleFeature handles PUT /api/v1/features/{feature_id}/toggle
func (r *RestAPI) ToggleFeature(
	ctx context.Context,
	req *generatedapi.ToggleFeatureRequest,
	params generatedapi.ToggleFeatureParams,
) (generatedapi.ToggleFeatureRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Ensure feature exists and get project ID
	feature, err := r.featuresUseCase.GetByID(ctx, featureID)
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

	updated, err := r.featuresUseCase.Toggle(ctx, featureID, req.Enabled)
	if err != nil {
		slog.Error("toggle feature failed", "error", err)
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}
		return nil, err
	}

	resp := &generatedapi.FeatureResponse{Feature: generatedapi.Feature{
		ID:             updated.ID.String(),
		ProjectID:      updated.ProjectID.String(),
		Key:            updated.Key,
		Name:           updated.Name,
		Description:    generatedapi.NewOptNilString(updated.Description),
		Kind:           generatedapi.FeatureKind(updated.Kind),
		DefaultVariant: updated.DefaultVariant,
		Enabled:        updated.Enabled,
		CreatedAt:      updated.CreatedAt,
		UpdatedAt:      updated.UpdatedAt,
	}}

	return resp, nil
}
