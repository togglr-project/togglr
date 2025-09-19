package domain

import "time"

type FeatureScheduleID string

type FeatureScheduleAction string

const (
	FeatureScheduleActionEnable  FeatureScheduleAction = "enable"
	FeatureScheduleActionDisable FeatureScheduleAction = "disable"
)

type FeatureSchedule struct {
	ID           FeatureScheduleID
	ProjectID    ProjectID
	FeatureID    FeatureID
	StartsAt     *time.Time
	EndsAt       *time.Time
	CronExpr     *string
	CronDuration *time.Duration
	Timezone     string
	Action       FeatureScheduleAction
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (id FeatureScheduleID) String() string { return string(id) }

func (a FeatureScheduleAction) String() string { return string(a) }
