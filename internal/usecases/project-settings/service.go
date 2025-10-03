package project_settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

const (
	// Auto-disable project settings keys
	AutoDisableEnabledKey          = "auto_disable_enabled"
	AutoDisableRequiresApprovalKey = "auto_disable_requires_approval"
	AutoDisableErrorThresholdKey   = "auto_disable_error_threshold"
	AutoDisableTimeWindowSecKey    = "auto_disable_time_window_sec"

	// Default values for auto-disable settings
	DefaultAutoDisableEnabled = true
	DefaultRequiresApproval   = false
	DefaultErrorThreshold     = 10
	DefaultTimeWindow         = 60 * time.Second
)

// Service provides project settings management functionality.
type Service struct {
	txManager           db.TxManager
	projectSettingsRepo contract.ProjectSettingsRepository
}

// New creates a new project settings use case.
func New(txManager db.TxManager, projectSettingsRepo contract.ProjectSettingsRepository) *Service {
	return &Service{
		txManager:           txManager,
		projectSettingsRepo: projectSettingsRepo,
	}
}

// Create creates a new project setting.
func (s *Service) Create(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	value any,
) (*domain.ProjectSetting, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("setting name cannot be empty")
	}

	// Marshal value to JSON to validate it
	_, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("invalid setting value: %w", err)
	}

	setting := &domain.ProjectSetting{
		ProjectID: projectID,
		Name:      name,
		Value:     value,
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.projectSettingsRepo.Create(ctx, setting)
	})
	if err != nil {
		return nil, fmt.Errorf("create project setting: %w", err)
	}

	return setting, nil
}

// GetByName retrieves a project setting by name.
func (s *Service) GetByName(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
) (*domain.ProjectSetting, error) {
	if name == "" {
		return nil, errors.New("setting name cannot be empty")
	}

	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("get project setting: %w", err)
	}

	return setting, nil
}

// Update updates a project setting.
func (s *Service) Update(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	value any,
) (*domain.ProjectSetting, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("setting name cannot be empty")
	}

	// Marshal value to JSON to validate it
	_, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("invalid setting value: %w", err)
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.projectSettingsRepo.Update(ctx, projectID, name, value)
	})
	if err != nil {
		return nil, fmt.Errorf("update project setting: %w", err)
	}

	// Return the updated setting
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("get updated project setting: %w", err)
	}

	return setting, nil
}

// Delete deletes a project setting.
func (s *Service) Delete(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
) error {
	if name == "" {
		return errors.New("setting name cannot be empty")
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.projectSettingsRepo.Delete(ctx, projectID, name)
	})
	if err != nil {
		return fmt.Errorf("delete project setting: %w", err)
	}

	return nil
}

// List retrieves project settings with pagination.
func (s *Service) List(
	ctx context.Context,
	projectID domain.ProjectID,
	page, perPage int,
) ([]*domain.ProjectSetting, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 20
	}

	if perPage > 100 {
		perPage = 100
	}

	settings, total, err := s.projectSettingsRepo.List(ctx, projectID, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("list project settings: %w", err)
	}

	return settings, total, nil
}

// GetBoolSetting retrieves a project setting as a boolean value with safe type conversion.
// Returns the default value if the setting doesn't exist or cannot be converted to bool.
func (s *Service) GetBoolSetting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue bool,
) (bool, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case bool:
		return val, nil
	case string:
		// Try to parse string representation of boolean
		if parsed, err := strconv.ParseBool(val); err == nil {
			return parsed, nil
		}
		return defaultValue, nil
	default:
		return defaultValue, nil
	}
}

// GetIntSetting retrieves a project setting as an integer value with safe type conversion.
// Returns the default value if the setting doesn't exist or cannot be converted to int.
func (s *Service) GetIntSetting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue int,
) (int, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		// Try to parse string representation of integer
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed, nil
		}
		return defaultValue, nil
	default:
		return defaultValue, nil
	}
}

// GetStringSetting retrieves a project setting as a string value with safe type conversion.
// Returns the default value if the setting doesn't exist or cannot be converted to string.
func (s *Service) GetStringSetting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue string,
) (string, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case string:
		return val, nil
	case bool:
		return strconv.FormatBool(val), nil
	case int:
		return strconv.Itoa(val), nil
	case int8:
		return strconv.Itoa(int(val)), nil
	case int16:
		return strconv.Itoa(int(val)), nil
	case int32:
		return strconv.Itoa(int(val)), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	default:
		return defaultValue, nil
	}
}

// GetFloat64Setting retrieves a project setting as a float64 value with safe type conversion.
// Returns the default value if the setting doesn't exist or cannot be converted to float64.
func (s *Service) GetFloat64Setting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue float64,
) (float64, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case string:
		// Try to parse string representation of float
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed, nil
		}
		return defaultValue, nil
	default:
		return defaultValue, nil
	}
}

// Auto-disable specific getters with predefined keys and defaults

// GetAutoDisableEnabled retrieves the auto-disable enabled setting for a project.
func (s *Service) GetAutoDisableEnabled(ctx context.Context, projectID domain.ProjectID) (bool, error) {
	return s.GetBoolSetting(ctx, projectID, AutoDisableEnabledKey, DefaultAutoDisableEnabled)
}

// GetAutoDisableRequiresApproval retrieves the auto-disable requires approval setting for a project.
func (s *Service) GetAutoDisableRequiresApproval(ctx context.Context, projectID domain.ProjectID) (bool, error) {
	return s.GetBoolSetting(ctx, projectID, AutoDisableRequiresApprovalKey, DefaultRequiresApproval)
}

// GetAutoDisableErrorThreshold retrieves the auto-disable error threshold setting for a project.
func (s *Service) GetAutoDisableErrorThreshold(ctx context.Context, projectID domain.ProjectID) (int, error) {
	return s.GetIntSetting(ctx, projectID, AutoDisableErrorThresholdKey, DefaultErrorThreshold)
}

// GetAutoDisableTimeWindowSec retrieves the auto-disable time window in seconds setting for a project.
func (s *Service) GetAutoDisableTimeWindowSec(ctx context.Context, projectID domain.ProjectID) (int, error) {
	return s.GetIntSetting(ctx, projectID, AutoDisableTimeWindowSecKey, int(DefaultTimeWindow.Seconds()))
}

// GetAutoDisableTimeWindow retrieves the auto-disable time window as time.Duration for a project.
func (s *Service) GetAutoDisableTimeWindow(ctx context.Context, projectID domain.ProjectID) (time.Duration, error) {
	secs, err := s.GetAutoDisableTimeWindowSec(ctx, projectID)
	if err != nil {
		return DefaultTimeWindow, err
	}
	return time.Duration(secs) * time.Second, nil
}
