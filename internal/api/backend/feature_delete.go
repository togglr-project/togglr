package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DeleteFeature handles DELETE /api/v1/features/{feature_id}
func (r *RestAPI) DeleteFeature(
	ctx context.Context,
	params generatedapi.DeleteFeatureParams,
) (generatedapi.DeleteFeatureRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Load feature to know project and to return 404 if it doesn't exist
	feature, err := r.featuresUseCase.GetByID(ctx, featureID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}
		slog.Error("get feature before delete failed", "error", err, "feature_id", featureID)
		return nil, err
	}

	// Check permissions to manage the project
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

	if err := r.featuresUseCase.Delete(ctx, featureID); err != nil {
		slog.Error("delete feature failed", "error", err, "feature_id", featureID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteFeatureNoContent{}, nil
}
