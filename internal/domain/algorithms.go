package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type AlgorithmType string

const (
	AlgorithmTypeUnknown          AlgorithmType = "unknown"
	AlgorithmTypeEpsilonGreedy                  = "epsilon-greedy"
	AlgorithmTypeThompsonSampling               = "thompson-sampling"
	AlgorithmTypeUCB                            = "ucb"
)

type AlgorithmKind string

const (
	AlgorithmKindBandit AlgorithmKind = "bandit"
)

type Algorithm struct {
	Slug            string
	Name            string
	Kind            AlgorithmKind
	Description     string
	DefaultSettings map[string]decimal.Decimal
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (alg *Algorithm) AlgorithmType() AlgorithmType {
	return AlgorithmSlugToType(alg.Slug)
}

func AlgorithmSlugToType(slug string) AlgorithmType {
	switch slug {
	case "epsilon-greedy":
		return AlgorithmTypeEpsilonGreedy
	case "thompson-sampling":
		return AlgorithmTypeThompsonSampling
	case "ucb":
		return AlgorithmTypeUCB
	default:
		return AlgorithmTypeUnknown
	}
}

func (algType AlgorithmType) Slug() string {
	return string(algType)
}
