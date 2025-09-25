package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type CategoriesUseCase interface {
	CreateCategory(
		ctx context.Context,
		name, slug string,
		description *string,
		color *string,
	) (domain.Category, error)
	GetCategory(ctx context.Context, id domain.CategoryID) (domain.Category, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(
		ctx context.Context,
		id domain.CategoryID,
		name, slug string,
		description *string,
		color *string,
	) (domain.Category, error)
	DeleteCategory(ctx context.Context, id domain.CategoryID) error
}

type CategoriesRepository interface {
	GetByID(ctx context.Context, id domain.CategoryID) (domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	Create(ctx context.Context, category *domain.CategoryDTO) (domain.CategoryID, error)
	Update(
		ctx context.Context,
		id domain.CategoryID,
		name, slug string,
		description *string,
		color *string,
	) error
	Delete(ctx context.Context, id domain.CategoryID) error
}
