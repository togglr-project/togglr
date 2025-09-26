package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetFeatureSchedule(
	ctx context.Context,
	params generatedapi.GetFeatureScheduleParams,
) (generatedapi.GetFeatureScheduleRes, error) {
	id := domain.FeatureScheduleID(params.ScheduleID)

	item, err := r.featureSchedulesUseCase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("schedule not found"),
			}}, nil
		}

		slog.Error("get schedule failed", "error", err)

		return nil, err
	}

	// Check access permissions for the owning project
	if err := r.permissionsService.CanAccessProject(ctx, item.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", item.ProjectID)

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

	resp := &generatedapi.FeatureScheduleResponse{Schedule: generatedapi.FeatureSchedule{
		ID:           item.ID.String(),
		ProjectID:    item.ProjectID.String(),
		FeatureID:    item.FeatureID.String(),
		StartsAt:     ptrToOptNilDateTime(item.StartsAt),
		EndsAt:       ptrToOptNilDateTime(item.EndsAt),
		CronExpr:     ptrToOptNilString(item.CronExpr),
		CronDuration: ptrToOptNilDuration(item.CronDuration),
		Timezone:     item.Timezone,
		Action:       generatedapi.FeatureScheduleAction(item.Action),
		CreatedAt:    item.CreatedAt,
	}}

	return resp, nil
}
