package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteFeatureSchedule(
	ctx context.Context,
	params generatedapi.DeleteFeatureScheduleParams,
) (generatedapi.DeleteFeatureScheduleRes, error) {
	id := domain.FeatureScheduleID(params.ScheduleID)

	// Load schedule to know project and to return 404 if it doesn't exist
	schedule, err := r.featureSchedulesUseCase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("schedule not found"),
			}}, nil
		}

		slog.Error("get schedule before delete failed", "error", err, "schedule_id", id)

		return nil, err
	}

	// Check permissions to manage the project
	if err := r.permissionsService.CanManageProject(ctx, schedule.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", schedule.ProjectID)

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

	// Guarded flow: if a feature is guarded, create a pending change and return 202
	pending, conflict, _, err := r.guardCheckAndMaybeCreatePending(
		ctx,
		GuardPendingInput{
			ProjectID:       schedule.ProjectID,
			EnvironmentID:   schedule.EnvironmentID,
			FeatureID:       schedule.FeatureID,
			Reason:          "Delete schedule via API",
			Origin:          "feature-schedule-delete",
			PrimaryEntity:   string(domain.EntityFeatureSchedule),
			PrimaryEntityID: string(id),
			Action:          domain.EntityActionDelete,
			ExtraChanges:    nil,
		},
	)
	if err != nil {
		slog.Error("guard check for schedule delete failed", "error", err)

		return nil, err
	}
	if conflict {
		return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
			Message: generatedapi.NewOptString("Feature is already locked by another pending change"),
		}}, nil
	}
	if pending != nil {
		return pending, nil
	}

	if err := r.featureSchedulesUseCase.Delete(ctx, id); err != nil {
		slog.Error("delete schedule failed", "error", err, "schedule_id", id)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("schedule not found"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteFeatureScheduleNoContent{}, nil
}
