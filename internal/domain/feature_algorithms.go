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
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FeatureAlgorithmDTO struct {
	EnvironmentID EnvironmentID
	FeatureID     FeatureID
	AlgorithmSlug string
	Settings      map[string]decimal.Decimal
}
