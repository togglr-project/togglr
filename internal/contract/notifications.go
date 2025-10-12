package contract

import (
	"context"
	"encoding/json"

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

type NotificationChannel interface {
	Type() domain.NotificationType
	Send(
		ctx context.Context,
		project *domain.Project,
		feature *domain.Feature,
		envKey string,
		config json.RawMessage,
		payload domain.FeatureNotificationPayload,
	) error
}
