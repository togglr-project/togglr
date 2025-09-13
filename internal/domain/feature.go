package domain

import (
	"time"
)

type FeatureID string

type FeatureKind string

const (
	FeatureKindBoolean      FeatureKind = "boolean"
	FeatureKindMultivariant FeatureKind = "multivariant"
)

type Feature struct {
	ID             FeatureID
	ProjectID      ProjectID
	Key            string      // machine name, e.g. "new_ui"
	Name           string      // human readable name
	Description    string      // optional description
	Kind           FeatureKind // "boolean" | "multivariant"
	DefaultVariant string      // "on"/"off" for boolean, or variant name
	FlagVariants   []FlagVariant
	Rules          []Rule
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (id FeatureID) String() string {
	return string(id)
}

func (kind FeatureKind) String() string {
	return string(kind)
}
