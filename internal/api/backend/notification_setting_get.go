package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetNotificationSetting(
	ctx context.Context,
	params generatedapi.GetNotificationSettingParams,
) (generatedapi.GetNotificationSettingRes, error) {
	settingID := domain.NotificationSettingID(params.SettingID)
	envKey := params.EnvironmentKey

	// Call the service
	setting, err := r.featureNotificationsUseCase.GetNotificationSetting(ctx, settingID)
	if err != nil {
		slog.Error("get notification setting failed", "error", err, "setting_id", settingID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		}

		return nil, err
	}

	// Convert a domain model to an API model
	apiSetting := dto.DomainNotificationSettingToAPI(setting, envKey)

	return &apiSetting, nil
}
