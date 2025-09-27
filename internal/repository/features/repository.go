package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/contract"
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

// Create inserts a new feature and returns the created entity.
//

func (r *Repository) Create(
	ctx context.Context,
	envID domain.EnvironmentID,
	feature domain.BasicFeature,
) (domain.BasicFeature, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO features (project_id, key, name, description, kind, rollout_key)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, project_id, key, name, description, kind, rollout_key, created_at, updated_at`

	var model baseFeatureModel

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
		rolloutKey,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Key,
		&model.Name,
		&model.Description,
		&model.Kind,
		&model.RolloutKey,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return domain.BasicFeature{}, fmt.Errorf("insert feature: %w", err)
	}

	newFeature := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newFeature.ProjectID,
		newFeature.ID,
		domain.EntityFeature,
		newFeature.ID.String(),
		domain.AuditActionCreate,
		nil,
		newFeature,
		envID,
	); err != nil {
		return domain.BasicFeature{}, fmt.Errorf("audit feature create: %w", err)
	}

	return newFeature, nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.FeatureID) (domain.BasicFeature, error) {
	executor := r.getExecutor(ctx)

	// Get basic feature info first
	const basicQuery = `SELECT * FROM features WHERE id = $1 LIMIT 1`
	rows, err := executor.Query(ctx, basicQuery, id)
	if err != nil {
		return domain.BasicFeature{}, fmt.Errorf("query feature by id: %w", err)
	}
	defer rows.Close()

	basicModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[baseFeatureModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BasicFeature{}, domain.ErrEntityNotFound
		}
		return domain.BasicFeature{}, fmt.Errorf("collect feature row: %w", err)
	}

	// Return basic feature without environment-specific data
	// This is used for audit logs where we don't need enabled/default_value
	return basicModel.toDomain(), nil
}

func (r *Repository) GetByIDWithEnvironment(ctx context.Context, id domain.FeatureID, environmentKey string) (domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM v_features_full WHERE id = $1 AND environment_key = $2 LIMIT 1`

	rows, err := executor.Query(ctx, query, id, environmentKey)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("query feature by id with environment: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureFullModel])
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

	// Prefer a deterministic environment row when fetching by key without explicit environment.
	// We choose prod if present, otherwise the first by environment_key order.
	const query = `
SELECT *
FROM v_features_full
WHERE key = $1
ORDER BY (environment_key = 'prod') DESC, environment_key
LIMIT 1`

	rows, err := executor.Query(ctx, query, key)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("query feature by key: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureFullModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Feature{}, domain.ErrEntityNotFound
		}

		return domain.Feature{}, fmt.Errorf("collect feature row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByKeyWithEnvironment(ctx context.Context, key, environmentKey string) (domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM v_features_full WHERE key = $1 AND environment_key = $2 LIMIT 1`

	rows, err := executor.Query(ctx, query, key, environmentKey)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("query feature by key with environment: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[featureFullModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Feature{}, domain.ErrEntityNotFound
		}
		return domain.Feature{}, fmt.Errorf("collect feature row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) List(ctx context.Context, environmentKey string) ([]domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM v_features_full WHERE environment_key = $1 ORDER BY created_at DESC`

	rows, err := executor.Query(ctx, query, environmentKey)
	if err != nil {
		return nil, fmt.Errorf("query features: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureFullModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature rows: %w", err)
	}

	features := make([]domain.Feature, 0, len(models))
	for _, m := range models {
		features = append(features, m.toDomain())
	}

	return features, nil
}

func (r *Repository) ListByProjectID(ctx context.Context, projectID domain.ProjectID, environmentKey string) ([]domain.Feature, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM v_features_full WHERE project_id = $1::uuid AND environment_key = $2 ORDER BY created_at DESC`

	rows, err := executor.Query(ctx, query, projectID, environmentKey)
	if err != nil {
		return nil, fmt.Errorf("query features by project_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureFullModel])
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
	environmentKey string,
	filter contract.FeaturesListFilter,
) ([]domain.Feature, int, error) {
	executor := r.getExecutor(ctx)

	builder := sq.Select("*").
		From("v_features_full vf").
		Where(sq.Eq{"vf.project_id": projectID, "vf.environment_key": environmentKey})

	countBuilder := sq.Select("COUNT(*)").
		From("v_features_full vf").
		Where(sq.Eq{"vf.project_id": projectID, "vf.environment_key": environmentKey})

	if filter.Kind != nil {
		builder = builder.Where(sq.Eq{"vf.kind": *filter.Kind})
		countBuilder = countBuilder.Where(sq.Eq{"vf.kind": *filter.Kind})
	}

	if filter.Enabled != nil {
		if *filter.Enabled {
			builder = builder.Where(sq.Eq{"vf.enabled": true})
			countBuilder = countBuilder.Where(sq.Eq{"vf.enabled": true})
		} else {
			builder = builder.Where(sq.Or{
				sq.Eq{"vf.enabled": false},
				sq.Eq{"vf.enabled": nil},
			})
			countBuilder = countBuilder.Where(sq.Or{
				sq.Eq{"vf.enabled": false},
				sq.Eq{"vf.enabled": nil},
			})
		}
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

	if len(filter.TagIDs) > 0 {
		// Use subquery to filter by tag IDs to avoid GROUP BY issues
		subquery := sq.Select("DISTINCT feature_id").From("feature_tags").Where(sq.Eq{"tag_id": filter.TagIDs})
		subquerySQL, subqueryArgs, _ := subquery.PlaceholderFormat(sq.Question).ToSql()

		builder = builder.Where(sq.Expr("vf.id IN ("+subquerySQL+")", subqueryArgs...))
		countBuilder = countBuilder.Where(sq.Expr("vf.id IN ("+subquerySQL+")", subqueryArgs...))
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

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[featureFullModel])
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
func (r *Repository) Update(
	ctx context.Context,
	envID domain.EnvironmentID,
	feature domain.BasicFeature,
) (domain.BasicFeature, error) {
	executor := r.getExecutor(ctx)

	// Read old state for audit purposes within the same transaction.
	oldFeature, err := r.GetByID(ctx, feature.ID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.BasicFeature{}, err
		}

		return domain.BasicFeature{}, fmt.Errorf("get feature before update: %w", err)
	}

	const query = `
UPDATE features
SET project_id = $1::uuid, key = $2, name = $3, description = $4, kind = $5, rollout_key = $6, updated_at = now()
WHERE id = $7
RETURNING id, project_id, key, name, description, kind, rollout_key, created_at, updated_at`

	var model baseFeatureModel

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
		rolloutKey,
		feature.ID,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Key,
		&model.Name,
		&model.Description,
		&model.Kind,
		&model.RolloutKey,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BasicFeature{}, domain.ErrEntityNotFound
		}

		return domain.BasicFeature{}, fmt.Errorf("update feature: %w", err)
	}

	newFeature := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newFeature.ProjectID,
		newFeature.ID,
		domain.EntityFeature,
		newFeature.ID.String(),
		domain.AuditActionUpdate,
		oldFeature,
		newFeature,
		envID,
	); err != nil {
		return domain.BasicFeature{}, fmt.Errorf("audit feature update: %w", err)
	}

	return newFeature, nil
}

func (r *Repository) Delete(ctx context.Context, envID domain.EnvironmentID, id domain.FeatureID) error {
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
		oldFeature.ID.String(),
		domain.AuditActionDelete,
		oldFeature,
		nil,
		envID,
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
