package contract

import (
	"context"
	"encoding/json"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type UserNotificationsUseCase interface {
	CreateNotification(
		ctx context.Context,
		userID domain.UserID,
		notificationType domain.UserNotificationType,
		content domain.UserNotificationContent,
	) error
	GetUserNotifications(
		ctx context.Context,
		userID domain.UserID,
		limit, offset uint,
	) ([]domain.UserNotification, error)
	GetUnreadCount(ctx context.Context, userID domain.UserID) (uint, error)
	MarkAsRead(ctx context.Context, notificationID domain.UserNotificationID) error
	MarkAllAsRead(ctx context.Context, userID domain.UserID) error
	DeleteOldNotifications(ctx context.Context, maxAge time.Duration, limit uint) (uint, error)
	TakePendingEmailNotifications(ctx context.Context, limit uint) ([]domain.UserNotification, error)
	MarkEmailAsSent(ctx context.Context, id domain.UserNotificationID) error
	MarkEmailAsFailed(ctx context.Context, id domain.UserNotificationID, reason string) error
}

type UserNotificationsRepository interface {
	Create(
		ctx context.Context,
		userID domain.UserID,
		notificationType domain.UserNotificationType,
		content json.RawMessage,
	) (domain.UserNotification, error)
	GetByID(ctx context.Context, id domain.UserNotificationID) (domain.UserNotification, error)
	GetByUserID(ctx context.Context, userID domain.UserID, limit, offset uint) ([]domain.UserNotification, error)
	GetUnreadCount(ctx context.Context, userID domain.UserID) (uint, error)
	MarkAsRead(ctx context.Context, id domain.UserNotificationID) error
	MarkAllAsRead(ctx context.Context, userID domain.UserID) error
	DeleteOld(ctx context.Context, maxAge time.Duration, limit uint) (uint, error)
	GetPendingEmailNotifications(ctx context.Context, limit uint) ([]domain.UserNotification, error)
	GetPendingEmailNotificationsForUpdate(ctx context.Context, limit uint) ([]domain.UserNotification, error)
	MarkEmailAsSent(ctx context.Context, id domain.UserNotificationID) error
	MarkEmailAsFailed(ctx context.Context, id domain.UserNotificationID, reason string) error
}
