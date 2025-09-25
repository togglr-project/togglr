package dto

import (
	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainCategoryToAPI converts domain Category to generated API Category.
func DomainCategoryToAPI(category domain.Category) generatedapi.Category {
	item := generatedapi.Category{
		ID:        uuid.MustParse(category.ID.String()),
		Name:      category.Name,
		Slug:      category.Slug,
		Kind:      generatedapi.CategoryKind(category.Kind),
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if category.Description != nil {
		item.Description = generatedapi.NewOptNilString(*category.Description)
	}

	if category.Color != nil {
		item.Color = generatedapi.NewOptNilString(*category.Color)
	}

	return item
}

// DomainCategoriesToAPI converts slice of domain Categories to slice of generated API Categories.
func DomainCategoriesToAPI(categories []domain.Category) []generatedapi.Category {
	resp := make([]generatedapi.Category, 0, len(categories))
	for _, category := range categories {
		resp = append(resp, DomainCategoryToAPI(category))
	}

	return resp
}

// APICategoryToDomain converts generated API Category to domain Category.
func APICategoryToDomain(category generatedapi.Category) domain.Category {
	item := domain.Category{
		ID:        domain.CategoryID(category.ID.String()),
		Name:      category.Name,
		Slug:      category.Slug,
		Kind:      domain.CategoryKind(category.Kind),
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if category.Description.IsSet() {
		item.Description = &category.Description.Value
	}

	if category.Color.IsSet() {
		item.Color = &category.Color.Value
	}

	return item
}
