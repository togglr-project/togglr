package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureAlgorithm struct {
	EnvironmentID EnvironmentID
	FeatureID     FeatureID
	AlgorithmID   AlgorithmID
	FlagVariantID *FlagVariantID
	Settings      map[string]decimal.Decimal
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
