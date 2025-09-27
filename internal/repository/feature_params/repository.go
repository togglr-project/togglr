package feature_params

import (
	"context"
	"errors"
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
	projectID domain.ProjectID,
	params domain.FeatureParams,
) (domain.FeatureParams, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO feature_params (feature_id, environment_id, enabled, default_value, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (feature_id, environment_id) DO UPDATE SET enabled = EXCLUDED.enabled, default_value = EXCLUDED.default_value, updated_at = NOW()
RETURNING feature_id, environment_id, enabled, default_value, created_at, updated_at`

	var model featureParamsModel

	err := executor.QueryRow(ctx, query,
		params.FeatureID,
		params.EnvironmentID,
		params.Enabled,
		params.DefaultValue,
		params.CreatedAt,
		params.UpdatedAt,
	).Scan(
		&model.FeatureID,
		&model.EnvironmentID,
		&model.Enabled,
		&model.DefaultValue,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return domain.FeatureParams{}, fmt.Errorf("insert feature_params: %w", err)
	}

	newParams := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		projectID,
		newParams.FeatureID,
		domain.EntityFeatureParams,
		newParams.FeatureID.String(),
		domain.AuditActionCreate,
		nil,
		newParams,
		newParams.EnvironmentID,
	); err != nil {
		return domain.FeatureParams{}, fmt.Errorf("audit feature_params create: %w", err)
	}

	return newParams, nil
}

func (r *Repository) GetByFeatureWithEnv(ctx context.Context, featureID domain.FeatureID, envID domain.EnvironmentID) (domain.FeatureParams, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_params WHERE feature_id = $1 AND environment_id = $2 LIMIT 1`

	rows, err := executor.Query(ctx, query, featureID, envID)
	if err != nil {
		return domain.FeatureParams{}, fmt.Errorf("query feature_params by feature_id and environment_id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureParamsModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FeatureParams{}, domain.ErrEntityNotFound
		}

		return domain.FeatureParams{}, fmt.Errorf("collect feature_params row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_params WHERE feature_id = $1 ORDER BY environment_id`

	rows, err := executor.Query(ctx, query, featureID)
	if err != nil {
		return nil, fmt.Errorf("query feature_params by feature_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureParamsModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_params rows: %w", err)
	}

	params := make([]domain.FeatureParams, 0, len(models))
	for _, m := range models {
		params = append(params, m.toDomain())
	}

	return params, nil
}

func (r *Repository) ListByEnvironmentID(ctx context.Context, envID domain.EnvironmentID) ([]domain.FeatureParams, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_params WHERE environment_id = $1 ORDER BY feature_id`

	rows, err := executor.Query(ctx, query, envID)
	if err != nil {
		return nil, fmt.Errorf("query feature_params by environment_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureParamsModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature_params rows: %w", err)
	}

	params := make([]domain.FeatureParams, 0, len(models))
	for _, m := range models {
		params = append(params, m.toDomain())
	}

	return params, nil
}

func (r *Repository) Update(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error) {
	executor := r.getExecutor(ctx)

	// Read old state for audit within the same transaction.
	oldParams, err := r.GetByFeatureWithEnv(ctx, params.FeatureID, params.EnvironmentID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.FeatureParams{}, err
		}

		return domain.FeatureParams{}, fmt.Errorf("get feature_params before update: %w", err)
	}

	const query = `
UPDATE feature_params
SET enabled = $1, default_value = $2, updated_at = $3
WHERE feature_id = $4 AND environment_id = $5
RETURNING feature_id, environment_id, enabled, default_value, created_at, updated_at`

	var model featureParamsModel
	err = executor.QueryRow(ctx, query,
		params.Enabled,
		params.DefaultValue,
		params.UpdatedAt,
		params.FeatureID,
		params.EnvironmentID,
	).Scan(
		&model.FeatureID,
		&model.EnvironmentID,
		&model.Enabled,
		&model.DefaultValue,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FeatureParams{}, domain.ErrEntityNotFound
		}

		return domain.FeatureParams{}, fmt.Errorf("update feature_params: %w", err)
	}

	newParams := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		projectID,
		newParams.FeatureID,
		domain.EntityFeatureParams,
		newParams.FeatureID.String(),
		domain.AuditActionUpdate,
		oldParams,
		newParams,
		newParams.EnvironmentID,
	); err != nil {
		return domain.FeatureParams{}, fmt.Errorf("audit feature_params update: %w", err)
	}

	return newParams, nil
}

func (r *Repository) Delete(ctx context.Context, projectID domain.ProjectID, featureID domain.FeatureID, envID domain.EnvironmentID) error {
	executor := r.getExecutor(ctx)

	oldParams, err := r.GetByFeatureWithEnv(ctx, featureID, envID)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		projectID,
		featureID,
		domain.EntityFeatureParams,
		fmt.Sprintf("%s-%d", featureID, envID),
		domain.AuditActionDelete,
		oldParams,
		nil,
		envID,
	); err != nil {
		return fmt.Errorf("audit feature_params delete: %w", err)
	}

	const query = `DELETE FROM feature_params WHERE feature_id = $1 AND environment_id = $2`

	ct, err := executor.Exec(ctx, query, featureID, envID)
	if err != nil {
		return fmt.Errorf("delete feature_params: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
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
