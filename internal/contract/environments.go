package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type EnvironmentsUseCase interface {
	Create(ctx context.Context, projectID domain.ProjectID, key, name string) (domain.Environment, error)
	GetByID(ctx context.Context, id domain.EnvironmentID) (domain.Environment, error)
	GetByProjectIDAndKey(ctx context.Context, projectID domain.ProjectID, key string) (domain.Environment, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Environment, error)
	Update(ctx context.Context, id domain.EnvironmentID, name string) (domain.Environment, error)
	Delete(ctx context.Context, id domain.EnvironmentID) error

	// Cached version
	GetByProjectIDAndKeyCached(ctx context.Context, projectID domain.ProjectID, key string) (domain.Environment, error)
}

type EnvironmentsRepository interface {
	Create(ctx context.Context, env domain.Environment) (domain.Environment, error)
	GetByID(ctx context.Context, id domain.EnvironmentID) (domain.Environment, error)
	GetByProjectIDAndKey(ctx context.Context, projectID domain.ProjectID, key string) (domain.Environment, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Environment, error)
	Update(ctx context.Context, env domain.Environment) (domain.Environment, error)
	Delete(ctx context.Context, id domain.EnvironmentID) error
}
