package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureOptimizerStatsRepository interface {
	LoadAll(ctx context.Context) ([]domain.FeatureOptimizerStats, error)
	InsertBatch(ctx context.Context, records []domain.FeatureOptimizerStats) error
}
