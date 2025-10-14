package feature_algorithm_stats

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) LoadAll(ctx context.Context) ([]domain.FeatureAlgorithmStats, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM monitoring.feature_algorithm_stats`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureAlgorithmStatsModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_algorithms rows: %w", err)
	}

	result := make([]domain.FeatureAlgorithmStats, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) LoadByFeatureEnvAlg(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	algSlug string,
) (domain.FeatureAlgorithmStats, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT * FROM monitoring.feature_algorithm_stats
WHERE feature_id = $1 AND environment_id = $2 AND algorithm_slug = $3`

	rows, err := executor.Query(ctx, query, featureID, envID, algSlug)
	if err != nil {
		return domain.FeatureAlgorithmStats{}, err
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureAlgorithmStatsModel])
	if err != nil {
		return domain.FeatureAlgorithmStats{}, fmt.Errorf("collect row: %w", err)
	}

	return model.toDomain(), nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
