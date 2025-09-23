package dto

import (
	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainTagToAPI converts domain Tag to generated API ProjectTag
func DomainTagToAPI(tag domain.Tag) generatedapi.ProjectTag {
	item := generatedapi.ProjectTag{
		ID:        uuid.MustParse(tag.ID.String()),
		ProjectID: uuid.MustParse(tag.ProjectID.String()),
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}

	if tag.Description != nil {
		item.Description = generatedapi.NewOptNilString(*tag.Description)
	}
	if tag.Color != nil {
		item.Color = generatedapi.NewOptNilString(*tag.Color)
	}

	// Convert category if present
	if tag.Category != nil {
		catItem := DomainCategoryToAPI(*tag.Category)
		item.Category = generatedapi.NewOptCategory(catItem)
	}

	// Convert category ID if present
	if tag.CategoryID != nil {
		catID, err := uuid.Parse(tag.CategoryID.String())
		if err == nil {
			item.CategoryID = generatedapi.NewOptNilUUID(catID)
		}
	}

	return item
}

// DomainTagsToAPI converts slice of domain Tags to slice of generated API ProjectTags
func DomainTagsToAPI(tags []domain.Tag) []generatedapi.ProjectTag {
	resp := make([]generatedapi.ProjectTag, 0, len(tags))
	for _, tag := range tags {
		resp = append(resp, DomainTagToAPI(tag))
	}
	return resp
}

// APITagToDomain converts generated API ProjectTag to domain Tag
func APITagToDomain(tag generatedapi.ProjectTag) domain.Tag {
	item := domain.Tag{
		ID:        domain.TagID(tag.ID.String()),
		ProjectID: domain.ProjectID(tag.ProjectID.String()),
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}

	if tag.Description.IsSet() {
		item.Description = &tag.Description.Value
	}
	if tag.Color.IsSet() {
		item.Color = &tag.Color.Value
	}

	// Convert category if present
	if tag.Category.IsSet() {
		category := APICategoryToDomain(tag.Category.Value)
		item.Category = &category
	}

	// Convert category ID if present
	if tag.CategoryID.IsSet() {
		categoryID := domain.CategoryID(tag.CategoryID.Value.String())
		item.CategoryID = &categoryID
	}

	return item
}
