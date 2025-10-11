package notification_settings

import (
	"encoding/json"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type notificationSettingModel struct {
	ID            uint            `db:"id"`
	ProjectID     string          `db:"project_id"`
	EnvironmentID uint64          `db:"environment_id"`
	Type          string          `db:"type"`
	Config        json.RawMessage `db:"config"`
	Enabled       bool            `db:"enabled"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}

func (m *notificationSettingModel) toDomain() domain.NotificationSetting {
	return domain.NotificationSetting{
		ID:            domain.NotificationSettingID(m.ID),
		ProjectID:     domain.ProjectID(m.ProjectID),
		EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
		Type:          domain.NotificationType(m.Type),
		Config:        m.Config,
		Enabled:       m.Enabled,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func settingFromDTO(dto domain.NotificationSettingDTO) notificationSettingModel {
	now := time.Now()

	if dto.Config == nil {
		dto.Config = json.RawMessage("{}")
	}

	return notificationSettingModel{
		ProjectID:     dto.ProjectID.String(),
		EnvironmentID: uint64(dto.EnvironmentID),
		Type:          string(dto.Type),
		Config:        dto.Config,
		Enabled:       dto.Enabled,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func settingFromDomain(setting domain.NotificationSetting) notificationSettingModel {
	return notificationSettingModel{
		ID:            uint(setting.ID),
		ProjectID:     setting.ProjectID.String(),
		EnvironmentID: uint64(setting.EnvironmentID),
		Type:          string(setting.Type),
		Config:        setting.Config,
		Enabled:       setting.Enabled,
		CreatedAt:     setting.CreatedAt,
		UpdatedAt:     setting.UpdatedAt,
	}
}
