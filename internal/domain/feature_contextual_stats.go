package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureContextualStats struct {
	ProjectID      ProjectID
	EnvironmentID  EnvironmentID
	FeatureID      FeatureID
	AlgorithmSlug  string
	VariantKey     string
	FeatureKey     string
	EnvironmentKey string
	FeatureDim     int
	MatrixA        []float64
	VectorB        []float64
	Pulls          uint64
	TotalReward    decimal.Decimal
	Successes      uint64
	Failures       uint64
	UpdatedAt      time.Time
}
