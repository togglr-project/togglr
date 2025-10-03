package tags

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
	simplecache "github.com/togglr-project/togglr/pkg/simple-cache"
)

const (
	AutoDisableTagSlug = "auto-disable"
	GuardedTagSlug     = "guarded"
	CacheTTL           = 10 * time.Minute
)

type Service struct {
	txManager    db.TxManager
	tagRepo      contract.TagsRepository
	categoryRepo contract.CategoriesRepository
	cache        *simplecache.Cache[string, domain.Tag]
}

func New(
	txManager db.TxManager,
	tagRepo contract.TagsRepository,
	categoryRepo contract.CategoriesRepository,
) *Service {
	service := &Service{
		txManager:    txManager,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
		cache:        simplecache.New[string, domain.Tag](),
	}
	service.cache.StartCleanup(2 * time.Minute)

	return service
}

func (s *Service) CreateTag(
	ctx context.Context,
	projectID domain.ProjectID,
	categoryID *domain.CategoryID,
	name, slug string,
	description *string,
	color *string,
) (domain.Tag, error) {
	if strings.TrimSpace(name) == "" {
		return domain.Tag{}, errors.New("name is required")
	}

	if strings.TrimSpace(slug) == "" {
		return domain.Tag{}, errors.New("slug is required")
	}

	_, err := s.tagRepo.GetByProjectAndSlug(ctx, projectID, slug)
	if err == nil {
		return domain.Tag{}, fmt.Errorf("tag with slug %s already exists in project", slug)
	}

	if !errors.Is(err, domain.ErrEntityNotFound) {
		return domain.Tag{}, fmt.Errorf("check tag existence: %w", err)
	}

	if categoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *categoryID)
		if err != nil {
			return domain.Tag{}, fmt.Errorf("category %s not found: %w", *categoryID, err)
		}
	}

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
	if strings.TrimSpace(name) == "" {
		return domain.Tag{}, errors.New("name is required")
	}

	if strings.TrimSpace(slug) == "" {
		return domain.Tag{}, errors.New("slug is required")
	}

	existingTag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get tag: %w", err)
	}

	if existingTag.Slug != slug {
		_, err := s.tagRepo.GetByProjectAndSlug(ctx, existingTag.ProjectID, slug)
		if err == nil {
			return domain.Tag{}, fmt.Errorf("tag with slug %s already exists in project", slug)
		}

		if !errors.Is(err, domain.ErrEntityNotFound) {
			return domain.Tag{}, fmt.Errorf("check tag existence: %w", err)
		}
	}

	if categoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *categoryID)
		if err != nil {
			return domain.Tag{}, fmt.Errorf("category %s not found: %w", *categoryID, err)
		}
	}

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

	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get updated tag: %w", err)
	}

	return tag, nil
}

func (s *Service) DeleteTag(ctx context.Context, id domain.TagID) error {
	_, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

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

func (s *Service) GetAutoDisableTag(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error) {
	tag, err := s.tagRepo.GetByProjectAndSlug(ctx, projectID, AutoDisableTagSlug)
	if err == nil {
		return tag, nil
	}

	if !errors.Is(err, domain.ErrEntityNotFound) {
		return domain.Tag{}, fmt.Errorf("get auto-disable tag: %w", err)
	}

	return s.createSystemTag(ctx, projectID, AutoDisableTagSlug, "Auto Disable", "Enables automatic feature disabling on errors")
}

func (s *Service) GetGuardedTag(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error) {
	tag, err := s.tagRepo.GetByProjectAndSlug(ctx, projectID, GuardedTagSlug)
	if err == nil {
		return tag, nil
	}

	if !errors.Is(err, domain.ErrEntityNotFound) {
		return domain.Tag{}, fmt.Errorf("get guarded tag: %w", err)
	}

	return s.createSystemTag(ctx, projectID, GuardedTagSlug, "Guarded", "Requires approval for changes")
}

func (s *Service) EnsureSystemTags(ctx context.Context, projectID domain.ProjectID) error {
	_, err := s.GetAutoDisableTag(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ensure auto-disable tag: %w", err)
	}

	_, err = s.GetGuardedTag(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ensure guarded tag: %w", err)
	}

	return nil
}

func (s *Service) GetAutoDisableTagCached(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error) {
	cacheKey := makeTagCacheKey(projectID, AutoDisableTagSlug)

	if cached, found := s.cache.Get(cacheKey); found {
		return cached, nil
	}

	tag, err := s.GetAutoDisableTag(ctx, projectID)
	if err != nil {
		return domain.Tag{}, err
	}

	s.cache.Set(cacheKey, tag, CacheTTL)

	return tag, nil
}

func (s *Service) GetGuardedTagCached(ctx context.Context, projectID domain.ProjectID) (domain.Tag, error) {
	cacheKey := makeTagCacheKey(projectID, GuardedTagSlug)

	if cached, found := s.cache.Get(cacheKey); found {
		return cached, nil
	}

	tag, err := s.GetGuardedTag(ctx, projectID)
	if err != nil {
		return domain.Tag{}, err
	}

	s.cache.Set(cacheKey, tag, CacheTTL)

	return tag, nil
}

func (s *Service) InvalidateCache(projectID domain.ProjectID) {
	keys := []string{
		makeTagCacheKey(projectID, AutoDisableTagSlug),
		makeTagCacheKey(projectID, GuardedTagSlug),
	}

	for _, key := range keys {
		s.cache.Delete(key)
	}
}

func (s *Service) createSystemTag(ctx context.Context, projectID domain.ProjectID, slug, name, description string) (domain.Tag, error) {
	descriptionPtr := &description
	color := "#ff6b6b"

	var id domain.TagID
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		id, err = s.tagRepo.Create(ctx, &domain.TagDTO{
			ProjectID:   projectID,
			CategoryID:  nil,
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

	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("get created system tag %s: %w", slug, err)
	}

	return tag, nil
}

func makeTagCacheKey(projectID domain.ProjectID, tagSlug string) string {
	return string(projectID) + ":" + tagSlug
}
