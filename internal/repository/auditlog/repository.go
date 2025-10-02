package auditlog

import (
	"context"
	"errors"
	"fmt"
	"time"

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

type auditLogModel struct {
	ID            uint64    `db:"id"`
	ProjectID     string    `db:"project_id"`
	FeatureID     string    `db:"feature_id"`
	EntityID      *string   `db:"entity_id"`
	RequestID     string    `db:"request_id"`
	Entity        string    `db:"entity"`
	Actor         string    `db:"actor"`
	Username      *string   `db:"username"`
	Action        string    `db:"action"`
	OldValue      []byte    `db:"old_value"`
	NewValue      []byte    `db:"new_value"`
	EnvironmentID int64     `db:"environment_id"`
	EnvKey        string    `db:"env_key"`
	CreatedAt     time.Time `db:"created_at"`
}

func (m auditLogModel) toDomain() domain.AuditLog {
	entityID := ""
	if m.EntityID != nil {
		entityID = *m.EntityID
	}

	username := ""
	if m.Username != nil {
		username = *m.Username
	}

	return domain.AuditLog{
		ID:            domain.AuditLogID(m.ID),
		ProjectID:     domain.ProjectID(m.ProjectID),
		FeatureID:     domain.FeatureID(m.FeatureID),
		EntityID:      entityID,
		RequestID:     m.RequestID,
		Entity:        domain.EntityType(m.Entity),
		Actor:         m.Actor,
		Username:      username,
		Action:        domain.AuditAction(m.Action),
		OldValue:      m.OldValue,
		NewValue:      m.NewValue,
		EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
		EnvKey:        m.EnvKey,
		CreatedAt:     m.CreatedAt,
	}
}

// ListSince returns audit_log rows with created_at strictly greater than the provided timestamp.
// Results are ordered ascending by created_at to allow deterministic processing.
func (r *Repository) ListSince(ctx context.Context, since time.Time) ([]domain.AuditLog, error) {
	exec := r.getExecutor(ctx)

	const query = `
SELECT audit_log.id, audit_log.project_id, feature_id, entity_id, request_id, entity,
       actor, username, action, old_value, new_value, environment_id, COALESCE(envs.key, '') AS env_key, audit_log.created_at
FROM audit_log
LEFT JOIN environments envs ON audit_log.environment_id = envs.id
WHERE audit_log.created_at > $1
ORDER BY created_at
`

	rows, err := exec.Query(ctx, query, since)
	if err != nil {
		return nil, fmt.Errorf("query audit_log since: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[auditLogModel])
	if err != nil {
		return nil, fmt.Errorf("collect audit_log rows: %w", err)
	}

	result := make([]domain.AuditLog, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

// ListChanges returns paginated audit log changes grouped by request_id with filtering.
func (r *Repository) ListChanges(
	ctx context.Context,
	filter domain.ChangesListFilter,
) (domain.ChangesListResult, error) {
	exec := r.getExecutor(ctx)

	// Build base query for changes
	builder := sq.Select(
		"audit_log.id", "audit_log.project_id", "audit_log.feature_id", "audit_log.entity_id", "audit_log.request_id",
		"audit_log.entity", "audit_log.actor", "audit_log.username", "audit_log.action",
		"audit_log.old_value", "audit_log.new_value", "audit_log.environment_id", "COALESCE(envs.key, '') AS env_key", "audit_log.created_at",
	).From("audit_log").
		LeftJoin("environments envs ON audit_log.environment_id = envs.id")

	// Apply filters
	builder = builder.Where(sq.Eq{"audit_log.project_id": filter.ProjectID})

	if filter.Actor != nil {
		builder = builder.Where(sq.Eq{"audit_log.actor": *filter.Actor})
	}

	if filter.Entity != nil {
		builder = builder.Where(sq.Eq{"audit_log.entity": *filter.Entity})
	}

	if filter.Action != nil {
		builder = builder.Where(sq.Eq{"audit_log.action": *filter.Action})
	}

	if filter.FeatureID != nil {
		builder = builder.Where(sq.Eq{"audit_log.feature_id": *filter.FeatureID})
	}

	if filter.From != nil {
		builder = builder.Where(sq.GtOrEq{"audit_log.created_at": *filter.From})
	}

	if filter.To != nil {
		builder = builder.Where(sq.LtOrEq{"audit_log.created_at": *filter.To})
	}

	// Apply sorting
	orderCol := "audit_log.created_at"

	switch filter.SortBy {
	case "created_at", "actor", "entity":
		orderCol = "audit_log." + filter.SortBy
	}

	orderDir := "DESC"
	if !filter.SortDesc {
		orderDir = "ASC"
	}

	builder = builder.OrderBy(fmt.Sprintf("%s %s", orderCol, orderDir))

	// Apply pagination
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

	// Execute query
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.ChangesListResult{}, fmt.Errorf("build changes query: %w", err)
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return domain.ChangesListResult{}, fmt.Errorf("query changes: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[auditLogModel])
	if err != nil {
		return domain.ChangesListResult{}, fmt.Errorf("collect changes rows: %w", err)
	}

	// Group changes by request_id
	changeGroups := make(map[string]*domain.ChangeGroup)

	for _, model := range models {
		auditLog := model.toDomain()

		// Get or create change group
		group, exists := changeGroups[auditLog.RequestID]
		if !exists {
			group = &domain.ChangeGroup{
				RequestID: auditLog.RequestID,
				Actor:     auditLog.Actor,
				Username:  auditLog.Username,
				CreatedAt: auditLog.CreatedAt,
				Changes:   make([]domain.Change, 0),
			}
			changeGroups[auditLog.RequestID] = group
		}

		// Convert to Change
		change := domain.Change{
			ID:       auditLog.ID,
			Entity:   auditLog.Entity,
			EntityID: auditLog.EntityID, // Using EntityID field
			Action:   auditLog.Action,
		}

		// Handle old_value
		if len(auditLog.OldValue) > 0 {
			change.OldValue = &auditLog.OldValue
		}

		// Handle new_value
		if len(auditLog.NewValue) > 0 {
			change.NewValue = &auditLog.NewValue
		}

		group.Changes = append(group.Changes, change)
	}

	// Convert map to slice
	items := make([]domain.ChangeGroup, 0, len(changeGroups))
	for _, group := range changeGroups {
		items = append(items, *group)
	}

	// Get total count
	countBuilder := sq.Select("COUNT(DISTINCT audit_log.request_id)").From("audit_log").
		LeftJoin("environments envs ON audit_log.environment_id = envs.id")
	countBuilder = countBuilder.Where(sq.Eq{"audit_log.project_id": filter.ProjectID})

	if filter.Actor != nil {
		countBuilder = countBuilder.Where(sq.Eq{"audit_log.actor": *filter.Actor})
	}

	if filter.Entity != nil {
		countBuilder = countBuilder.Where(sq.Eq{"audit_log.entity": *filter.Entity})
	}

	if filter.Action != nil {
		countBuilder = countBuilder.Where(sq.Eq{"audit_log.action": *filter.Action})
	}

	if filter.FeatureID != nil {
		countBuilder = countBuilder.Where(sq.Eq{"audit_log.feature_id": *filter.FeatureID})
	}

	if filter.From != nil {
		countBuilder = countBuilder.Where(sq.GtOrEq{"audit_log.created_at": *filter.From})
	}

	if filter.To != nil {
		countBuilder = countBuilder.Where(sq.LtOrEq{"audit_log.created_at": *filter.To})
	}

	countQuery, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.ChangesListResult{}, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if err := exec.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return domain.ChangesListResult{}, fmt.Errorf("count changes: %w", err)
	}

	return domain.ChangesListResult{
		ProjectID: filter.ProjectID,
		Items:     items,
		Total:     total,
	}, nil
}

// ListByProjectIDFiltered returns paginated audit logs for a project with filters and sorting.
func (r *Repository) ListByProjectIDFiltered(
	ctx context.Context,
	opts contract.AuditLogListFilter,
) (items []domain.AuditLog, total int, err error) {
	exec := r.getExecutor(ctx)

	builder := sq.Select(
		"audit_log.id", "audit_log.project_id", "audit_log.feature_id", "audit_log.entity_id", "audit_log.request_id",
		"audit_log.entity", "audit_log.actor", "audit_log.username", "audit_log.action",
		"audit_log.old_value", "audit_log.new_value", "audit_log.environment_id", "COALESCE(envs.key, '') AS env_key", "audit_log.created_at",
	).From("audit_log").
		LeftJoin("environments envs ON audit_log.environment_id = envs.id")

	builder = builder.Where(sq.Eq{"audit_log.project_id": opts.ProjectID})

	if opts.EnvironmentKey != nil {
		builder = builder.Where(sq.Eq{"envs.key": *opts.EnvironmentKey})
	}
	if opts.Entity != nil {
		builder = builder.Where(sq.Eq{"audit_log.entity": *opts.Entity})
	}
	if opts.EntityID != nil {
		builder = builder.Where(sq.Eq{"audit_log.entity_id": *opts.EntityID})
	}
	if opts.Actor != nil {
		builder = builder.Where(sq.Eq{"audit_log.actor": *opts.Actor})
	}
	if opts.From != nil {
		builder = builder.Where(sq.GtOrEq{"audit_log.created_at": *opts.From})
	}
	if opts.To != nil {
		builder = builder.Where(sq.LtOrEq{"audit_log.created_at": *opts.To})
	}

	orderCol := "audit_log.created_at"
	switch opts.SortBy {
	case "environment_key":
		orderCol = "envs.key"
	case "entity":
		orderCol = "audit_log.entity"
	case "entity_id":
		orderCol = "audit_log.entity_id"
	case "actor":
		orderCol = "audit_log.actor"
	case "action":
		orderCol = "audit_log.action"
	case "username":
		orderCol = "audit_log.username"
	case "created_at":
		orderCol = "audit_log.created_at"
	}
	orderDir := "DESC"
	if !opts.SortDesc {
		orderDir = "ASC"
	}
	builder = builder.OrderBy(fmt.Sprintf("%s %s", orderCol, orderDir))

	page := opts.Page
	perPage := opts.PerPage
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	builder = builder.Limit(uint64(perPage)).Offset(uint64(offset))

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build audit logs query: %w", err)
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query audit logs: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[auditLogModel])
	if err != nil {
		return nil, 0, fmt.Errorf("collect audit_log rows: %w", err)
	}

	items = make([]domain.AuditLog, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	// Count total
	countBuilder := sq.Select("COUNT(*)").From("audit_log").
		LeftJoin("environments envs ON audit_log.environment_id = envs.id")
	countBuilder = countBuilder.Where(sq.Eq{"audit_log.project_id": opts.ProjectID})
	if opts.EnvironmentKey != nil {
		countBuilder = countBuilder.Where(sq.Eq{"envs.key": *opts.EnvironmentKey})
	}
	if opts.Entity != nil {
		countBuilder = countBuilder.Where(sq.Eq{"audit_log.entity": *opts.Entity})
	}
	if opts.EntityID != nil {
		countBuilder = countBuilder.Where(sq.Eq{"audit_log.entity_id": *opts.EntityID})
	}
	if opts.Actor != nil {
		countBuilder = countBuilder.Where(sq.Eq{"audit_log.actor": *opts.Actor})
	}
	if opts.From != nil {
		countBuilder = countBuilder.Where(sq.GtOrEq{"audit_log.created_at": *opts.From})
	}
	if opts.To != nil {
		countBuilder = countBuilder.Where(sq.LtOrEq{"audit_log.created_at": *opts.To})
	}
	countQuery, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count audit logs query: %w", err)
	}
	if err := exec.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	return items, total, nil
}

// GetByID returns a single audit log entry by id.
func (r *Repository) GetByID(ctx context.Context, id domain.AuditLogID) (domain.AuditLog, error) {
	exec := r.getExecutor(ctx)

	builder := sq.Select(
		"audit_log.id", "audit_log.project_id", "audit_log.feature_id", "audit_log.entity_id", "audit_log.request_id",
		"audit_log.entity", "audit_log.actor", "audit_log.username", "audit_log.action",
		"audit_log.old_value", "audit_log.new_value", "audit_log.environment_id", "COALESCE(envs.key, '') AS env_key", "audit_log.created_at",
	).From("audit_log").
		LeftJoin("environments envs ON audit_log.environment_id = envs.id").
		Where(sq.Eq{"audit_log.id": id})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.AuditLog{}, fmt.Errorf("build audit log by id query: %w", err)
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return domain.AuditLog{}, fmt.Errorf("query audit log by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[auditLogModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AuditLog{}, domain.ErrEntityNotFound
		}

		return domain.AuditLog{}, fmt.Errorf("collect audit_log row: %w", err)
	}

	return model.toDomain(), nil
}

//nolint:ireturn // repository executor pattern
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
