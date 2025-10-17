package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureAlgorithmStatsRepository interface {
	LoadAll(ctx context.Context) ([]domain.FeatureAlgorithmStats, error)
	LoadByFeatureEnvAlg(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
		algSlug string,
	) (domain.FeatureAlgorithmStats, error)
	InsertBatch(ctx context.Context, records []domain.FeatureAlgorithmStats) error
}
