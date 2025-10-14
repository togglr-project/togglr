package feature_algorithms

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

func (r *Repository) Create(
	ctx context.Context,
	featureAlgorithm domain.FeatureAlgorithmDTO,
) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO feature_algorithms (
	feature_id,
	environment_id,
	algorithm_slug,
	settings,
	enabled,
	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`

	_, err := executor.Exec(
		ctx,
		query,
		featureAlgorithm.FeatureID,
		featureAlgorithm.EnvironmentID,
		featureAlgorithm.AlgorithmSlug,
		featureAlgorithm.Settings,
		featureAlgorithm.Enabled,
	)

	return err
}

func (r *Repository) Update(
	ctx context.Context,
	featureAlgorithm domain.FeatureAlgorithmDTO,
) error {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE feature_algorithms
SET
	settings = $1,
	algorithm_slug = $2,
	enabled = $3,
	updated_at = NOW()
WHERE
	feature_id = $4 AND
	environment_id = $5`

	_, err := executor.Exec(
		ctx,
		query,
		featureAlgorithm.Settings,
		featureAlgorithm.AlgorithmSlug,
		featureAlgorithm.Enabled,
		featureAlgorithm.FeatureID,
		featureAlgorithm.EnvironmentID,
	)

	return err
}

func (r *Repository) Delete(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) error {
	executor := r.getExecutor(ctx)

	const query = `DELETE FROM feature_algorithms WHERE feature_id = $1 AND environment_id = $2`

	_, err := executor.Exec(ctx, query, featureID, envID)

	return err
}

func (r *Repository) ListByFeatureID(
	ctx context.Context,
	featureID domain.FeatureID,
) ([]domain.FeatureAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_algorithms WHERE feature_id = $1`

	rows, err := executor.Query(ctx, query, featureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureAlgorithmModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_algorithms rows: %w", err)
	}

	result := make([]domain.FeatureAlgorithm, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) ListByFeatureIDWithEnvID(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) ([]domain.FeatureAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT * FROM feature_algorithms
WHERE feature_id = $1 AND environment_id = $2`

	rows, err := executor.Query(ctx, query, featureID, envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureAlgorithmModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_algorithms rows: %w", err)
	}

	result := make([]domain.FeatureAlgorithm, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) ListEnabled(ctx context.Context) ([]domain.FeatureAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_algorithms WHERE enabled = true`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureAlgorithmModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_algorithms rows: %w", err)
	}

	result := make([]domain.FeatureAlgorithm, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) ListAll(ctx context.Context) ([]domain.FeatureAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_algorithms`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureAlgorithmModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_algorithms rows: %w", err)
	}

	result := make([]domain.FeatureAlgorithm, 0, len(models))
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
