package errorreports

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

var _ contract.ErrorReportRepository = (*Repository)(nil)

// Repository implements contract.ErrorReportRepository using Postgres/TimescaleDB.
type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository { //nolint:ireturn // conforms to repository constructors
	return &Repository{db: pool}
}

func (r *Repository) Insert(ctx context.Context, report domain.ErrorReport) error {
	exec := getExecutor(ctx, r.db)

	const query = `
		insert into monitoring.error_reports (
			project_id, feature_id, environment_id,
			error_type, error_message, context, event_id
		) values ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := exec.Exec(ctx, query, report.ProjectID, report.FeatureID, report.EnvironmentID, report.ErrorType, report.ErrorMessage, report.Context, report.EventID)
	if err != nil {
		return fmt.Errorf("insert error report: %w", err)
	}

	return nil
}

func (r *Repository) CountRecent(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	window time.Duration,
) (int, error) {
	exec := getExecutor(ctx, r.db)

	const query = `
		select count(*)
		from monitoring.error_reports
		where feature_id = $1 and environment_id = $2 and created_at > now() - make_interval(secs => $3::double precision)
	`

	var cnt int
	secs := int(window.Seconds())
	if err := exec.QueryRow(ctx, query, featureID, envID, secs).Scan(&cnt); err != nil {
		return 0, fmt.Errorf("count recent error reports: %w", err)
	}

	return cnt, nil
}

func (r *Repository) GetHealth(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	window time.Duration,
) (domain.FeatureHealth, error) {
	exec := getExecutor(ctx, r.db)

	const query = `
		with recent as (
			select count(*) as cnt
			from monitoring.error_reports
			where feature_id = $1 and environment_id = $2 and created_at > now() - $3::interval
		), last_err as (
			select created_at as last_at
			from monitoring.error_reports
			where feature_id = $1 and environment_id = $2
			order by created_at desc
			limit 1
		)
		select coalesce(recent.cnt, 0) as cnt, coalesce(last_err.last_at, to_timestamp(0)) as last_at
		from recent cross join last_err
	`

	var (
		cnt     int
		lastErr time.Time
	)

	// If there is no last_err row, cross join fails; use left joins instead.
	const query2 = `
		with recent as (
			select count(*) as cnt
			from monitoring.error_reports
			where feature_id = $1 and environment_id = $2 and created_at > now() - make_interval(secs => $3::double precision)
		), last_err as (
			select created_at as last_at
			from monitoring.error_reports
			where feature_id = $1 and environment_id = $2
			order by created_at desc
			limit 1
		)
		select coalesce(recent.cnt, 0) as cnt, coalesce(last_err.last_at, to_timestamp(0)) as last_at
		from recent left join last_err on true
	`

	secs2 := int(window.Seconds())
	if err := exec.QueryRow(ctx, query2, featureID, envID, secs2).Scan(&cnt, &lastErr); err != nil {
		return domain.FeatureHealth{}, fmt.Errorf("get feature health: %w", err)
	}

	// Enabled flag is computed in use case with feature params; default here.
	return domain.FeatureHealth{
		FeatureID:     featureID,
		EnvironmentID: envID,
		Enabled:       true,
		Status:        "healthy",
		ErrorRate:     0,
		LastErrorAt:   lastErr,
	}, nil
}

// helper to get tx from context
//
//nolint:ireturn // internal helper
func getExecutor(ctx context.Context, fallback db.Tx) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return fallback
}
