package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeatureOptimizerStats struct {
	ProjectID      ProjectID
	EnvironmentID  EnvironmentID
	FeatureID      FeatureID
	AlgorithmSlug  string
	FeatureKey     string
	EnvironmentKey string
	Iteration      uint64
	CurrentValue   decimal.Decimal
	BestValue      decimal.Decimal
	BestReward     decimal.Decimal
	MetricSum      decimal.Decimal
	LastError      decimal.Decimal
	Integral       decimal.Decimal
	StepSize       decimal.Decimal
	Temperature    decimal.Decimal
	UpdatedAt      time.Time
}
