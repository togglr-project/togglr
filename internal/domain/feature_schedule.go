package domain

import "time"

type FeatureScheduleID string

type FeatureScheduleAction string

const (
	FeatureScheduleActionEnable  FeatureScheduleAction = "enable"
	FeatureScheduleActionDisable FeatureScheduleAction = "disable"
)

type FeatureSchedule struct {
	ID            FeatureScheduleID     `db:"id" pk:"true"`
	ProjectID     ProjectID             `db:"project_id"`
	FeatureID     FeatureID             `db:"feature_id"`
	EnvironmentID EnvironmentID         `db:"environment_id"`
	StartsAt      *time.Time            `db:"starts_at" editable:"true"`
	EndsAt        *time.Time            `db:"ends_at" editable:"true"`
	CronExpr      *string               `db:"cron_expr" editable:"true"`
	CronDuration  *time.Duration        `db:"cron_duration" editable:"true"`
	Timezone      string                `db:"timezone" editable:"true"`
	Action        FeatureScheduleAction `db:"action" editable:"true"`
	CreatedAt     time.Time             `db:"created_at"`
	UpdatedAt     time.Time             `db:"updated_at"`
}

func (id FeatureScheduleID) String() string { return string(id) }

func (a FeatureScheduleAction) String() string { return string(a) }
