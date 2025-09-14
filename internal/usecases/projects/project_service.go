package projects

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
)

type ProjectService struct {
	projectRepo contract.ProjectsRepository
}

func New(
	projectRepo contract.ProjectsRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

func (s *ProjectService) GetProject(ctx context.Context, id domain.ProjectID) (domain.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

func (s *ProjectService) CreateProject(
	ctx context.Context,
	name, description string,
) (domain.Project, error) {
	project := domain.ProjectDTO{
		Name:        name,
		Description: description,
		APIKey:      uuid.NewString(),
	}

	id, err := s.projectRepo.Create(ctx, &project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("create project: %w", err)
	}

	return domain.Project{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *ProjectService) List(ctx context.Context) ([]domain.Project, error) {
	return s.projectRepo.List(ctx)
}

func (s *ProjectService) UpdateInfo(
	ctx context.Context,
	id domain.ProjectID,
	name, description string,
) (domain.Project, error) {
	// Check if the project exists
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to get project: %w", err)
	}

	// Update the project
	err = s.projectRepo.Update(ctx, id, name, description)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to update project: %w", err)
	}

	// Return the updated project with extended info
	project.Name = name
	project.Description = description

	return project, nil
}

func (s *ProjectService) ArchiveProject(ctx context.Context, id domain.ProjectID) error {
	// Check if the project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Archive the project
	err = s.projectRepo.Archive(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to archive project: %w", err)
	}

	slog.Info("project archived", "project_id", id)

	return nil
}
