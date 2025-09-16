package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// ListFeatureSchedules handles GET /api/v1/features/{feature_id}/schedules
func (r *RestAPI) ListFeatureSchedules(
	ctx context.Context,
	params generatedapi.ListFeatureSchedulesParams,
) (generatedapi.ListFeatureSchedulesRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Ensure feature exists and get its project to check access rights
	feature, err := r.featuresUseCase.GetByID(ctx, featureID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}
		slog.Error("get feature for list schedules failed", "error", err)
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

	items, err := r.featureSchedulesUseCase.ListByFeatureID(ctx, featureID)
	if err != nil {
		slog.Error("list schedules by feature failed", "error", err)
		return nil, err
	}

	resp := make(generatedapi.ListFeatureSchedulesResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, generatedapi.FeatureSchedule{
			ID:        it.ID.String(),
			ProjectID: it.ProjectID.String(),
			FeatureID: it.FeatureID.String(),
			StartsAt:  ptrToOptNilDateTime(it.StartsAt),
			EndsAt:    ptrToOptNilDateTime(it.EndsAt),
			CronExpr:  ptrToOptNilString(it.CronExpr),
			Timezone:  it.Timezone,
			Action:    generatedapi.FeatureScheduleAction(it.Action),
			CreatedAt: it.CreatedAt,
		})
	}

	return &resp, nil
}
