package flagvariants

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
	return &Repository{db: pool}
}

func (r *Repository) Create(ctx context.Context, v domain.FlagVariant) (domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	var (
		query string
		args  []any
	)

	if v.ID != "" {
		// Use client-provided ID
		query = `
INSERT INTO flag_variants (id, project_id, feature_id, environment_id, name, rollout_percent)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, project_id, feature_id, environment_id, name, rollout_percent`
		args = []any{v.ID, v.ProjectID, v.FeatureID, int64(v.EnvironmentID), v.Name, int(v.RolloutPercent)}
	} else {
		query = `
INSERT INTO flag_variants (project_id, feature_id, environment_id, name, rollout_percent)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, project_id, feature_id, environment_id, name, rollout_percent`
		args = []any{v.ProjectID, v.FeatureID, int64(v.EnvironmentID), v.Name, int(v.RolloutPercent)}
	}

	var model flagVariantModel
	if err := executor.QueryRow(ctx, query, args...).Scan(
		&model.ID,
		&model.ProjectID,
		&model.FeatureID,
		&model.EnvironmentID,
		&model.Name,
		&model.RolloutPercent,
	); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("insert flag_variant: %w", err)
	}

	newVariant := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newVariant.ProjectID,
		newVariant.FeatureID,
		domain.EntityFlagVariant,
		newVariant.ID.String(),
		domain.AuditActionCreate,
		nil,
		newVariant,
		newVariant.EnvironmentID,
	); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("audit flag_variant create: %w", err)
	}

	return newVariant, nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.FlagVariantID) (domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM flag_variants WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.FlagVariant{}, fmt.Errorf("query flag_variant by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[flagVariantModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FlagVariant{}, domain.ErrEntityNotFound
		}

		return domain.FlagVariant{}, fmt.Errorf("collect flag_variant row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM flag_variants ORDER BY name ASC`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query flag_variants: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[flagVariantModel])
	if err != nil {
		return nil, fmt.Errorf("collect flag_variant rows: %w", err)
	}

	items := make([]domain.FlagVariant, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) ListExtended(ctx context.Context) ([]domain.FlagVariantExtended, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT fv.id,
       fv.project_id,
       fv.feature_id,
       fv.environment_id,
       fv.name,
       fv.rollout_percent,
       e.key AS env_key,
       f.key AS feature_key
FROM flag_variants fv
INNER JOIN public.environments e on e.id = fv.environment_id
Inner join public.features f on f.id = fv.feature_id
ORDER BY name`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query flag_variants: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[flagVariantExtModel])
	if err != nil {
		return nil, fmt.Errorf("collect flag_variant rows: %w", err)
	}

	items := make([]domain.FlagVariantExtended, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM flag_variants WHERE feature_id = $1 ORDER BY name ASC`

	rows, err := executor.Query(ctx, query, featureID)
	if err != nil {
		return nil, fmt.Errorf("query flag_variants by feature_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[flagVariantModel])
	if err != nil {
		return nil, fmt.Errorf("collect flag_variant rows: %w", err)
	}

	items := make([]domain.FlagVariant, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) ListByFeatureIDWithEnvID(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
) ([]domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM flag_variants WHERE feature_id = $1 AND environment_id = $2 ORDER BY name ASC`

	rows, err := executor.Query(ctx, query, featureID, envID)
	if err != nil {
		return nil, fmt.Errorf("query flag_variants by feature_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[flagVariantModel])
	if err != nil {
		return nil, fmt.Errorf("collect flag_variant rows: %w", err)
	}

	items := make([]domain.FlagVariant, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) Update(ctx context.Context, v domain.FlagVariant) (domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	// Read old state for audit within the same transaction.
	oldVariant, err := r.GetByID(ctx, v.ID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.FlagVariant{}, err
		}

		return domain.FlagVariant{}, fmt.Errorf("get flag_variant before update: %w", err)
	}

	const query = `
UPDATE flag_variants
SET project_id = $1::uuid, feature_id = $2, environment_id = $3, name = $4, rollout_percent = $5
WHERE id = $6
RETURNING id, project_id, feature_id, environment_id, name, rollout_percent`

	var model flagVariantModel
	if err := executor.QueryRow(
		ctx,
		query,
		v.ProjectID,
		v.FeatureID,
		int64(v.EnvironmentID),
		v.Name,
		int(v.RolloutPercent),
		v.ID,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.FeatureID,
		&model.EnvironmentID,
		&model.Name,
		&model.RolloutPercent,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FlagVariant{}, domain.ErrEntityNotFound
		}

		return domain.FlagVariant{}, fmt.Errorf("update flag_variant: %w", err)
	}

	newVariant := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newVariant.ProjectID,
		newVariant.FeatureID,
		domain.EntityFlagVariant,
		newVariant.ID.String(),
		domain.AuditActionUpdate,
		oldVariant,
		newVariant,
		newVariant.EnvironmentID,
	); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("audit flag_variant update: %w", err)
	}

	return newVariant, nil
}

func (r *Repository) Delete(ctx context.Context, id domain.FlagVariantID) error {
	executor := r.getExecutor(ctx)

	oldVariant, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		oldVariant.ProjectID,
		oldVariant.FeatureID,
		domain.EntityFlagVariant,
		oldVariant.ID.String(),
		domain.AuditActionDelete,
		oldVariant,
		nil,
		oldVariant.EnvironmentID,
	); err != nil {
		return fmt.Errorf("audit flag_variant delete: %w", err)
	}

	const query = `DELETE FROM flag_variants WHERE id = $1`

	ct, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete flag_variant: %w", err)
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
