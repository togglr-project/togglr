package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// ListProjectFeatures handles GET /api/v1/projects/{project_id}/features
func (r *RestAPI) ListProjectFeatures(
	ctx context.Context,
	params generatedapi.ListProjectFeaturesParams,
) (generatedapi.ListProjectFeaturesRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check access permissions for the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", projectID)

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

	items, err := r.featuresUseCase.ListByProjectID(ctx, projectID)
	if err != nil {
		slog.Error("list project features failed", "error", err)
		return nil, err
	}

	resp := make(generatedapi.ListFeaturesResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, generatedapi.Feature{
			ID:             it.ID.String(),
			ProjectID:      it.ProjectID.String(),
			Key:            it.Key,
			Name:           it.Name,
			Description:    generatedapi.NewOptNilString(it.Description),
			Kind:           generatedapi.FeatureKind(it.Kind),
			DefaultVariant: it.DefaultVariant,
			Enabled:        it.Enabled,
			CreatedAt:      it.CreatedAt,
			UpdatedAt:      it.UpdatedAt,
		})
	}

	return &resp, nil
}
