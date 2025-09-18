package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/etoggle/internal/contract"
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
INSERT INTO features (project_id, key, name, description, kind, default_variant, enabled, rollout_key)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, project_id, key, name, description, kind, default_variant, enabled, rollout_key, created_at, updated_at`

	var model featureModel

	var desc any
	if feature.Description != "" {
		desc = feature.Description
	} else {
		desc = sql.NullString{}
	}

	var rolloutKey sql.NullString
	if feature.RolloutKey != "" {
		rolloutKey = sql.NullString{Valid: true, String: feature.RolloutKey.String()}
	}

	err := executor.QueryRow(ctx, query,
		feature.ProjectID,
		feature.Key,
		feature.Name,
		desc,
		feature.Kind,
		feature.DefaultVariant,
		feature.Enabled,
		rolloutKey,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Key,
		&model.Name,
		&model.Description,
		&model.Kind,
		&model.DefaultVariant,
		&model.Enabled,
		&model.RolloutKey,
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

// ListByProjectIDFiltered returns features list by project with optional filters and pagination.
func (r *Repository) ListByProjectIDFiltered(
	ctx context.Context,
	projectID domain.ProjectID,
	filter contract.FeaturesListFilter,
) ([]domain.Feature, int, error) {
	executor := r.getExecutor(ctx)

	builder := sq.Select("*").From("features").Where(sq.Eq{"project_id": projectID})
	countBuilder := sq.Select("COUNT(*)").From("features").Where(sq.Eq{"project_id": projectID})

	if filter.Kind != nil {
		builder = builder.Where(sq.Eq{"kind": *filter.Kind})
		countBuilder = countBuilder.Where(sq.Eq{"kind": *filter.Kind})
	}
	if filter.Enabled != nil {
		builder = builder.Where(sq.Eq{"enabled": *filter.Enabled})
		countBuilder = countBuilder.Where(sq.Eq{"enabled": *filter.Enabled})
	}
	if filter.TextSelector != nil && *filter.TextSelector != "" {
		pattern := fmt.Sprintf("%%%s%%", *filter.TextSelector)
		or := sq.Or{
			sq.Expr("key ILIKE ?", pattern),
			sq.Expr("name ILIKE ?", pattern),
			sq.Expr("COALESCE(description, '') ILIKE ?", pattern),
			sq.Expr("COALESCE(rollout_key, '') ILIKE ?", pattern),
		}
		builder = builder.Where(or)
		countBuilder = countBuilder.Where(or)
	}

	// Sorting with whitelist
	orderCol := "created_at"
	switch filter.SortBy {
	case "name", "key", "enabled", "kind", "created_at", "updated_at":
		orderCol = filter.SortBy
	}
	orderDir := "DESC"
	if !filter.SortDesc {
		orderDir = "ASC"
	}
	builder = builder.OrderBy(fmt.Sprintf("%s %s", orderCol, orderDir))

	// Pagination
	page := filter.Page
	perPage := filter.PerPage
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	builder = builder.Limit(uint64(perPage)).Offset(uint64(offset))

	// Build and execute list query
	listSQL, listArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build features list sql: %w", err)
	}
	rows, err := executor.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query features filtered: %w", err)
	}
	defer rows.Close()
	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureModel])
	if err != nil {
		return nil, 0, fmt.Errorf("collect features rows: %w", err)
	}
	items := make([]domain.Feature, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	// Count total
	countSQL, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build features count sql: %w", err)
	}
	var total int
	if err := executor.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count features: %w", err)
	}

	return items, total, nil
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
SET project_id = $1, key = $2, name = $3, description = $4, kind = $5, default_variant = $6, enabled = $7, rollout_key = $8, updated_at = now()
WHERE id = $9
RETURNING id, project_id, key, name, description, kind, default_variant, enabled, rollout_key, created_at, updated_at`

	var model featureModel

	var desc any
	if feature.Description != "" {
		desc = feature.Description
	} else {
		desc = sql.NullString{}
	}

	var rolloutKey sql.NullString
	if feature.RolloutKey != "" {
		rolloutKey = sql.NullString{Valid: true, String: feature.RolloutKey.String()}
	}

	err = executor.QueryRow(ctx, query,
		feature.ProjectID,
		feature.Key,
		feature.Name,
		desc,
		feature.Kind,
		feature.DefaultVariant,
		feature.Enabled,
		rolloutKey,
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
		&model.RolloutKey,
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
