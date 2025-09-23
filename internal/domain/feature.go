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

type Feature struct {
	ID             FeatureID
	ProjectID      ProjectID
	Key            string      // machine name, e.g. "new_ui"
	Name           string      // human readable name
	Description    string      // optional description
	Kind           FeatureKind // "simple" | "multivariant"
	DefaultVariant string      // any value for simple, or variant name for multivariant
	RolloutKey     RuleAttribute
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
