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
	ID          FeatureID
	ProjectID   ProjectID
	Key         string      // machine name, e.g. "new_ui"
	Name        string      // human readable name
	Description string      // optional description
	Kind        FeatureKind // "simple" | "multivariant"
	RolloutKey  RuleAttribute
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Feature struct {
	BasicFeature

	Enabled      bool   // whether the feature is enabled in the specified environment
	DefaultValue string // default value for the feature in the specified environment
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

// GetParamsForEnvironment returns the feature parameters for a specific environment
// This method will be implemented in the repository layer
func (f *Feature) GetParamsForEnvironment(envID EnvironmentID, params []FeatureParams) *FeatureParams {
	for _, p := range params {
		if p.EnvironmentID == envID {
			return &p
		}
	}
	return nil
}
