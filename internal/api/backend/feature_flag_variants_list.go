package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ListFeatureFlagVariants handles GET /api/v1/features/{feature_id}/variants
func (r *RestAPI) ListFeatureFlagVariants(
	ctx context.Context,
	params generatedapi.ListFeatureFlagVariantsParams,
) (generatedapi.ListFeatureFlagVariantsRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Ensure feature exists and get its project to check access rights
	feature, err := r.featuresUseCase.GetByID(ctx, featureID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}
		slog.Error("get feature for list variants failed", "error", err)
		return nil, err
	}

	// Check access permissions for the owning project
	if err := r.permissionsService.CanAccessProject(ctx, feature.ProjectID); err != nil {
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

	items, err := r.flagVariantsUseCase.ListByFeatureID(ctx, featureID)
	if err != nil {
		slog.Error("list flag variants by feature failed", "error", err)
		return nil, err
	}

	resp := make(generatedapi.ListFlagVariantsResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, generatedapi.FlagVariant{
			ID:             it.ID.String(),
			FeatureID:      it.FeatureID.String(),
			Name:           it.Name,
			RolloutPercent: int(it.RolloutPercent),
		})
	}

	return &resp, nil
}
