package contract

import (
	"context"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type ProjectSettingsRepository interface {
	// Set sets a project setting
	Set(
		ctx context.Context,
		projectID domain.ProjectID,
		name string,
		value interface{},
	) error

	// Get retrieves a project setting
	Get(
		ctx context.Context,
		projectID domain.ProjectID,
		name string,
	) (domain.ProjectSetting, error)

	// GetAll retrieves all settings for a project
	GetAll(ctx context.Context, projectID domain.ProjectID) ([]domain.ProjectSetting, error)

	// Delete removes a project setting
	Delete(ctx context.Context, projectID domain.ProjectID, name string) error

	// Create creates a new project setting
	Create(ctx context.Context, setting *domain.ProjectSetting) error

	// GetByName retrieves a project setting by name
	GetByName(ctx context.Context, projectID domain.ProjectID, name string) (*domain.ProjectSetting, error)

	// Update updates a project setting
	Update(ctx context.Context, projectID domain.ProjectID, name string, value any) error

	// List retrieves project settings with pagination
	List(ctx context.Context, projectID domain.ProjectID, page, perPage int) ([]*domain.ProjectSetting, int, error)
}

// ProjectSettingsUseCase defines the interface for project settings operations.
type ProjectSettingsUseCase interface {
	Create(ctx context.Context, projectID domain.ProjectID, name string, value any) (*domain.ProjectSetting, error)
	GetByName(ctx context.Context, projectID domain.ProjectID, name string) (*domain.ProjectSetting, error)
	Update(ctx context.Context, projectID domain.ProjectID, name string, value any) (*domain.ProjectSetting, error)
	Delete(ctx context.Context, projectID domain.ProjectID, name string) error
	List(ctx context.Context, projectID domain.ProjectID, page, perPage int) ([]*domain.ProjectSetting, int, error)

	// Typed getters with safe type conversion and default values
	GetBoolSetting(ctx context.Context, projectID domain.ProjectID, name string, defaultValue bool) (bool, error)
	GetIntSetting(ctx context.Context, projectID domain.ProjectID, name string, defaultValue int) (int, error)
	GetStringSetting(ctx context.Context, projectID domain.ProjectID, name string, defaultValue string) (string, error)
	GetFloat64Setting(ctx context.Context, projectID domain.ProjectID, name string, defaultValue float64) (float64, error)

	// Auto-disable specific getters with predefined keys and defaults
	GetAutoDisableEnabled(ctx context.Context, projectID domain.ProjectID) (bool, error)
	GetAutoDisableRequiresApproval(ctx context.Context, projectID domain.ProjectID) (bool, error)
	GetAutoDisableErrorThreshold(ctx context.Context, projectID domain.ProjectID) (int, error)
	GetAutoDisableTimeWindowSec(ctx context.Context, projectID domain.ProjectID) (int, error)
	GetAutoDisableTimeWindow(ctx context.Context, projectID domain.ProjectID) (time.Duration, error)

	// Cached versions of auto-disable getters
	GetAutoDisableEnabledCached(ctx context.Context, projectID domain.ProjectID) (bool, error)
	GetAutoDisableRequiresApprovalCached(ctx context.Context, projectID domain.ProjectID) (bool, error)
	GetAutoDisableErrorThresholdCached(ctx context.Context, projectID domain.ProjectID) (int, error)
	GetAutoDisableTimeWindowSecCached(ctx context.Context, projectID domain.ProjectID) (int, error)
	GetAutoDisableTimeWindowCached(ctx context.Context, projectID domain.ProjectID) (time.Duration, error)
}
