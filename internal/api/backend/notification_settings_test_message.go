package apibackend

import (
	"context"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) SendTestNotification(
	ctx context.Context,
	params generatedapi.SendTestNotificationParams,
) (generatedapi.SendTestNotificationRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	envKey := params.EnvironmentKey

	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		slog.Error("get environment failed", "error", err, "project_id", projectID, "env_key", envKey)

		return nil, err
	}

	err = r.featureNotificationsUseCase.SendTestNotification(
		ctx,
		projectID,
		env.ID,
		domain.NotificationSettingID(params.SettingID),
	)
	if err != nil {
		slog.Error("failed to send test notification", "error", err)

		return nil, err
	}

	return &generatedapi.SendTestNotificationNoContent{}, nil
}
