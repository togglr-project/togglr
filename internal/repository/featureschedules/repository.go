package featureschedules

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

func New(pool *pgxpool.Pool) *Repository { //nolint:ireturn // follows existing pattern
	return &Repository{db: pool}
}

//nolint:lll // long SQL strings are fine
func (r *Repository) Create(ctx context.Context, s domain.FeatureSchedule) (domain.FeatureSchedule, error) {
	exec := r.getExecutor(ctx)

	var (
		query string
		args  []any
	)

	starts := sql.NullTime{}
	if s.StartsAt != nil {
		starts.Valid = true
		starts.Time = *s.StartsAt
	}
	ends := sql.NullTime{}
	if s.EndsAt != nil {
		ends.Valid = true
		ends.Time = *s.EndsAt
	}
	cron := sql.NullString{}
	if s.CronExpr != nil {
		cron.Valid = true
		cron.String = *s.CronExpr
	}

	if s.ID != "" {
		query = `
INSERT INTO feature_schedules (id, project_id, feature_id, starts_at, ends_at, cron_expr, timezone, action)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, project_id, feature_id, starts_at, ends_at, cron_expr, timezone, action, created_at`
		args = []any{s.ID, s.ProjectID, s.FeatureID, starts, ends, cron, s.Timezone, s.Action}
	} else {
		query = `
INSERT INTO feature_schedules (project_id, feature_id, starts_at, ends_at, cron_expr, timezone, action)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, project_id, feature_id, starts_at, ends_at, cron_expr, timezone, action, created_at`
		args = []any{s.ProjectID, s.FeatureID, starts, ends, cron, s.Timezone, s.Action}
	}

	var m scheduleModel
	if err := exec.QueryRow(ctx, query, args...).Scan(
		&m.ID,
		&m.ProjectID,
		&m.FeatureID,
		&m.StartsAt,
		&m.EndsAt,
		&m.CronExpr,
		&m.Timezone,
		&m.Action,
		&m.CreatedAt,
	); err != nil {
		return domain.FeatureSchedule{}, fmt.Errorf("insert feature_schedule: %w", err)
	}

	created := m.toDomain()
	if err := auditlog.Write(
		ctx,
		exec,
		created.ProjectID,
		created.FeatureID,
		domain.EntityFeatureSchedule,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionCreate,
		nil,
		created,
	); err != nil {
		return domain.FeatureSchedule{}, fmt.Errorf("audit schedule create: %w", err)
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.FeatureScheduleID) (domain.FeatureSchedule, error) {
	exec := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_schedules WHERE id = $1 LIMIT 1`

	rows, err := exec.Query(ctx, query, id)
	if err != nil {
		return domain.FeatureSchedule{}, fmt.Errorf("query schedule by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[scheduleModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FeatureSchedule{}, domain.ErrEntityNotFound
		}
		return domain.FeatureSchedule{}, fmt.Errorf("collect schedule row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.FeatureSchedule, error) {
	exec := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_schedules ORDER BY created_at`

	rows, err := exec.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query schedules: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[scheduleModel])
	if err != nil {
		return nil, fmt.Errorf("collect schedule rows: %w", err)
	}

	items := make([]domain.FeatureSchedule, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureSchedule, error) {
	exec := r.getExecutor(ctx)

	const query = `SELECT * FROM feature_schedules WHERE feature_id = $1 ORDER BY created_at`

	rows, err := exec.Query(ctx, query, featureID)
	if err != nil {
		return nil, fmt.Errorf("query schedules by feature_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[scheduleModel])
	if err != nil {
		return nil, fmt.Errorf("collect schedule rows: %w", err)
	}

	items := make([]domain.FeatureSchedule, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

//nolint:lll // long SQL strings are fine
func (r *Repository) Update(ctx context.Context, s domain.FeatureSchedule) (domain.FeatureSchedule, error) {
	exec := r.getExecutor(ctx)

	old, err := r.GetByID(ctx, s.ID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.FeatureSchedule{}, err
		}
		return domain.FeatureSchedule{}, fmt.Errorf("get schedule before update: %w", err)
	}

	starts := sql.NullTime{}
	if s.StartsAt != nil {
		starts.Valid = true
		starts.Time = *s.StartsAt
	}
	ends := sql.NullTime{}
	if s.EndsAt != nil {
		ends.Valid = true
		ends.Time = *s.EndsAt
	}
	cron := sql.NullString{}
	if s.CronExpr != nil {
		cron.Valid = true
		cron.String = *s.CronExpr
	}

	const query = `
UPDATE feature_schedules
SET feature_id = $1, starts_at = $2, ends_at = $3, cron_expr = $4, timezone = $5, action = $6
WHERE id = $7
RETURNING id, project_id, feature_id, starts_at, ends_at, cron_expr, timezone, action, created_at`

	var m scheduleModel
	if err := exec.QueryRow(ctx, query,
		s.FeatureID,
		starts,
		ends,
		cron,
		s.Timezone,
		s.Action,
		s.ID,
	).Scan(
		&m.ID,
		&m.ProjectID,
		&m.FeatureID,
		&m.StartsAt,
		&m.EndsAt,
		&m.CronExpr,
		&m.Timezone,
		&m.Action,
		&m.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FeatureSchedule{}, domain.ErrEntityNotFound
		}
		return domain.FeatureSchedule{}, fmt.Errorf("update schedule: %w", err)
	}

	updated := m.toDomain()
	if err := auditlog.Write(
		ctx,
		exec,
		updated.ProjectID,
		updated.FeatureID,
		domain.EntityFeatureSchedule,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionUpdate,
		old,
		updated,
	); err != nil {
		return domain.FeatureSchedule{}, fmt.Errorf("audit schedule update: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id domain.FeatureScheduleID) error {
	exec := r.getExecutor(ctx)

	old, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		exec,
		old.ProjectID,
		old.FeatureID,
		domain.EntityFeatureSchedule,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionDelete,
		old,
		nil,
	); err != nil {
		return fmt.Errorf("audit schedule delete: %w", err)
	}

	const query = `DELETE FROM feature_schedules WHERE id = $1`

	ct, err := exec.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete schedule: %w", err)
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
