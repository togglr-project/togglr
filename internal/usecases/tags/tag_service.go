package tags

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

const (
	// System tag slugs for auto-disable functionality
	AutoDisableTagSlug = "auto-disable"
	GuardedTagSlug     = "guarded"
)

type Service struct {
	txManager    db.TxManager
	tagRepo      contract.TagsRepository
	categoryRepo contract.CategoriesRepository
}

func New(
	txManager db.TxManager,
	tagRepo contract.TagsRepository,
	categoryRepo contract.CategoriesRepository,
) *Service {
	return &Service{
		txManager:    txManager,
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
		return domain.Tag{}, errors.New("name is required")
	}

	if strings.TrimSpace(slug) == "" {
		return domain.Tag{}, errors.New("slug is required")
	}

	// Check if tag with this slug already exists in the project
	_, err := s.tagRepo.GetByProjectAndSlug(ctx, projectID, slug)
	if err == nil {
		return domain.Tag{}, fmt.Errorf("tag with slug %s already exists in project", slug)
	}

	if !errors.Is(err, domain.ErrEntityNotFound) {
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

	var id domain.TagID
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		id, err = s.tagRepo.Create(ctx, tagDTO)

		return err
	})
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
		return domain.Tag{}, errors.New("name is required")
	}

	if strings.TrimSpace(slug) == "" {
		return domain.Tag{}, errors.New("slug is required")
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

		if !errors.Is(err, domain.ErrEntityNotFound) {
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
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.tagRepo.Update(
			ctx,
			id,
			categoryID,
			strings.TrimSpace(name),
			strings.TrimSpace(slug),
			description,
			color,
		)
	})
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
	// Check if a tag exists
	_, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

	// Delete tag
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.tagRepo.Delete(ctx, id)
	})
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

// System tag specific methods

// GetAutoDisableTag retrieves the auto-disable tag for a project, creating it if it doesn't exist.
func (s *Service) GetAutoDisableTag(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error) {
	tag, err := s.tagRepo.GetByProjectAndSlug(ctx, projectID, AutoDisableTagSlug)
	if err == nil {
		return tag, nil
	}

	if !errors.Is(err, domain.ErrEntityNotFound) {
		return domain.Tag{}, fmt.Errorf("get auto-disable tag: %w", err)
	}

	// Create auto-disable tag if it doesn't exist
	return s.createSystemTag(ctx, projectID, AutoDisableTagSlug, "Auto Disable", "Enables automatic feature disabling on errors")
}

// GetGuardedTag retrieves the guarded tag for a project, creating it if it doesn't exist.
func (s *Service) GetGuardedTag(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error) {
	tag, err := s.tagRepo.GetByProjectAndSlug(ctx, projectID, GuardedTagSlug)
	if err == nil {
		return tag, nil
	}

	if !errors.Is(err, domain.ErrEntityNotFound) {
		return domain.Tag{}, fmt.Errorf("get guarded tag: %w", err)
	}

	// Create guarded tag if it doesn't exist
	return s.createSystemTag(ctx, projectID, GuardedTagSlug, "Guarded", "Requires approval for changes")
}

// EnsureSystemTags ensures that all system tags exist for a project.
func (s *Service) EnsureSystemTags(ctx context.Context, projectID domain.ProjectID) error {
	// Ensure auto-disable tag exists
	_, err := s.GetAutoDisableTag(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ensure auto-disable tag: %w", err)
	}

	// Ensure guarded tag exists
	_, err = s.GetGuardedTag(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ensure guarded tag: %w", err)
	}

	return nil
}

// createSystemTag creates a system tag with predefined properties.
func (s *Service) createSystemTag(ctx context.Context, projectID domain.ProjectID, slug, name, description string) (domain.Tag, error) {
	descriptionPtr := &description
	color := "#ff6b6b" // Default red color for system tags

	var id domain.TagID
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		id, err = s.tagRepo.Create(ctx, &domain.TagDTO{
			ProjectID:   projectID,
			CategoryID:  nil, // System tags don't belong to categories
			Name:        name,
			Slug:        slug,
			Description: descriptionPtr,
			Color:       &color,
		})
		return err
	})
	if err != nil {
		return domain.Tag{}, fmt.Errorf("create system tag %s: %w", slug, err)
	}

	// Get created tag
	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get created system tag %s: %w", slug, err)
	}

	return tag, nil
}
