package contract

import (
	"context"
	"encoding/json"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureNotificationRepository interface {
	AddNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
		featureID domain.FeatureID,
		payload json.RawMessage,
	) error
	GetByID(ctx context.Context, id domain.FeatureNotificationID) (domain.FeatureNotification, error)
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

	SendTestNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
		notificationSettingID domain.NotificationSettingID,
	) error
}
