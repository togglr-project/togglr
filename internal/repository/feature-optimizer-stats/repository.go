package feature_optimizer_stats

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

func (r *Repository) LoadAll(ctx context.Context) ([]domain.FeatureOptimizerStats, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM monitoring.feature_optimizer_stats`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureOptimizerStatsModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_optimizer_stats rows: %w", err)
	}

	result := make([]domain.FeatureOptimizerStats, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) InsertBatch(ctx context.Context, records []domain.FeatureOptimizerStats) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO monitoring.feature_optimizer_stats
(project_id, feature_id, environment_id, algorithm_slug, feature_key, environment_key,
	iteration, current_value, best_value, best_reward, metric_sum, last_error, integral,
	step_size, temperature, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW())
ON CONFLICT (feature_id, environment_id, algorithm_slug) DO UPDATE SET
	iteration     = EXCLUDED.iteration,
	current_value = EXCLUDED.current_value,
	best_value    = EXCLUDED.best_value,
	best_reward   = EXCLUDED.best_reward,
	metric_sum    = EXCLUDED.metric_sum,
	last_error    = EXCLUDED.last_error,
	integral      = EXCLUDED.integral,
	step_size     = EXCLUDED.step_size,
	temperature   = EXCLUDED.temperature,
	updated_at    = NOW();`

	batch := &pgx.Batch{}
	for _, record := range records {
		batch.Queue(query, record.ProjectID, record.FeatureID, record.EnvironmentID, record.AlgorithmSlug,
			record.FeatureKey, record.EnvironmentKey,
			record.Iteration, record.CurrentValue, record.BestValue, record.BestReward,
			record.MetricSum, record.LastError, record.Integral, record.StepSize, record.Temperature)
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

//nolint:ireturn
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
