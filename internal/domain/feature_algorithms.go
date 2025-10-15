package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureAlgorithm struct {
	EnvironmentID EnvironmentID
	FeatureID     FeatureID
	AlgorithmSlug string
	Settings      map[string]decimal.Decimal
	Enabled       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FeatureAlgorithmExtended struct {
	FeatureAlgorithm
	FeatureKey string
	EnvKey     string
}

type FeatureAlgorithmDTO struct {
	EnvironmentID EnvironmentID
	FeatureID     FeatureID
	AlgorithmSlug string
	Enabled       bool
	Settings      map[string]decimal.Decimal
}
