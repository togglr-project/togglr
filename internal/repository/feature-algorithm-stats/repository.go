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

func (r *Repository) InsertBatch(ctx context.Context, records []domain.FeatureAlgorithmStats) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO monitoring.feature_algorithm_stats
(feature_id, environment_id, algorithm_slug, variant_key,
	evaluations, successes, failures, metric_sum, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())
ON CONFLICT (feature_id, environment_id, algorithm_slug, variant_key) DO UPDATE SET
	evaluations = EXCLUDED.evaluations,
	successes   = EXCLUDED.successes,
	failures    = EXCLUDED.failures,
	metric_sum  = EXCLUDED.metric_sum,
	updated_at  = NOW();`

	batch := &pgx.Batch{}
	for _, record := range records {
		// skip zeros
		if record.Evaluations == 0 && record.Successes == 0 && record.Failures == 0 && record.MetricSum.IsZero() {
			continue
		}

		batch.Queue(query, record.FeatureID, record.EnvironmentID, record.AlgorithmSlug, record.VariantKey,
			record.Evaluations, record.Successes, record.Failures, record.MetricSum)
	}

	if batch.Len() == 0 {
		return nil
	}

	br := executor.SendBatch(ctx, batch)
	defer func() { _ = br.Close() }()

	for range batch.Len() {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}

	return nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
