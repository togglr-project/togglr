package domain

import "time"

type FeatureParams struct {
	FeatureID     FeatureID     `db:"feature_id"     pk:"true"`
	EnvironmentID EnvironmentID `db:"environment_id"`
	Enabled       bool          `db:"enabled"        editable:"true"`
	DefaultValue  string        `db:"default_value"  editable:"true"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
}

type FeatureParamsDTO struct {
	EnvironmentID EnvironmentID
	Enabled       bool
	DefaultValue  string
}
