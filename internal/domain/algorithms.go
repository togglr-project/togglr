package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type AlgorithmID string

type AlgorithmType uint8

const (
	AlgorithmTypeUnknown AlgorithmType = iota
	AlgorithmTypeEpsilonGreedy
	AlgorithmTypeThompsonSampling
	AlgorithmTypeUCB
)

type AlgorithmKind string

const (
	AlgorithmKindBandit AlgorithmKind = "bandit"
)

type Algorithm struct {
	ID              AlgorithmID
	Name            string
	Slug            string
	Kind            AlgorithmKind
	Description     string
	DefaultSettings map[string]decimal.Decimal
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
