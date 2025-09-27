package domain

import "time"

type FeatureParams struct {
	FeatureID     FeatureID
	EnvironmentID EnvironmentID
	Enabled       bool
	DefaultValue  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FeatureParamsDTO struct {
	EnvironmentID EnvironmentID
	Enabled       bool
	DefaultValue  string
}
