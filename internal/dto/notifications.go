package dto

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainNotificationSettingToAPI converts domain.NotificationSetting to generatedapi.NotificationSetting.
func DomainNotificationSettingToAPI(
	setting domain.NotificationSetting,
	envKey string,
) generatedapi.NotificationSetting {
	return generatedapi.NotificationSetting{
		ID:             uint(setting.ID),
		ProjectID:      uuid.MustParse(setting.ProjectID.String()),
		EnvironmentID:  uint(setting.EnvironmentID),
		EnvironmentKey: envKey,
		Type:           string(setting.Type),
		Config:         string(setting.Config),
		Enabled:        setting.Enabled,
		CreatedAt:      setting.CreatedAt,
		UpdatedAt:      setting.UpdatedAt,
	}
}

// MakeNotificationSettingDTO converts generatedapi.CreateNotificationSettingRequest to domain.NotificationSettingDTO.
func MakeNotificationSettingDTO(
	req *generatedapi.CreateNotificationSettingRequest,
	projectID domain.ProjectID,
	envID domain.EnvironmentID,
) domain.NotificationSettingDTO {
	return domain.NotificationSettingDTO{
		ProjectID:     projectID,
		EnvironmentID: envID,
		Type:          domain.NotificationType(req.Type),
		Config:        json.RawMessage(req.Config),
		Enabled:       req.Enabled.Value,
	}
}

// UpdateNotificationSettingFromRequest updates a domain.NotificationSetting
// from generatedapi.UpdateNotificationSettingRequest.
func UpdateNotificationSettingFromRequest(
	setting domain.NotificationSetting,
	req *generatedapi.UpdateNotificationSettingRequest,
) domain.NotificationSetting {
	if req.Type.IsSet() {
		setting.Type = domain.NotificationType(req.Type.Value)
	}

	if req.Enabled.Set {
		setting.Enabled = req.Enabled.Value
	}

	if req.Config.Set {
		setting.Config = json.RawMessage(req.Config.Value)
	}

	return setting
}

// MakeListNotificationSettingsResponse converts a slice of domain.NotificationSetting
// to generatedapi.ListNotificationSettingsResponse.
func MakeListNotificationSettingsResponse(
	settings []domain.NotificationSetting,
	envKey string,
) generatedapi.ListNotificationSettingsResponse {
	apiSettings := make([]generatedapi.NotificationSetting, len(settings))
	for i, setting := range settings {
		apiSettings[i] = DomainNotificationSettingToAPI(setting, envKey)
	}

	return generatedapi.ListNotificationSettingsResponse{
		NotificationSettings: apiSettings,
	}
}
