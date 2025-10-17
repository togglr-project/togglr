package environments

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	simplecache "github.com/togglr-project/togglr/pkg/simple-cache"
)

const (
	CacheTTL = 30 * time.Minute
)

type Service struct {
	envRepo contract.EnvironmentsRepository
	cache   *simplecache.Cache[string, domain.Environment]
}

func New(envRepo contract.EnvironmentsRepository) *Service {
	service := &Service{
		envRepo: envRepo,
		cache:   simplecache.New[string, domain.Environment](),
	}
	service.cache.StartCleanup(5 * time.Minute)

	return service
}

func (s *Service) Create(
	ctx context.Context,
	projectID domain.ProjectID,
	key, name string,
) (domain.Environment, error) {
	apiKey := uuid.New().String()

	env := domain.Environment{
		ProjectID: projectID,
		Key:       key,
		Name:      name,
		APIKey:    apiKey,
		CreatedAt: time.Now(),
	}

	return s.envRepo.Create(ctx, env)
}

func (s *Service) GetByID(ctx context.Context, id domain.EnvironmentID) (domain.Environment, error) {
	return s.envRepo.GetByID(ctx, id)
}

func (s *Service) GetByProjectIDAndKey(
	ctx context.Context,
	projectID domain.ProjectID,
	key string,
) (domain.Environment, error) {
	return s.envRepo.GetByProjectIDAndKey(ctx, projectID, key)
}

func (s *Service) ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Environment, error) {
	return s.envRepo.ListByProjectID(ctx, projectID)
}

func (s *Service) Update(ctx context.Context, id domain.EnvironmentID, name string) (domain.Environment, error) {
	env, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Environment{}, fmt.Errorf("get environment: %w", err)
	}

	env.Name = name

	return s.envRepo.Update(ctx, env)
}

func (s *Service) Delete(ctx context.Context, id domain.EnvironmentID) error {
	return s.envRepo.Delete(ctx, id)
}

func (s *Service) GetByIDCached(ctx context.Context, id domain.EnvironmentID) (domain.Environment, error) {
	cacheKey := makeEnvIDCacheKey(id)

	if cached, found := s.cache.Get(cacheKey); found {
		return cached, nil
	}

	env, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Environment{}, err
	}

	s.cache.Set(cacheKey, env, CacheTTL)

	return env, nil
}

func (s *Service) GetByProjectIDAndKeyCached(
	ctx context.Context,
	projectID domain.ProjectID,
	key string,
) (domain.Environment, error) {
	cacheKey := makeEnvironmentCacheKey(projectID, key)

	if cached, found := s.cache.Get(cacheKey); found {
		return cached, nil
	}

	env, err := s.envRepo.GetByProjectIDAndKey(ctx, projectID, key)
	if err != nil {
		return domain.Environment{}, err
	}

	s.cache.Set(cacheKey, env, CacheTTL)

	return env, nil
}

// func (s *Service) InvalidateCache(projectID domain.ProjectID, envKey string) {
//	cacheKey := makeEnvironmentCacheKey(projectID, envKey)
//	s.cache.Delete(cacheKey)
//}
//
// func (s *Service) InvalidateProjectCache(projectID domain.ProjectID) {
//	s.cache.Clear()
//}

func makeEnvironmentCacheKey(projectID domain.ProjectID, envKey string) string {
	return string(projectID) + ":" + envKey
}

func makeEnvIDCacheKey(envID domain.EnvironmentID) string {
	return envID.String()
}
