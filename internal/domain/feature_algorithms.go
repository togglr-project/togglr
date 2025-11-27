package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureAlgorithmID string

type FeatureAlgorithm struct {
	ID                FeatureAlgorithmID         `db:"id"                    pk:"true"`
	ProjectID         ProjectID                  `db:"project_id"`
	EnvironmentID     EnvironmentID              `db:"environment_id"`
	FeatureID         FeatureID                  `db:"feature_id"`
	AlgorithmSlug     *string                    `db:"algorithm_slug"        editable:"true"`
	CustomAlgorithmID *CustomAlgorithmID         `db:"custom_algorithm_id"   editable:"true"`
	Settings          map[string]decimal.Decimal `db:"settings"              editable:"true"`
	Enabled           bool                       `db:"enabled"               editable:"true"`
	CreatedAt         time.Time                  `db:"created_at"`
	UpdatedAt         time.Time                  `db:"updated_at"`
}

// IsCustom returns true if this feature uses a custom WASM algorithm.
func (fa *FeatureAlgorithm) IsCustom() bool {
	return fa.CustomAlgorithmID != nil
}

// GetAlgorithmSlug returns the algorithm slug (empty if custom algorithm).
func (fa *FeatureAlgorithm) GetAlgorithmSlug() string {
	if fa.AlgorithmSlug != nil {
		return *fa.AlgorithmSlug
	}
	return ""
}

type FeatureAlgorithmExtended struct {
	FeatureAlgorithm
	FeatureKey string
	EnvKey     string
}

type FeatureAlgorithmDTO struct {
	ProjectID         ProjectID
	EnvironmentID     EnvironmentID
	FeatureID         FeatureID
	AlgorithmSlug     *string
	CustomAlgorithmID *CustomAlgorithmID
	Enabled           bool
	Settings          map[string]decimal.Decimal
}

func (id FeatureAlgorithmID) String() string {
	return string(id)
}
