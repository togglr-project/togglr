package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureContextualStatsRepository interface {
	LoadAll(ctx context.Context) ([]domain.FeatureContextualStats, error)
	InsertBatch(ctx context.Context, records []domain.FeatureContextualStats) error
}
