package feature_algorithms

import (
	"context"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

var _ contract.FeatureAlgorithmsUseCase = (*Service)(nil)

type Service struct {
	txManager db.TxManager
	repo      contract.FeatureAlgorithmsRepository
}

func New(
	txManager db.TxManager,
	repo contract.FeatureAlgorithmsRepository,
) *Service {
	return &Service{
		txManager: txManager,
		repo:      repo,
	}
}

func (s Service) Create(ctx context.Context, featureAlgorithm domain.FeatureAlgorithmDTO) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.repo.Create(ctx, featureAlgorithm)
	})
}

func (s Service) Update(ctx context.Context, featureAlgorithm domain.FeatureAlgorithm) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.repo.Update(ctx, featureAlgorithm)
	})
}

func (s Service) DeleteByFeatureIDWithEnvID(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.repo.DeleteByFeatureIDWithEnvID(ctx, featureID, envID)
	})
}

func (s Service) GetByFeatureIDWithEnvID(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) (domain.FeatureAlgorithm, error) {
	return s.repo.GetByFeatureIDWithEnvID(ctx, featureID, envID)
}
