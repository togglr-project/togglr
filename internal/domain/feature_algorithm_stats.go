package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureAlgorithmStats struct {
	FeatureID     FeatureID
	EnvironmentID EnvironmentID
	AlgorithmSlug string
	VariantKey    string
	Evaluations   uint64
	Successes     uint64
	Failures      uint64
	MetricSum     decimal.Decimal
	UpdatedAt     time.Time
}
