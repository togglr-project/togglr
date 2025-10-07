package apibackend

import (
	"context"

	appctx "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetUserNotifications(
	ctx context.Context,
	params generatedapi.GetUserNotificationsParams,
) (generatedapi.GetUserNotificationsRes, error) {
	userID := appctx.UserID(ctx)

	limit := params.Limit.Or(50)
	offset := params.Offset.Or(0)

	notifications, err := r.userNotificationsUseCase.GetUserNotifications(ctx, userID, limit, offset)
	if err != nil {
		return nil, r.NewError(ctx, err)
	}

	dtoNotifications := make([]generatedapi.UserNotification, 0, len(notifications))
	for _, notification := range notifications {
		notif, err := dto.UserNotificationToDTO(notification)
		if err != nil {
			return nil, err
		}

		dtoNotifications = append(dtoNotifications, notif)
	}

	return &generatedapi.UserNotificationsResponse{
		Notifications: dtoNotifications,
		Total:         len(dtoNotifications),
	}, nil
}
