package apibackend

import (
	"context"
	"log/slog"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ListAllFeatureSchedules handles GET /api/v1/feature-schedules.
func (r *RestAPI) ListAllFeatureSchedules(
	ctx context.Context,
) (generatedapi.ListAllFeatureSchedulesRes, error) {
	items, err := r.featureSchedulesUseCase.List(ctx)
	if err != nil {
		slog.Error("list all schedules failed", "error", err)

		return nil, err
	}

	// Filter schedules based on user's project access
	resp := make(generatedapi.ListFeatureSchedulesResponse, 0, len(items))
	for _, it := range items {
		// Check if user can view this project
		if err := r.permissionsService.CanViewProject(ctx, it.ProjectID); err != nil {
			// Skip this schedule if user doesn't have access to the project
			continue
		}

		resp = append(resp, generatedapi.FeatureSchedule{
			ID:           it.ID.String(),
			ProjectID:    it.ProjectID.String(),
			FeatureID:    it.FeatureID.String(),
			StartsAt:     ptrToOptNilDateTime(it.StartsAt),
			EndsAt:       ptrToOptNilDateTime(it.EndsAt),
			CronExpr:     ptrToOptNilString(it.CronExpr),
			CronDuration: ptrToOptNilDuration(it.CronDuration),
			Timezone:     it.Timezone,
			Action:       generatedapi.FeatureScheduleAction(it.Action),
			CreatedAt:    it.CreatedAt,
		})
	}

	return &resp, nil
}
