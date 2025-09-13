package permissions

import (
	"context"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
)

// Service handles permission checks for various operations.
type Service struct {
	projectRepo contract.ProjectsRepository
}

// New creates a new permissions service.
func New(
	projectRepo contract.ProjectsRepository,
) *Service {
	return &Service{
		projectRepo: projectRepo,
	}
}

// CanAccessProject checks if a user can access a project.
func (s *Service) CanAccessProject(ctx context.Context, projectID domain.ProjectID) error {
	// Get the project to check its team
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	isSuper := etogglcontext.IsSuper(ctx)
	if isSuper {
		return nil
	}

	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// CanManageProject checks if a user can manage a project (create, update, delete).
func (s *Service) CanManageProject(ctx context.Context, projectID domain.ProjectID) error {
	// Get the project to check its team
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	isSuper := etogglcontext.IsSuper(ctx)
	if isSuper {
		return nil
	}

	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// GetAccessibleProjects returns all projects that a user can access.
func (s *Service) GetAccessibleProjects(
	ctx context.Context,
	projects []domain.Project,
) ([]domain.Project, error) {
	return projects, nil
}
