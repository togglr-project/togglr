package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) UpdateFeatureSchedule(
	ctx context.Context,
	req *generatedapi.UpdateFeatureScheduleRequest,
	params generatedapi.UpdateFeatureScheduleParams,
) (generatedapi.UpdateFeatureScheduleRes, error) {
	id := domain.FeatureScheduleID(params.ScheduleID)

	// Load existing schedule to validate existence and determine project for permissions.
	current, err := r.featureSchedulesUseCase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("schedule not found"),
			}}, nil
		}

		slog.Error("get schedule before update failed", "error", err, "schedule_id", id)

		return nil, err
	}

	// Check permissions on the owning project
	if err := r.permissionsService.CanManageProject(ctx, current.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", current.ProjectID)

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

	sch := domain.FeatureSchedule{
		ID:           id,
		ProjectID:    current.ProjectID,
		FeatureID:    current.FeatureID,
		StartsAt:     optNilDateTimeToPtr(req.StartsAt),
		EndsAt:       optNilDateTimeToPtr(req.EndsAt),
		CronExpr:     optNilStringToPtr(req.CronExpr),
		CronDuration: optNilDurationToPtr(req.CronDuration),
		Timezone:     req.Timezone,
		Action:       domain.FeatureScheduleAction(req.Action),
	}

	// Build changes diff for pending change payload
	changes := buildFeatureScheduleChangeDiff(current, sch)

	// Guarded flow: if feature is guarded, create a pending change and return 202
	pending, conflict, _, err := r.guardCheckAndMaybeCreatePending(
		ctx,
		GuardPendingInput{
			ProjectID:       sch.ProjectID,
			EnvironmentID:   current.EnvironmentID,
			FeatureID:       sch.FeatureID,
			Reason:          "Update schedule via API",
			Origin:          "feature-schedule-update",
			PrimaryEntity:   string(domain.EntityFeatureSchedule),
			PrimaryEntityID: string(id),
			Action:          domain.EntityActionUpdate,
			ExtraChanges:    changes,
		},
	)
	if err != nil {
		slog.Error("guard check for schedule update failed", "error", err)

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

	updated, err := r.featureSchedulesUseCase.Update(ctx, sch)
	if err != nil {
		slog.Error("update schedule failed", "error", err, "schedule_id", id)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("schedule not found"),
			}}, nil
		}

		return nil, err
	}

	resp := &generatedapi.FeatureScheduleResponse{Schedule: generatedapi.FeatureSchedule{
		ID:           updated.ID.String(),
		ProjectID:    updated.ProjectID.String(),
		FeatureID:    updated.FeatureID.String(),
		StartsAt:     ptrToOptNilDateTime(updated.StartsAt),
		EndsAt:       ptrToOptNilDateTime(updated.EndsAt),
		CronExpr:     ptrToOptNilString(updated.CronExpr),
		CronDuration: ptrToOptNilDuration(updated.CronDuration),
		Timezone:     updated.Timezone,
		Action:       generatedapi.FeatureScheduleAction(updated.Action),
		CreatedAt:    updated.CreatedAt,
	}}

	return resp, nil
}
