package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type AlgorithmsRepository interface {
	GetBySlug(ctx context.Context, slug string) (domain.Algorithm, error)
	List(ctx context.Context) ([]domain.Algorithm, error)
}
