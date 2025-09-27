package feature_params

import (
	"context"
	"fmt"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
	paramsRepo contract.FeatureParamsRepository
}

func New(paramsRepo contract.FeatureParamsRepository) *Service {
	return &Service{
		paramsRepo: paramsRepo,
	}
}

func (s *Service) GetByFeatureAndEnvironment(
	ctx context.Context,
	featureID domain.FeatureID,
	environmentID domain.EnvironmentID,
) (domain.FeatureParams, error) {
	return s.paramsRepo.GetByFeatureAndEnvironment(ctx, featureID, environmentID)
}

func (s *Service) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error) {
	return s.paramsRepo.ListByFeatureID(ctx, featureID)
}

func (s *Service) Update(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error) {
	params.UpdatedAt = time.Now()

	existing, err := s.paramsRepo.GetByFeatureAndEnvironment(ctx, params.FeatureID, params.EnvironmentID)
	if err != nil {
		if err == domain.ErrEntityNotFound {
			return s.paramsRepo.Create(ctx, projectID, params)
		}
		return domain.FeatureParams{}, fmt.Errorf("get existing params: %w", err)
	}

	existing.Enabled = params.Enabled
	existing.DefaultValue = params.DefaultValue
	existing.UpdatedAt = params.UpdatedAt

	return s.paramsRepo.Update(ctx, projectID, existing)
}

func (s *Service) Delete(
	ctx context.Context,
	projectID domain.ProjectID,
	featureID domain.FeatureID,
	environmentID domain.EnvironmentID,
) error {
	return s.paramsRepo.Delete(ctx, projectID, featureID, environmentID)
}
