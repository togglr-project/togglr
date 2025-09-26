package project_settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
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
