package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type ProjectsUseCase interface {
	CreateProject(ctx context.Context, name, description string) (domain.Project, error)
	GetProject(ctx context.Context, id domain.ProjectID) (domain.Project, error)
	List(ctx context.Context) ([]domain.Project, error)
	UpdateInfo(ctx context.Context, id domain.ProjectID, name, description string) (domain.Project, error)
	ArchiveProject(ctx context.Context, id domain.ProjectID) error
	ListChanges(ctx context.Context, filter domain.ChangesListFilter) (domain.ChangesListResult, error)
	// GetAPIKeyForEnvironment returns the API key for a specific environment
	GetAPIKeyForEnvironment(ctx context.Context, projectID domain.ProjectID, environmentID domain.EnvironmentID) (string, error)
}

type ProjectsRepository interface {
	GetByID(ctx context.Context, id domain.ProjectID) (domain.Project, error)
	GetByAPIKey(ctx context.Context, apiKey string) (domain.ProjectWithEnv, error)
	Create(ctx context.Context, project *domain.ProjectDTO) (domain.ProjectID, error)
	List(ctx context.Context) ([]domain.Project, error)
	Update(ctx context.Context, id domain.ProjectID, name, description string) error
	Archive(ctx context.Context, id domain.ProjectID) error
}
