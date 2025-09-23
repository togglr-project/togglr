package tags

import (
	"context"
	"fmt"
	"strings"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
)

type Service struct {
	tagRepo      contract.TagsRepository
	categoryRepo contract.CategoriesRepository
}

func New(tagRepo contract.TagsRepository, categoryRepo contract.CategoriesRepository) *Service {
	return &Service{
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *Service) CreateTag(
	ctx context.Context,
	projectID domain.ProjectID,
	categoryID *domain.CategoryID,
	name, slug string,
	description *string,
	color *string,
) (domain.Tag, error) {
	// Validate inputs
	if strings.TrimSpace(name) == "" {
		return domain.Tag{}, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(slug) == "" {
		return domain.Tag{}, fmt.Errorf("slug is required")
	}

	// Check if tag with this slug already exists in the project
	_, err := s.tagRepo.GetByProjectAndSlug(ctx, projectID, slug)
	if err == nil {
		return domain.Tag{}, fmt.Errorf("tag with slug %s already exists in project", slug)
	}
	if err != domain.ErrEntityNotFound {
		return domain.Tag{}, fmt.Errorf("check tag existence: %w", err)
	}

	// Validate category ID if provided
	if categoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *categoryID)
		if err != nil {
			return domain.Tag{}, fmt.Errorf("category %s not found: %w", *categoryID, err)
		}
	}

	// Create tag
	tagDTO := &domain.TagDTO{
		ProjectID:   projectID,
		CategoryID:  categoryID,
		Name:        strings.TrimSpace(name),
		Slug:        strings.TrimSpace(slug),
		Description: description,
		Color:       color,
	}

	id, err := s.tagRepo.Create(ctx, tagDTO)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("create tag: %w", err)
	}

	// Get created tag
	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get created tag: %w", err)
	}

	return tag, nil
}

func (s *Service) GetTag(ctx context.Context, id domain.TagID) (domain.Tag, error) {
	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get tag: %w", err)
	}

	return tag, nil
}

func (s *Service) ListProjectTags(
	ctx context.Context,
	projectID domain.ProjectID,
	categoryID *domain.CategoryID,
) ([]domain.Tag, error) {
	tags, err := s.tagRepo.ListByProject(ctx, projectID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("list project tags: %w", err)
	}

	return tags, nil
}

func (s *Service) UpdateTag(
	ctx context.Context,
	id domain.TagID,
	categoryID *domain.CategoryID,
	name, slug string,
	description *string,
	color *string,
) (domain.Tag, error) {
	// Validate inputs
	if strings.TrimSpace(name) == "" {
		return domain.Tag{}, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(slug) == "" {
		return domain.Tag{}, fmt.Errorf("slug is required")
	}

	// Check if tag exists
	existingTag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get tag: %w", err)
	}

	// Check if slug is already taken by another tag in the same project
	if existingTag.Slug != slug {
		_, err := s.tagRepo.GetByProjectAndSlug(ctx, existingTag.ProjectID, slug)
		if err == nil {
			return domain.Tag{}, fmt.Errorf("tag with slug %s already exists in project", slug)
		}
		if err != domain.ErrEntityNotFound {
			return domain.Tag{}, fmt.Errorf("check tag existence: %w", err)
		}
	}

	// Validate category ID if provided
	if categoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *categoryID)
		if err != nil {
			return domain.Tag{}, fmt.Errorf("category %s not found: %w", *categoryID, err)
		}
	}

	// Update tag
	err = s.tagRepo.Update(ctx, id, categoryID, strings.TrimSpace(name), strings.TrimSpace(slug), description, color)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("update tag: %w", err)
	}

	// Get updated tag
	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get updated tag: %w", err)
	}

	return tag, nil
}

func (s *Service) DeleteTag(ctx context.Context, id domain.TagID) error {
	// Check if tag exists
	_, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

	// Delete tag
	err = s.tagRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}

	return nil
}

func (s *Service) CreateTagsFromCategories(ctx context.Context, projectID domain.ProjectID) error {
	err := s.tagRepo.CreateFromCategories(ctx, projectID)
	if err != nil {
		return fmt.Errorf("create tags from categories: %w", err)
	}

	return nil
}
