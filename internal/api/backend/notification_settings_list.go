package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListNotificationSettings(
	ctx context.Context,
	params generatedapi.ListNotificationSettingsParams,
) (generatedapi.ListNotificationSettingsRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	envKey := params.EnvironmentKey

	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		default:
			slog.Error("get environment failed", "error", err, "project_id", projectID, "env_key", envKey)

			return nil, err
		}
	}

	// Call the service to list settings
	settings, err := r.featureNotificationsUseCase.ListNotificationSettings(ctx, projectID, env.ID)
	if err != nil {
		slog.Error("list notification settings failed", "error", err, "project_id", projectID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		default:
			return nil, err
		}
	}

	// Convert domain models to API response
	response := dto.MakeListNotificationSettingsResponse(settings, envKey)

	return &response, nil
}
