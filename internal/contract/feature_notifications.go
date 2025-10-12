package contract

import (
	"context"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureNotificationRepository interface {
	AddNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
		featureID domain.FeatureID,
		payload domain.FeatureNotificationPayload,
	) error
	GetByID(ctx context.Context, id domain.FeatureNotificationID) (domain.FeatureNotification, error)
	TakePending(ctx context.Context, limit uint) ([]domain.FeatureNotification, error)
	TakePendingForUpdate(ctx context.Context, limit uint) ([]domain.FeatureNotification, error)
	MarkAsSent(ctx context.Context, id domain.FeatureNotificationID) error
	MarkAsFailed(ctx context.Context, id domain.FeatureNotificationID, reason string) error
	MarkAsSkipped(ctx context.Context, id domain.FeatureNotificationID, reason string) error
	DeleteOld(ctx context.Context, maxAge time.Duration, limit uint) (uint, error)
}

type FeatureNotificationsUseCase interface {
	// Notification Settings
	CreateNotificationSetting(
		ctx context.Context,
		settingDTO domain.NotificationSettingDTO,
	) (domain.NotificationSetting, error)
	GetNotificationSetting(
		ctx context.Context,
		id domain.NotificationSettingID,
	) (domain.NotificationSetting, error)
	UpdateNotificationSetting(
		ctx context.Context,
		setting domain.NotificationSetting,
	) error
	DeleteNotificationSetting(
		ctx context.Context,
		id domain.NotificationSettingID,
	) error
	ListNotificationSettings(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
	) ([]domain.NotificationSetting, error)

	// FeatureNotifications
	TakePendingNotificationsWithSettings(
		ctx context.Context,
		limit uint,
	) ([]domain.FeatureNotificationWithSettings, error)
	MarkNotificationAsSent(ctx context.Context, id domain.FeatureNotificationID) error
	MarkNotificationAsFailed(ctx context.Context, id domain.FeatureNotificationID, reason string) error
	MarkNotificationAsSkipped(ctx context.Context, id domain.FeatureNotificationID, reason string) error

	SendTestNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
		notificationSettingID domain.NotificationSettingID,
	) error
}
