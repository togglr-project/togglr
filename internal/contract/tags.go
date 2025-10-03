package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type TagsUseCase interface {
	CreateTag(
		ctx context.Context,
		projectID domain.ProjectID,
		categoryID *domain.CategoryID,
		name, slug string,
		description *string,
		color *string,
	) (domain.Tag, error)
	GetTag(ctx context.Context, id domain.TagID) (domain.Tag, error)
	ListProjectTags(
		ctx context.Context,
		projectID domain.ProjectID,
		categoryID *domain.CategoryID,
	) ([]domain.Tag, error)
	UpdateTag(
		ctx context.Context,
		id domain.TagID,
		categoryID *domain.CategoryID,
		name, slug string,
		description *string,
		color *string,
	) (domain.Tag, error)
	DeleteTag(ctx context.Context, id domain.TagID) error
	CreateTagsFromCategories(
		ctx context.Context,
		projectID domain.ProjectID,
	) error

	// System tag specific methods
	GetAutoDisableTag(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error)
	GetGuardedTag(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error)
	EnsureSystemTags(ctx context.Context, projectID domain.ProjectID) error
}

type TagsRepository interface {
	GetByID(ctx context.Context, id domain.TagID) (domain.Tag, error)
	GetByProjectAndSlug(
		ctx context.Context,
		projectID domain.ProjectID,
		slug string,
	) (domain.Tag, error)
	ListByProject(
		ctx context.Context,
		projectID domain.ProjectID,
		categoryID *domain.CategoryID,
	) ([]domain.Tag, error)
	Create(ctx context.Context, tag *domain.TagDTO) (domain.TagID, error)
	Update(
		ctx context.Context,
		id domain.TagID,
		categoryID *domain.CategoryID,
		name, slug string,
		description *string,
		color *string,
	) error
	Delete(ctx context.Context, id domain.TagID) error
	CreateFromCategories(ctx context.Context, projectID domain.ProjectID) error
}
