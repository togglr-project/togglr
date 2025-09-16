package apibackend

import (
	"context"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// ListAllFeatureSchedules handles GET /api/v1/feature-schedules
func (r *RestAPI) ListAllFeatureSchedules(
	ctx context.Context,
) (generatedapi.ListAllFeatureSchedulesRes, error) {
	// Allow only users with global feature.manage permission
	ok, err := r.permissionsService.HasGlobalPermission(ctx, domain.PermFeatureManage)
	if err != nil {
		slog.Error("check global permission failed", "error", err)
		return nil, err
	}
	if !ok {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	items, err := r.featureSchedulesUseCase.List(ctx)
	if err != nil {
		slog.Error("list all schedules failed", "error", err)
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
