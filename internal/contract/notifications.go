package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type Emailer interface {
	SendResetPasswordEmail(ctx context.Context, email, token string) error
	Send2FACodeEmail(ctx context.Context, email, code, action string) error
	SendUserNotificationEmail(
		ctx context.Context,
		toEmail string,
		notifType domain.UserNotificationType,
		content domain.UserNotificationContent,
	) error
}
