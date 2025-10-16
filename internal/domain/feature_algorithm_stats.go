package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureAlgorithmStats struct {
	ProjectID      ProjectID
	EnvironmentID  EnvironmentID
	FeatureID      FeatureID
	AlgorithmSlug  string
	VariantKey     string
	FeatureKey     string
	EnvironmentKey string
	Evaluations    uint64
	Successes      uint64
	Failures       uint64
	MetricSum      decimal.Decimal
	UpdatedAt      time.Time
}
