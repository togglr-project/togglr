package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureAlgorithmID string

type FeatureAlgorithm struct {
	ID            FeatureAlgorithmID         `db:"id"             pk:"true"`
	ProjectID     ProjectID                  `db:"project_id"`
	EnvironmentID EnvironmentID              `db:"environment_id"`
	FeatureID     FeatureID                  `db:"feature_id"`
	AlgorithmSlug string                     `db:"algorithm_slug" editable:"true"`
	Settings      map[string]decimal.Decimal `db:"settings"       editable:"true"`
	Enabled       bool                       `db:"enabled"        editable:"true"`
	CreatedAt     time.Time                  `db:"created_at"`
	UpdatedAt     time.Time                  `db:"updated_at"`
}

type FeatureAlgorithmExtended struct {
	FeatureAlgorithm
	FeatureKey string
	EnvKey     string
}

type FeatureAlgorithmDTO struct {
	ProjectID     ProjectID
	EnvironmentID EnvironmentID
	FeatureID     FeatureID
	AlgorithmSlug string
	Enabled       bool
	Settings      map[string]decimal.Decimal
}

func (id FeatureAlgorithmID) String() string {
	return string(id)
}
