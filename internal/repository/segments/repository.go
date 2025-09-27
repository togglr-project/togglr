package segments

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

// Create inserts a new segment and returns it.
//
//nolint:lll // long query strings are acceptable
func (r *Repository) Create(ctx context.Context, segment domain.Segment) (domain.Segment, error) {
	executor := r.getExecutor(ctx)

	condsData, err := json.Marshal(segment.Conditions)
	if err != nil {
		return domain.Segment{}, fmt.Errorf("marshal conditions: %w", err)
	}

	var (
		query string
		args  []any
	)

	if segment.ID != "" {
		query = `
INSERT INTO segments (id, project_id, name, description, conditions)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, project_id, name, description, conditions, created_at, updated_at`
		args = []any{segment.ID, segment.ProjectID, segment.Name, sql.NullString{String: segment.Description, Valid: segment.Description != ""}, condsData}
	} else {
		query = `
INSERT INTO segments (project_id, name, description, conditions)
VALUES ($1, $2, $3, $4)
RETURNING id, project_id, name, description, conditions, created_at, updated_at`
		args = []any{segment.ProjectID, segment.Name, sql.NullString{String: segment.Description, Valid: segment.Description != ""}, condsData}
	}

	var model segmentModel
	if err := executor.QueryRow(ctx, query, args...).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Name,
		&model.Description,
		&model.Conditions,
		&model.CreatedAt,
		&model.UpdatedAt,
	); err != nil {
		return domain.Segment{}, fmt.Errorf("insert segment: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.SegmentID) (domain.Segment, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM segments WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Segment{}, fmt.Errorf("query segment by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[segmentModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Segment{}, domain.ErrEntityNotFound
		}

		return domain.Segment{}, fmt.Errorf("collect segment row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Segment, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM segments WHERE project_id = $1::uuid ORDER BY name`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query segments by project_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[segmentModel])
	if err != nil {
		return nil, fmt.Errorf("collect segment rows: %w", err)
	}

	items := make([]domain.Segment, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

// ListByProjectIDFiltered returns segments list by project with optional filters and pagination.
func (r *Repository) ListByProjectIDFiltered(
	ctx context.Context,
	projectID domain.ProjectID,
	filter contract.SegmentsListFilter,
) ([]domain.Segment, int, error) {
	executor := r.getExecutor(ctx)

	builder := sq.Select("*").From("segments").Where(sq.Eq{"project_id": projectID})
	countBuilder := sq.Select("COUNT(*)").From("segments").Where(sq.Eq{"project_id": projectID})

	if filter.TextSelector != nil && *filter.TextSelector != "" {
		pattern := fmt.Sprintf("%%%s%%", *filter.TextSelector)
		or := sq.Or{
			sq.Expr("name ILIKE ?", pattern),
			sq.Expr("COALESCE(description, '') ILIKE ?", pattern),
		}
		builder = builder.Where(or)
		countBuilder = countBuilder.Where(or)
	}

	// Sorting with whitelist
	orderCol := "created_at"

	switch filter.SortBy {
	case "name", "created_at", "updated_at":
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
		return nil, 0, fmt.Errorf("build segments list sql: %w", err)
	}

	rows, err := executor.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query segments filtered: %w", err)
	}

	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[segmentModel])
	if err != nil {
		return nil, 0, fmt.Errorf("collect segments rows: %w", err)
	}

	items := make([]domain.Segment, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	// Count total
	countSQL, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build segments count sql: %w", err)
	}

	var total int
	if err := executor.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count segments: %w", err)
	}

	return items, total, nil
}

// Update updates segment fields and returns the updated row.
//

func (r *Repository) Update(ctx context.Context, segment domain.Segment) (domain.Segment, error) {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE segments
SET name = $1, description = $2, conditions = $3, updated_at = now()
WHERE id = $4
RETURNING id, project_id, name, description, conditions, created_at, updated_at`

	condsData, err := json.Marshal(segment.Conditions)
	if err != nil {
		return domain.Segment{}, fmt.Errorf("marshal conditions: %w", err)
	}

	var model segmentModel
	if err := executor.QueryRow(ctx, query,
		segment.Name,
		sql.NullString{String: segment.Description, Valid: segment.Description != ""},
		condsData,
		segment.ID,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Name,
		&model.Description,
		&model.Conditions,
		&model.CreatedAt,
		&model.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Segment{}, domain.ErrEntityNotFound
		}

		return domain.Segment{}, fmt.Errorf("update segment: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id domain.SegmentID) error {
	executor := r.getExecutor(ctx)

	const query = `DELETE FROM segments WHERE id = $1`

	ct, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete segment: %w", err)
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
