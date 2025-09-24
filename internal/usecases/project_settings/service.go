package project_settings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// Service provides project settings management functionality.
type Service struct {
	projectSettingsRepo contract.ProjectSettingsRepository
}

// New creates a new project settings use case.
func New(projectSettingsRepo contract.ProjectSettingsRepository) *Service {
	return &Service{
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
		return nil, fmt.Errorf("setting name cannot be empty")
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

	err = s.projectSettingsRepo.Create(ctx, setting)
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
		return nil, fmt.Errorf("setting name cannot be empty")
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
		return nil, fmt.Errorf("setting name cannot be empty")
	}

	// Marshal value to JSON to validate it
	_, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("invalid setting value: %w", err)
	}

	err = s.projectSettingsRepo.Update(ctx, projectID, name, value)
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
		return fmt.Errorf("setting name cannot be empty")
	}

	err := s.projectSettingsRepo.Delete(ctx, projectID, name)
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
