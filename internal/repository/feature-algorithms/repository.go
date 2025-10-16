package feature_algorithms

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/repository/auditlog"
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
	project_id,
	feature_id,
	environment_id,
	algorithm_slug,
	settings,
	enabled,
	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING id`

	var id string
	err := executor.QueryRow(
		ctx,
		query,
		featureAlgorithm.ProjectID,
		featureAlgorithm.FeatureID,
		featureAlgorithm.EnvironmentID,
		featureAlgorithm.AlgorithmSlug,
		featureAlgorithm.Settings,
		featureAlgorithm.Enabled,
	).Scan(&id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		featureAlgorithm.ProjectID,
		featureAlgorithm.FeatureID,
		domain.EntityFeatureAlgorithm,
		id,
		domain.AuditActionCreate,
		nil,
		featureAlgorithm,
		featureAlgorithm.EnvironmentID,
	); err != nil {
		return fmt.Errorf("audit feature_algorithm create: %w", err)
	}

	return nil
}

func (r *Repository) Update(
	ctx context.Context,
	featureAlgorithm domain.FeatureAlgorithm,
) error {
	executor := r.getExecutor(ctx)

	// Read old state for audit within the same transaction.
	oldFeatAlg, err := r.GetByID(ctx, featureAlgorithm.ID)
	if err != nil {
		return fmt.Errorf("get feature_algorithm before update: %w", err)
	}

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

	_, err = executor.Exec(
		ctx,
		query,
		featureAlgorithm.Settings,
		featureAlgorithm.AlgorithmSlug,
		featureAlgorithm.Enabled,
		featureAlgorithm.FeatureID,
		featureAlgorithm.EnvironmentID,
	)
	if err != nil {
		return fmt.Errorf("update feature_algorithm: %w", err)
	}

	if err := auditlog.Write(
		ctx,
		executor,
		featureAlgorithm.ProjectID,
		featureAlgorithm.FeatureID,
		domain.EntityFeatureAlgorithm,
		featureAlgorithm.ID.String(),
		domain.AuditActionUpdate,
		oldFeatAlg,
		featureAlgorithm,
		featureAlgorithm.EnvironmentID,
	); err != nil {
		return fmt.Errorf("audit feature_algorithm create: %w", err)
	}

	return nil
}

func (r *Repository) Delete(
	ctx context.Context,
	id domain.FeatureAlgorithmID,
) error {
	executor := r.getExecutor(ctx)

	// Read old state for audit within the same transaction.
	oldFeatAlg, err := r.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get feature_algorithm before update: %w", err)
	}

	const query = `DELETE FROM feature_algorithms WHERE id = $1`

	_, err = executor.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		oldFeatAlg.ProjectID,
		oldFeatAlg.FeatureID,
		domain.EntityFeatureAlgorithm,
		id.String(),
		domain.AuditActionDelete,
		oldFeatAlg,
		nil,
		oldFeatAlg.EnvironmentID,
	); err != nil {
		return fmt.Errorf("audit feature_algorithm create: %w", err)
	}

	return nil
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

func (r *Repository) GetByID(ctx context.Context, id domain.FeatureAlgorithmID) (domain.FeatureAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_algorithms WHERE id = $1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.FeatureAlgorithm{}, err
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureAlgorithmModel])
	if err != nil {
		return domain.FeatureAlgorithm{}, fmt.Errorf("collect feature_algorithms rows: %w", err)
	}

	return model.toDomain(), nil
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

func (r *Repository) ListAllExtended(ctx context.Context) ([]domain.FeatureAlgorithmExtended, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT fa.feature_id,
       fa.algorithm_slug,
       fa.settings,
       fa.environment_id, 
       fa.enabled,
       fa.created_at,
       fa.updated_at,
       e.key AS env_key,
       f.key AS feature_key
FROM feature_algorithms fa
INNER JOIN public.environments e on e.id = fa.environment_id
INNER JOIN public.features f on f.id = fa.feature_id`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureAlgorithmExtModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_algorithms rows: %w", err)
	}

	result := make([]domain.FeatureAlgorithmExtended, 0, len(models))
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
