package categories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
	categoryRepo contract.CategoriesRepository
}

func New(categoryRepo contract.CategoriesRepository) *Service {
	return &Service{
		categoryRepo: categoryRepo,
	}
}

func (s *Service) CreateCategory(
	ctx context.Context,
	name, slug string,
	description *string,
	color *string,
	categoryType domain.CategoryType,
) (domain.Category, error) {
	// Validate inputs
	if strings.TrimSpace(name) == "" {
		return domain.Category{}, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(slug) == "" {
		return domain.Category{}, fmt.Errorf("slug is required")
	}

	// Check if the category with this slug already exists
	_, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err == nil {
		return domain.Category{}, fmt.Errorf("category with slug %s already exists", slug)
	}
	if !errors.Is(err, domain.ErrEntityNotFound) {
		return domain.Category{}, fmt.Errorf("check category existence: %w", err)
	}

	// Create category
	categoryDTO := &domain.CategoryDTO{
		Name:        strings.TrimSpace(name),
		Slug:        strings.TrimSpace(slug),
		Description: description,
		Color:       color,
		Kind:        domain.CategoryKindUser,
		Type:        categoryType,
	}

	id, err := s.categoryRepo.Create(ctx, categoryDTO)
	if err != nil {
		return domain.Category{}, fmt.Errorf("create category: %w", err)
	}

	// Get created category
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("get created category: %w", err)
	}

	return category, nil
}

func (s *Service) GetCategory(ctx context.Context, id domain.CategoryID) (domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("get category: %w", err)
	}

	return category, nil
}

func (s *Service) ListCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return categories, nil
}

func (s *Service) UpdateCategory(
	ctx context.Context,
	id domain.CategoryID,
	name, slug string,
	description *string,
	color *string,
) (domain.Category, error) {
	// Validate inputs
	if strings.TrimSpace(name) == "" {
		return domain.Category{}, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(slug) == "" {
		return domain.Category{}, fmt.Errorf("slug is required")
	}

	// Check if category exists
	existingCategory, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("get category: %w", err)
	}

	// Check if slug is already taken by another category
	if existingCategory.Slug != slug {
		_, err := s.categoryRepo.GetBySlug(ctx, slug)
		if err == nil {
			return domain.Category{}, fmt.Errorf("category with slug %s already exists", slug)
		}
		if !errors.Is(err, domain.ErrEntityNotFound) {
			return domain.Category{}, fmt.Errorf("check category existence: %w", err)
		}
	}

	// Update category
	err = s.categoryRepo.Update(ctx, id, strings.TrimSpace(name), strings.TrimSpace(slug), description, color)
	if err != nil {
		return domain.Category{}, fmt.Errorf("update category: %w", err)
	}

	// Get an updated category
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("get updated category: %w", err)
	}

	return category, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id domain.CategoryID) error {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get category: %w", err)
	}

	// Delete category
	err = s.categoryRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
