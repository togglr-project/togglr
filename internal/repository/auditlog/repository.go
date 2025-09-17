package auditlog

import (
	"context"
	"fmt"
	"time"

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
	Entity    string    `db:"entity"`
	Actor     string    `db:"actor"`
	Action    string    `db:"action"`
	OldValue  []byte    `db:"old_value"`
	NewValue  []byte    `db:"new_value"`
	CreatedAt time.Time `db:"created_at"`
}

func (m auditLogModel) toDomain() domain.AuditLog {
	return domain.AuditLog{
		ID:        domain.AuditLogID(m.ID),
		ProjectID: domain.ProjectID(m.ProjectID),
		FeatureID: domain.FeatureID(m.FeatureID),
		Entity:    domain.EntityType(m.Entity),
		Actor:     m.Actor,
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
SELECT id, project_id, feature_id, entity, actor, action, old_value, new_value, created_at
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

//nolint:ireturn // repository executor pattern
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}
