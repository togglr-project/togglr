package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureAlgorithmsRepository interface {
	Create(
		ctx context.Context,
		featureAlgorithm domain.FeatureAlgorithmDTO,
	) error
	Update(
		ctx context.Context,
		featureAlgorithm domain.FeatureAlgorithmDTO,
	) error
	Delete(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) error
	ListByFeatureID(
		ctx context.Context,
		featureID domain.FeatureID,
	) ([]domain.FeatureAlgorithm, error)
	ListEnabled(ctx context.Context) ([]domain.FeatureAlgorithm, error)
	ListAll(ctx context.Context) ([]domain.FeatureAlgorithm, error)
}
