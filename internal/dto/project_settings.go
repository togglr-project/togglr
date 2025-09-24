package dto

import (
	"encoding/json"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainProjectSettingToAPI converts domain ProjectSetting to generated API ProjectSetting
func DomainProjectSettingToAPI(setting domain.ProjectSetting) generatedapi.ProjectSetting {
	// Convert value to any type for JSON marshaling
	var value any
	if setting.Value != nil {
		// If Value is already json.RawMessage, unmarshal it
		if rawMsg, ok := setting.Value.(json.RawMessage); ok {
			json.Unmarshal(rawMsg, &value)
		} else {
			// Otherwise, use the value as is
			value = setting.Value
		}
	}

	return generatedapi.ProjectSetting{
		ID:        setting.ID,
		ProjectID: setting.ProjectID.String(),
		Name:      setting.Name,
		Value:     generatedapi.ProjectSettingValue{},
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}

// DomainProjectSettingsToAPI converts slice of domain ProjectSettings to slice of generated API ProjectSettings
func DomainProjectSettingsToAPI(settings []*domain.ProjectSetting) []generatedapi.ProjectSetting {
	resp := make([]generatedapi.ProjectSetting, 0, len(settings))
	for _, setting := range settings {
		resp = append(resp, DomainProjectSettingToAPI(*setting))
	}
	return resp
}

// APIProjectSettingToDomain converts generated API ProjectSetting to domain ProjectSetting
func APIProjectSettingToDomain(setting generatedapi.ProjectSetting) domain.ProjectSetting {
	// Convert value to any
	// ProjectSettingValue is empty struct, so we'll just use nil
	value := any(nil)

	return domain.ProjectSetting{
		ID:        int(setting.ID),
		ProjectID: domain.ProjectID(setting.ProjectID),
		Name:      setting.Name,
		Value:     value,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}
