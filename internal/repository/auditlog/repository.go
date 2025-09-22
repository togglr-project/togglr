package auditlog

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

type auditLogModel struct {
	ID        uint64    `db:"id"`
	ProjectID string    `db:"project_id"`
	FeatureID string    `db:"feature_id"`
	EntityID  *string   `db:"entity_id"`
	RequestID string    `db:"request_id"`
	Entity    string    `db:"entity"`
	Actor     string    `db:"actor"`
	Username  *string   `db:"username"`
	Action    string    `db:"action"`
	OldValue  []byte    `db:"old_value"`
	NewValue  []byte    `db:"new_value"`
	CreatedAt time.Time `db:"created_at"`
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
		ID:        domain.AuditLogID(m.ID),
		ProjectID: domain.ProjectID(m.ProjectID),
		FeatureID: domain.FeatureID(m.FeatureID),
		EntityID:  entityID,
		RequestID: m.RequestID,
		Entity:    domain.EntityType(m.Entity),
		Actor:     m.Actor,
		Username:  username,
		Action:    domain.AuditAction(m.Action),
		OldValue:  m.OldValue,
		NewValue:  m.NewValue,
		CreatedAt: m.CreatedAt,
	}
}

// ListSince returns audit_log rows with created_at strictly greater than the provided timestamp.
// Results are ordered ascending by created_at to allow deterministic processing.
func (r *Repository) ListSince(ctx context.Context, since time.Time) ([]domain.AuditLog, error) {
	exec := r.getExecutor(ctx)

	const query = `
SELECT id, project_id, feature_id, entity_id, request_id, entity, actor, username, action, old_value, new_value, created_at
FROM audit_log
WHERE created_at > $1
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
		"id", "project_id", "feature_id", "entity_id", "request_id", "entity", "actor", "username", "action",
		"old_value", "new_value", "created_at",
	).From("audit_log")

	// Apply filters
	builder = builder.Where(sq.Eq{"project_id": filter.ProjectID})

	if filter.Actor != nil {
		builder = builder.Where(sq.Eq{"actor": *filter.Actor})
	}
	if filter.Entity != nil {
		builder = builder.Where(sq.Eq{"entity": *filter.Entity})
	}
	if filter.Action != nil {
		builder = builder.Where(sq.Eq{"action": *filter.Action})
	}
	if filter.FeatureID != nil {
		builder = builder.Where(sq.Eq{"feature_id": *filter.FeatureID})
	}
	if filter.From != nil {
		builder = builder.Where(sq.GtOrEq{"created_at": *filter.From})
	}
	if filter.To != nil {
		builder = builder.Where(sq.LtOrEq{"created_at": *filter.To})
	}

	// Apply sorting
	orderCol := "created_at"
	switch filter.SortBy {
	case "created_at", "actor", "entity":
		orderCol = filter.SortBy
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
	countBuilder := sq.Select("COUNT(DISTINCT request_id)").From("audit_log")
	countBuilder = countBuilder.Where(sq.Eq{"project_id": filter.ProjectID})

	if filter.Actor != nil {
		countBuilder = countBuilder.Where(sq.Eq{"actor": *filter.Actor})
	}
	if filter.Entity != nil {
		countBuilder = countBuilder.Where(sq.Eq{"entity": *filter.Entity})
	}
	if filter.Action != nil {
		countBuilder = countBuilder.Where(sq.Eq{"action": *filter.Action})
	}
	if filter.FeatureID != nil {
		countBuilder = countBuilder.Where(sq.Eq{"feature_id": *filter.FeatureID})
	}
	if filter.From != nil {
		countBuilder = countBuilder.Where(sq.GtOrEq{"created_at": *filter.From})
	}
	if filter.To != nil {
		countBuilder = countBuilder.Where(sq.LtOrEq{"created_at": *filter.To})
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

//nolint:ireturn // repository executor pattern
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}
