package customalgorithmstats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

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

type statsModel struct {
	ProjectID      string          `db:"project_id"`
	FeatureID      string          `db:"feature_id"`
	EnvironmentID  int64           `db:"environment_id"`
	AlgorithmID    string          `db:"algorithm_id"`
	VariantKey     string          `db:"variant_key"`
	FeatureKey     string          `db:"feature_key"`
	EnvironmentKey string          `db:"environment_key"`
	State          json.RawMessage `db:"state"`
	Evaluations    int64           `db:"evaluations"`
	Successes      int64           `db:"successes"`
	Failures       int64           `db:"failures"`
	MetricSum      decimal.Decimal `db:"metric_sum"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

func (m *statsModel) toDomain() domain.CustomAlgorithmStats {
	return domain.CustomAlgorithmStats{
		ProjectID:      domain.ProjectID(m.ProjectID),
		FeatureID:      domain.FeatureID(m.FeatureID),
		EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
		AlgorithmID:    domain.CustomAlgorithmID(m.AlgorithmID),
		VariantKey:     m.VariantKey,
		FeatureKey:     m.FeatureKey,
		EnvironmentKey: m.EnvironmentKey,
		State:          m.State,
		Evaluations:    uint64(m.Evaluations),
		Successes:      uint64(m.Successes),
		Failures:       uint64(m.Failures),
		MetricSum:      m.MetricSum,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (r *Repository) Upsert(ctx context.Context, stats domain.CustomAlgorithmStats) error {
	executor := r.getExecutor(ctx)

	stateJSON := stats.State
	if stateJSON == nil {
		stateJSON = []byte("{}")
	}

	const query = `
INSERT INTO monitoring.custom_algorithm_stats
(project_id, feature_id, environment_id, algorithm_id, variant_key, feature_key, environment_key,
 state, evaluations, successes, failures, metric_sum, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
ON CONFLICT (feature_id, environment_id, algorithm_id, variant_key) DO UPDATE SET
    state = EXCLUDED.state,
    evaluations = EXCLUDED.evaluations,
    successes = EXCLUDED.successes,
    failures = EXCLUDED.failures,
    metric_sum = EXCLUDED.metric_sum,
    updated_at = NOW()`

	_, err := executor.Exec(ctx, query,
		stats.ProjectID, stats.FeatureID, stats.EnvironmentID, stats.AlgorithmID,
		stats.VariantKey, stats.FeatureKey, stats.EnvironmentKey,
		stateJSON, stats.Evaluations, stats.Successes, stats.Failures, stats.MetricSum)
	if err != nil {
		return fmt.Errorf("upsert custom algorithm stats: %w", err)
	}

	return nil
}

func (r *Repository) UpsertBatch(ctx context.Context, records []domain.CustomAlgorithmStats) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO monitoring.custom_algorithm_stats
(project_id, feature_id, environment_id, algorithm_id, variant_key, feature_key, environment_key,
 state, evaluations, successes, failures, metric_sum, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
ON CONFLICT (feature_id, environment_id, algorithm_id, variant_key) DO UPDATE SET
    state = EXCLUDED.state,
    evaluations = EXCLUDED.evaluations,
    successes = EXCLUDED.successes,
    failures = EXCLUDED.failures,
    metric_sum = EXCLUDED.metric_sum,
    updated_at = NOW()`

	batch := &pgx.Batch{}
	for _, stats := range records {
		stateJSON := stats.State
		if stateJSON == nil {
			stateJSON = []byte("{}")
		}

		batch.Queue(query,
			stats.ProjectID, stats.FeatureID, stats.EnvironmentID, stats.AlgorithmID,
			stats.VariantKey, stats.FeatureKey, stats.EnvironmentKey,
			stateJSON, stats.Evaluations, stats.Successes, stats.Failures, stats.MetricSum)
	}

	if batch.Len() == 0 {
		return nil
	}

	br := executor.SendBatch(ctx, batch)
	defer func() { _ = br.Close() }()

	for range batch.Len() {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("batch upsert failed: %w", err)
		}
	}

	return nil
}

func (r *Repository) Get(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	algID domain.CustomAlgorithmID,
	variantKey string,
) (domain.CustomAlgorithmStats, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT * FROM monitoring.custom_algorithm_stats
WHERE feature_id = $1 AND environment_id = $2 AND algorithm_id = $3 AND variant_key = $4`

	rows, err := executor.Query(ctx, query, featureID, envID, algID, variantKey)
	if err != nil {
		return domain.CustomAlgorithmStats{}, fmt.Errorf("query custom algorithm stats: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[statsModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CustomAlgorithmStats{}, domain.ErrEntityNotFound
		}
		return domain.CustomAlgorithmStats{}, fmt.Errorf("collect row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByFeature(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	algID domain.CustomAlgorithmID,
) ([]domain.CustomAlgorithmStats, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT * FROM monitoring.custom_algorithm_stats
WHERE feature_id = $1 AND environment_id = $2 AND algorithm_id = $3`

	rows, err := executor.Query(ctx, query, featureID, envID, algID)
	if err != nil {
		return nil, fmt.Errorf("query custom algorithm stats: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[statsModel])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	result := make([]domain.CustomAlgorithmStats, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	algID domain.CustomAlgorithmID,
) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM monitoring.custom_algorithm_stats
WHERE feature_id = $1 AND environment_id = $2 AND algorithm_id = $3`

	_, err := executor.Exec(ctx, query, featureID, envID, algID)
	if err != nil {
		return fmt.Errorf("delete custom algorithm stats: %w", err)
	}

	return nil
}

func (r *Repository) LoadAll(ctx context.Context) ([]domain.CustomAlgorithmStats, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM monitoring.custom_algorithm_stats`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all custom algorithm stats: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[statsModel])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	result := make([]domain.CustomAlgorithmStats, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) LoadByAlgorithmID(ctx context.Context, algID domain.CustomAlgorithmID) ([]domain.CustomAlgorithmStats, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM monitoring.custom_algorithm_stats WHERE algorithm_id = $1`

	rows, err := executor.Query(ctx, query, algID)
	if err != nil {
		return nil, fmt.Errorf("query custom algorithm stats by algorithm: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[statsModel])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	result := make([]domain.CustomAlgorithmStats, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}
