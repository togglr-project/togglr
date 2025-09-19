package apibackend

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// GetFeatureTimeline handles GET /api/v1/features/{feature_id}/timeline
func (r *RestAPI) GetFeatureTimeline(
	ctx context.Context,
	params generatedapi.GetFeatureTimelineParams,
) (generatedapi.GetFeatureTimelineRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Load feature with extended data (schedules needed) and check access to its project
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

	// Build timeline using the feature processor
	loc := params.From.Location()
	if locReq, err := time.LoadLocation(params.Location); err == nil {
		loc = locReq
	} else {
		slog.Error("invalid location", "location", loc)
	}

	from := params.From.In(loc)
	to := params.To.In(loc)

	events, err := r.featureProcessor.BuildFeatureTimeline(feature, from, to)
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
