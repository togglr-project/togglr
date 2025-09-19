package featureschedules

import (
	"database/sql"
	"time"

	"github.com/rom8726/etoggle/internal/domain"
)

type scheduleModel struct {
	ID        string         `db:"id"`
	ProjectID string         `db:"project_id"`
	FeatureID string         `db:"feature_id"`
	StartsAt  sql.NullTime   `db:"starts_at"`
	EndsAt    sql.NullTime   `db:"ends_at"`
	CronExpr  sql.NullString `db:"cron_expr"`
	Timezone  string         `db:"timezone"`
	Action    string         `db:"action"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (m *scheduleModel) toDomain() domain.FeatureSchedule {
	var (
		startsAt *time.Time
		endsAt   *time.Time
		cronStr  *string
	)
	if m.StartsAt.Valid {
		startsAt = &m.StartsAt.Time
	}
	if m.EndsAt.Valid {
		endsAt = &m.EndsAt.Time
	}
	if m.CronExpr.Valid {
		cron := m.CronExpr.String
		cronStr = &cron
	}

	return domain.FeatureSchedule{
		ID:        domain.FeatureScheduleID(m.ID),
		ProjectID: domain.ProjectID(m.ProjectID),
		FeatureID: domain.FeatureID(m.FeatureID),
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		CronExpr:  cronStr,
		Timezone:  m.Timezone,
		Action:    domain.FeatureScheduleAction(m.Action),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
