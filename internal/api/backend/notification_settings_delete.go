package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteNotificationSetting(
	ctx context.Context,
	params generatedapi.DeleteNotificationSettingParams,
) (generatedapi.DeleteNotificationSettingRes, error) {
	settingID := domain.NotificationSettingID(params.SettingID)

	// Call the service to delete the setting
	err := r.featureNotificationsUseCase.DeleteNotificationSetting(ctx, settingID)
	if err != nil {
		slog.Error("delete notification setting failed", "error", err, "setting_id", settingID)
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

	// Return a success response (204 No Content)
	return &generatedapi.DeleteNotificationSettingNoContent{}, nil
}
