package email

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type UsersRepository interface {
	FetchByIDs(ctx context.Context, ids []domain.UserID) ([]domain.User, error)
}

type ProjectsRepository interface {
	GetByID(ctx context.Context, id domain.ProjectID) (domain.Project, error)
}
