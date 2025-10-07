package apibackend

import (
	"context"

	appctx "github.com/togglr-project/togglr/internal/context"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) MarkAllNotificationsAsRead(ctx context.Context) (generatedapi.MarkAllNotificationsAsReadRes, error) {
	userID := appctx.UserID(ctx)

	err := r.userNotificationsUseCase.MarkAllAsRead(ctx, userID)
	if err != nil {
		return nil, r.NewError(ctx, err)
	}

	return &generatedapi.MarkAllNotificationsAsReadNoContent{}, nil
}
