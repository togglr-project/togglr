package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) UpdateNotificationSetting(
	ctx context.Context,
	req *generatedapi.UpdateNotificationSettingRequest,
	params generatedapi.UpdateNotificationSettingParams,
) (generatedapi.UpdateNotificationSettingRes, error) {
	settingID := domain.NotificationSettingID(params.SettingID)
	envKey := params.EnvironmentKey

	// Get the existing setting
	setting, err := r.featureNotificationsUseCase.GetNotificationSetting(ctx, settingID)
	if err != nil {
		slog.Error("get notification setting failed", "error", err, "setting_id", settingID)

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

	// Update the setting with the values from the request
	updatedSetting := dto.UpdateNotificationSettingFromRequest(setting, req)

	// Call the service to update the setting
	err = r.featureNotificationsUseCase.UpdateNotificationSetting(ctx, updatedSetting)
	if err != nil {
		slog.Error("update notification setting failed", "error", err, "setting_id", settingID)

		return nil, err
	}

	// Get the updated setting
	setting, err = r.featureNotificationsUseCase.GetNotificationSetting(ctx, settingID)
	if err != nil {
		slog.Error("get updated notification setting failed", "error", err, "setting_id", settingID)

		return nil, err
	}

	// Convert a domain model to an API model
	apiSetting := dto.DomainNotificationSettingToAPI(setting, envKey)

	return &apiSetting, nil
}
