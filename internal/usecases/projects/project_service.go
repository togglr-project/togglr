package projects

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type ProjectService struct {
	projectRepo  contract.ProjectsRepository
	auditLogRepo contract.AuditLogRepository
	tagsUseCase  contract.TagsUseCase
}

func New(
	projectRepo contract.ProjectsRepository,
	auditLogRepo contract.AuditLogRepository,
	tagsUseCase contract.TagsUseCase,
) *ProjectService {
	return &ProjectService{
		projectRepo:  projectRepo,
		auditLogRepo: auditLogRepo,
		tagsUseCase:  tagsUseCase,
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

	// Create tags from system categories
	// err = s.tagsUseCase.CreateTagsFromCategories(ctx, id)
	// if err != nil {
	//	slog.Error("failed to create tags from categories", "error", err, "project_id", id)
	//	// Don't fail project creation if tag creation fails
	//}

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

func (s *ProjectService) ListChanges(
	ctx context.Context,
	filter domain.ChangesListFilter,
) (domain.ChangesListResult, error) {
	// Verify that the project exists
	_, err := s.projectRepo.GetByID(ctx, filter.ProjectID)
	if err != nil {
		return domain.ChangesListResult{}, fmt.Errorf("get project: %w", err)
	}

	// Get changes from audit log repository
	result, err := s.auditLogRepo.ListChanges(ctx, filter)
	if err != nil {
		return domain.ChangesListResult{}, fmt.Errorf("list changes: %w", err)
	}

	return result, nil
}
