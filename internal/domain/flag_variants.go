package domain

type FlagVariantID string

type FlagVariant struct {
	ID             FlagVariantID `db:"id"              pk:"true"`
	ProjectID      ProjectID     `db:"project_id"`
	FeatureID      FeatureID     `db:"feature_id"`
	EnvironmentID  EnvironmentID `db:"environment_id"`
	Name           string        `db:"name"            editable:"true"` // e.g. "A", "B"
	RolloutPercent uint8         `db:"rollout_percent" editable:"true"` // % of traffic (0..100)
}

type FlagVariantExtended struct {
	FlagVariant
	FeatureKey string
	EnvKey     string
}

func (id FlagVariantID) String() string {
	return string(id)
}
