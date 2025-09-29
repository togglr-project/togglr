package domain

import (
	"time"
)

type FeatureID string

type FeatureKind string

const (
	FeatureKindSimple       FeatureKind = "simple"
	FeatureKindMultivariant FeatureKind = "multivariant"
)

type BasicFeature struct {
	ID          FeatureID     `db:"id" pk:"true"`
	ProjectID   ProjectID     `db:"project_id"`
	Key         string        `db:"key"`                         // machine name, e.g. "new_ui"
	Kind        FeatureKind   `db:"kind"`                        // "simple" | "multivariant"
	Name        string        `db:"name" editable:"true"`        // human readable name
	Description string        `db:"description" editable:"true"` // optional description
	RolloutKey  RuleAttribute `db:"rollout_key" editable:"true"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}

type Feature struct {
	BasicFeature

	EnvironmentID EnvironmentID
	Enabled       bool   // whether the feature is enabled in the specified environment
	DefaultValue  string // default value for the feature in the specified environment
}

type FeatureExtended struct {
	Feature

	FlagVariants []FlagVariant
	Rules        []Rule
	Schedules    []FeatureSchedule
	Tags         []Tag
}

func (id FeatureID) String() string {
	return string(id)
}

func (kind FeatureKind) String() string {
	return string(kind)
}

// GuardedResult represents the result of a guarded operation.
type GuardedResult struct {
	Pending        bool           // true if operation created a pending change
	ChangeConflict bool           // true if there's a conflict with existing pending changes
	PendingChange  *PendingChange // the created pending change (if any)
	Error          error          // any error that occurred
}

func (f *Feature) ConvertToFeatureParams() FeatureParams {
	return FeatureParams{
		FeatureID:     f.ID,
		EnvironmentID: f.EnvironmentID,
		Enabled:       f.Enabled,
		DefaultValue:  f.DefaultValue,
	}
}
