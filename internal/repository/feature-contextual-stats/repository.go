package feature_contextual_stats

import (
	"context"
	"encoding/json"
	"fmt"

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

type rawModel struct {
	ProjectID      string          `db:"project_id"`
	EnvironmentID  int64           `db:"environment_id"`
	FeatureID      string          `db:"feature_id"`
	AlgorithmSlug  string          `db:"algorithm_slug"`
	VariantKey     string          `db:"variant_key"`
	FeatureKey     string          `db:"feature_key"`
	EnvironmentKey string          `db:"environment_key"`
	FeatureDim     int             `db:"feature_dim"`
	MatrixA        []byte          `db:"matrix_a"`
	VectorB        []byte          `db:"vector_b"`
	Pulls          uint64          `db:"pulls"`
	TotalReward    decimal.Decimal `db:"total_reward"`
	Successes      uint64          `db:"successes"`
	Failures       uint64          `db:"failures"`
}

func (r *Repository) LoadAll(ctx context.Context) ([]domain.FeatureContextualStats, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT project_id, environment_id, feature_id, algorithm_slug, variant_key,
		feature_key, environment_key, feature_dim, matrix_a, vector_b, pulls, total_reward, successes, failures
		FROM monitoring.feature_contextual_stats`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_contextual_stats rows: %w", err)
	}

	result := make([]domain.FeatureContextualStats, 0, len(models))
	for _, m := range models {
		var matrixA []float64
		var vectorB []float64

		if len(m.MatrixA) > 0 {
			if err := json.Unmarshal(m.MatrixA, &matrixA); err != nil {
				return nil, fmt.Errorf("unmarshal matrix_a: %w", err)
			}
		}

		if len(m.VectorB) > 0 {
			if err := json.Unmarshal(m.VectorB, &vectorB); err != nil {
				return nil, fmt.Errorf("unmarshal vector_b: %w", err)
			}
		}

		result = append(result, domain.FeatureContextualStats{
			ProjectID:      domain.ProjectID(m.ProjectID),
			EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
			FeatureID:      domain.FeatureID(m.FeatureID),
			AlgorithmSlug:  m.AlgorithmSlug,
			VariantKey:     m.VariantKey,
			FeatureKey:     m.FeatureKey,
			EnvironmentKey: m.EnvironmentKey,
			FeatureDim:     m.FeatureDim,
			MatrixA:        matrixA,
			VectorB:        vectorB,
			Pulls:          m.Pulls,
			TotalReward:    m.TotalReward,
			Successes:      m.Successes,
			Failures:       m.Failures,
		})
	}

	return result, nil
}

func (r *Repository) InsertBatch(ctx context.Context, records []domain.FeatureContextualStats) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO monitoring.feature_contextual_stats
(project_id, feature_id, environment_id, algorithm_slug, variant_key, feature_key, environment_key,
	feature_dim, matrix_a, vector_b, pulls, total_reward, successes, failures, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,NOW())
ON CONFLICT (feature_id, environment_id, algorithm_slug, variant_key) DO UPDATE SET
	feature_dim  = EXCLUDED.feature_dim,
	matrix_a     = EXCLUDED.matrix_a,
	vector_b     = EXCLUDED.vector_b,
	pulls        = EXCLUDED.pulls,
	total_reward = EXCLUDED.total_reward,
	successes    = EXCLUDED.successes,
	failures     = EXCLUDED.failures,
	updated_at   = NOW();`

	batch := &pgx.Batch{}
	for _, record := range records {
		matrixAJSON, err := json.Marshal(record.MatrixA)
		if err != nil {
			return fmt.Errorf("marshal matrix_a: %w", err)
		}

		vectorBJSON, err := json.Marshal(record.VectorB)
		if err != nil {
			return fmt.Errorf("marshal vector_b: %w", err)
		}

		batch.Queue(query, record.ProjectID, record.FeatureID, record.EnvironmentID, record.AlgorithmSlug,
			record.VariantKey, record.FeatureKey, record.EnvironmentKey,
			record.FeatureDim, matrixAJSON, vectorBJSON, record.Pulls, record.TotalReward,
			record.Successes, record.Failures)
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
