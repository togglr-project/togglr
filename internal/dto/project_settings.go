package dto

import (
	"encoding/json"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainProjectSettingToAPI converts domain ProjectSetting to generated API ProjectSetting.
func DomainProjectSettingToAPI(setting domain.ProjectSetting) generatedapi.ProjectSetting {
	var value string
	if setting.Value != nil {
		data, _ := json.Marshal(setting.Value) //nolint:errchkjson // it's ok
		value = string(data)
	}

	return generatedapi.ProjectSetting{
		ID:        setting.ID,
		ProjectID: setting.ProjectID.String(),
		Name:      setting.Name,
		Value:     value,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}

// DomainProjectSettingsToAPI converts slice of domain ProjectSettings to slice of generated API ProjectSettings.
func DomainProjectSettingsToAPI(settings []*domain.ProjectSetting) []generatedapi.ProjectSetting {
	resp := make([]generatedapi.ProjectSetting, 0, len(settings))
	for _, setting := range settings {
		resp = append(resp, DomainProjectSettingToAPI(*setting))
	}

	return resp
}

// APIProjectSettingToDomain converts generated API ProjectSetting to domain ProjectSetting.
func APIProjectSettingToDomain(setting generatedapi.ProjectSetting) domain.ProjectSetting {
	var value any
	_ = json.Unmarshal([]byte(setting.Value), &value)

	return domain.ProjectSetting{
		ID:        setting.ID,
		ProjectID: domain.ProjectID(setting.ProjectID),
		Name:      setting.Name,
		Value:     value,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}
