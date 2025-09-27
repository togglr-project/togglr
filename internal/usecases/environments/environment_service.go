package environments

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
	envRepo contract.EnvironmentsRepository
}

func New(envRepo contract.EnvironmentsRepository) *Service {
	return &Service{
		envRepo: envRepo,
	}
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
