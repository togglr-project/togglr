package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type NotificationSettingsRepository interface {
	CreateSetting(
		ctx context.Context,
		settingDTO domain.NotificationSettingDTO,
	) (domain.NotificationSetting, error)
	GetSettingByID(
		ctx context.Context,
		id domain.NotificationSettingID,
	) (domain.NotificationSetting, error)
	UpdateSetting(ctx context.Context, setting domain.NotificationSetting) error
	DeleteSetting(ctx context.Context, id domain.NotificationSettingID) error
	ListSettings(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
	) ([]domain.NotificationSetting, error)
	ListSettingsAll(
		ctx context.Context,
		projectID domain.ProjectID,
	) ([]domain.NotificationSetting, error)
	CountSettings(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
	) (uint, error)
}
