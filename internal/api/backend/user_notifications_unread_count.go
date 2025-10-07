package apibackend

import (
	"context"

	appctx "github.com/togglr-project/togglr/internal/context"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetUnreadNotificationsCount(
	ctx context.Context,
) (generatedapi.GetUnreadNotificationsCountRes, error) {
	userID := appctx.UserID(ctx)

	count, err := r.userNotificationsUseCase.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, r.NewError(ctx, err)
	}

	return &generatedapi.UnreadCountResponse{
		Count: count,
	}, nil
}
