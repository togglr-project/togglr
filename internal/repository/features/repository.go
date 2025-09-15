package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/internal/repository/auditlog"
	"github.com/rom8726/etoggle/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

// Create inserts a new feature and returns the created entity.
//
//nolint:lll // long query string is acceptable
func (r *Repository) Create(ctx context.Context, feature domain.Feature) (domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO features (project_id, key, name, description, kind, default_variant, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, project_id, key, name, description, kind, default_variant, enabled, created_at, updated_at`

	var model featureModel

	var desc any
	if feature.Description != "" {
		desc = feature.Description
	} else {
		desc = sql.NullString{}
	}

	err := executor.QueryRow(ctx, query,
		feature.ProjectID,
		feature.Key,
		feature.Name,
		desc,
		feature.Kind,
		feature.DefaultVariant,
		feature.Enabled,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Key,
		&model.Name,
		&model.Description,
		&model.Kind,
		&model.DefaultVariant,
		&model.Enabled,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("insert feature: %w", err)
	}

	newFeature := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newFeature.ProjectID,
		newFeature.ID,
		domain.EntityFeature,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionCreate,
		nil,
		newFeature,
	); err != nil {
		return domain.Feature{}, fmt.Errorf("audit feature create: %w", err)
	}

	return newFeature, nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.FeatureID) (domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM features WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("query feature by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Feature{}, domain.ErrEntityNotFound
		}

		return domain.Feature{}, fmt.Errorf("collect feature row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByKey(ctx context.Context, key string) (domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM features WHERE key = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, key)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("query feature by key: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Feature{}, domain.ErrEntityNotFound
		}

		return domain.Feature{}, fmt.Errorf("collect feature row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM features ORDER BY created_at DESC`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query features: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature rows: %w", err)
	}

	features := make([]domain.Feature, 0, len(models))
	for _, m := range models {
		features = append(features, m.toDomain())
	}

	return features, nil
}

func (r *Repository) ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM features WHERE project_id = $1 ORDER BY created_at DESC`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query features by project_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature rows: %w", err)
	}

	features := make([]domain.Feature, 0, len(models))
	for _, m := range models {
		features = append(features, m.toDomain())
	}

	return features, nil
}

// Update updates existing feature by ID and returns updated entity.
//
//nolint:lll // long query string is acceptable
func (r *Repository) Update(ctx context.Context, feature domain.Feature) (domain.Feature, error) {
	executor := r.getExecutor(ctx)

	// Read old state for audit purposes within the same transaction.
	oldFeature, err := r.GetByID(ctx, feature.ID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.Feature{}, err
		}
		return domain.Feature{}, fmt.Errorf("get feature before update: %w", err)
	}

	const query = `
UPDATE features
SET project_id = $1, key = $2, name = $3, description = $4, kind = $5, default_variant = $6, enabled = $7, updated_at = now()
WHERE id = $8
RETURNING id, project_id, key, name, description, kind, default_variant, enabled, created_at, updated_at`

	var model featureModel

	var desc any
	if feature.Description != "" {
		desc = feature.Description
	} else {
		desc = sql.NullString{}
	}

	err = executor.QueryRow(ctx, query,
		feature.ProjectID,
		feature.Key,
		feature.Name,
		desc,
		feature.Kind,
		feature.DefaultVariant,
		feature.Enabled,
		feature.ID,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Key,
		&model.Name,
		&model.Description,
		&model.Kind,
		&model.DefaultVariant,
		&model.Enabled,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Feature{}, domain.ErrEntityNotFound
		}

		return domain.Feature{}, fmt.Errorf("update feature: %w", err)
	}

	newFeature := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newFeature.ProjectID,
		newFeature.ID,
		domain.EntityFeature,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionUpdate,
		oldFeature,
		newFeature,
	); err != nil {
		return domain.Feature{}, fmt.Errorf("audit feature update: %w", err)
	}

	return newFeature, nil
}

func (r *Repository) Delete(ctx context.Context, id domain.FeatureID) error {
	executor := r.getExecutor(ctx)

	// Read old state and write audit log before deletion. Note: due to FK cascade, this audit row
	// will also be deleted together with the feature, but we still log the action transactionally.
	oldFeature, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		oldFeature.ProjectID,
		id,
		domain.EntityFeature,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionDelete,
		oldFeature,
		nil,
	); err != nil {
		return fmt.Errorf("audit feature delete: %w", err)
	}

	const query = `DELETE FROM features WHERE id = $1`

	ct, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete feature: %w", err)
	}

	// If nothing was deleted, return not found
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
