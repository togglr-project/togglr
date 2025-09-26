package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateFeatureSchedule(
	ctx context.Context,
	req *generatedapi.CreateFeatureScheduleRequest,
	params generatedapi.CreateFeatureScheduleParams,
) (generatedapi.CreateFeatureScheduleRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Ensure feature exists and get its project
	feature, err := r.featuresUseCase.GetByID(ctx, featureID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature for schedule create failed", "error", err)

		return nil, err
	}

	// Check permissions on the owning project
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

	sch := domain.FeatureSchedule{
		ProjectID:    feature.ProjectID,
		FeatureID:    featureID,
		StartsAt:     optNilDateTimeToPtr(req.StartsAt),
		EndsAt:       optNilDateTimeToPtr(req.EndsAt),
		CronExpr:     optNilStringToPtr(req.CronExpr),
		CronDuration: optNilDurationToPtr(req.CronDuration),
		Timezone:     req.Timezone,
		Action:       domain.FeatureScheduleAction(req.Action),
	}

	created, err := r.featureSchedulesUseCase.Create(ctx, sch)
	if err != nil {
		slog.Error("create feature schedule failed", "error", err)

		return nil, err
	}

	resp := &generatedapi.FeatureScheduleResponse{Schedule: generatedapi.FeatureSchedule{
		ID:           created.ID.String(),
		ProjectID:    created.ProjectID.String(),
		FeatureID:    created.FeatureID.String(),
		StartsAt:     ptrToOptNilDateTime(created.StartsAt),
		EndsAt:       ptrToOptNilDateTime(created.EndsAt),
		CronExpr:     ptrToOptNilString(created.CronExpr),
		CronDuration: ptrToOptNilDuration(created.CronDuration),
		Timezone:     created.Timezone,
		Action:       generatedapi.FeatureScheduleAction(created.Action),
		CreatedAt:    created.CreatedAt,
	}}

	return resp, nil
}
