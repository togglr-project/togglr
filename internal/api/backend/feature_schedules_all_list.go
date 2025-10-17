//nolint:nilerr // false positive
package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ListAllFeatureSchedules handles GET /api/v1/projects/{project_id}/env/{environment_key}/feature-schedules.
func (r *RestAPI) ListAllFeatureSchedules(
	ctx context.Context,
	params generatedapi.ListAllFeatureSchedulesParams,
) (generatedapi.ListAllFeatureSchedulesRes, error) {
	projectID := domain.ProjectID(params.ProjectID)
	envKey := params.EnvironmentKey

	// Check if user can view this project
	if err := r.permissionsService.CanViewProject(ctx, projectID); err != nil {
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

	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
			Message: generatedapi.NewOptString("environment not found"),
		}}, nil
	}

	items, err := r.featureSchedulesUseCase.ListByProjectIDEnvID(ctx, projectID, env.ID)
	if err != nil {
		slog.Error("list all schedules failed", "error", err)

		return nil, err
	}

	// Filter schedules based on user's project access
	resp := make(generatedapi.ListFeatureSchedulesResponse, 0, len(items))
	for _, it := range items {
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
