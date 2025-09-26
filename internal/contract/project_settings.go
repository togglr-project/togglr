package contract

import (
	"context"

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
}
