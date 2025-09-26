package apibackend

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// TestFeatureTimeline handles POST /api/v1/features/{feature_id}/timeline/test.
func (r *RestAPI) TestFeatureTimeline(
	ctx context.Context,
	req *generatedapi.TestFeatureTimelineRequest,
	params generatedapi.TestFeatureTimelineParams,
) (generatedapi.TestFeatureTimelineRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Load a feature and check access to its project
	feature, err := r.featuresUseCase.GetExtendedByID(ctx, featureID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature extended failed", "error", err)

		return nil, err
	}

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

	// Parse time parameters
	loc := params.From.Location()
	if locReq, err := time.LoadLocation(params.Location); err == nil {
		loc = locReq
	} else {
		slog.Error("invalid location", "location", loc)
	}

	from := params.From.In(loc)
	to := params.To.In(loc)

	// Convert test schedules to domain schedules
	testSchedules := make([]domain.FeatureSchedule, 0, len(req.Schedules))

	for _, sched := range req.Schedules {
		schedule := domain.FeatureSchedule{
			ProjectID: feature.ProjectID,
			FeatureID: featureID,
			Timezone:  sched.Timezone,
			Action:    domain.FeatureScheduleAction(sched.Action),
		}

		// Parse starts_at if provided
		if sched.StartsAt.Set {
			if startsAt, err := time.Parse(time.RFC3339, sched.StartsAt.Value.Format(time.RFC3339)); err == nil {
				schedule.StartsAt = &startsAt
			} else {
				slog.Error("invalid starts_at", "starts_at", sched.StartsAt.Value, "error", err)
			}
		}

		// Parse ends_at if provided
		if sched.EndsAt.Set {
			if endsAt, err := time.Parse(time.RFC3339, sched.EndsAt.Value.Format(time.RFC3339)); err == nil {
				schedule.EndsAt = &endsAt
			} else {
				slog.Error("invalid ends_at", "ends_at", sched.EndsAt.Value, "error", err)
			}
		}

		// Set cron_expr if provided
		if sched.CronExpr.Set {
			schedule.CronExpr = &sched.CronExpr.Value
		}

		// Set cron_duration if provided
		if sched.CronDuration.Set {
			if duration, err := time.ParseDuration(sched.CronDuration.Value.String()); err == nil {
				schedule.CronDuration = &duration
			} else {
				slog.Error("invalid cron_duration", "cron_duration", sched.CronDuration.Value, "error", err)
			}
		}

		testSchedules = append(testSchedules, schedule)
	}

	// Create a test feature with the provided schedules
	testFeature := feature
	testFeature.Schedules = testSchedules

	// Build timeline using the feature processor with test data
	events, err := r.featureProcessor.BuildFeatureTimeline(testFeature, from, to)
	if err != nil {
		slog.Error("build feature timeline failed", "error", err)

		return nil, err
	}

	respEvents := make([]generatedapi.FeatureTimelineEvent, 0, len(events))
	for _, e := range events {
		respEvents = append(respEvents, generatedapi.FeatureTimelineEvent{
			Time:    e.Time,
			Enabled: e.Enabled,
		})
	}

	return &generatedapi.FeatureTimelineResponse{Events: respEvents}, nil
}
