package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type AlgorithmType string

const (
	AlgorithmTypeUnknown          AlgorithmType = "unknown"
	AlgorithmTypeEpsilonGreedy    AlgorithmType = "epsilon-greedy"
	AlgorithmTypeThompsonSampling AlgorithmType = "thompson-sampling"
	AlgorithmTypeUCB              AlgorithmType = "ucb"
	AlgorithmTypeHillClimb        AlgorithmType = "hill_climb"
	AlgorithmTypePIDController    AlgorithmType = "pid_controller"
	AlgorithmTypeBayesOpt         AlgorithmType = "bayes_opt"
	AlgorithmTypeCEM              AlgorithmType = "cem"
	AlgorithmTypeSimAnnealing     AlgorithmType = "simulated_annealing"
	// Contextual bandits
	AlgorithmTypeLinUCB             AlgorithmType = "lin_ucb"
	AlgorithmTypeContextualThompson AlgorithmType = "contextual_thompson"
	AlgorithmTypeContextualEpsilon  AlgorithmType = "contextual_epsilon"
)

type AlgorithmKind string

const (
	AlgorithmKindUnknown          AlgorithmKind = "unknown"
	AlgorithmKindBandit           AlgorithmKind = "bandit"
	AlgorithmKindOptimizer        AlgorithmKind = "optimizer"
	AlgorithmKindContextualBandit AlgorithmKind = "contextual_bandit"
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
	case "hill_climb":
		return AlgorithmTypeHillClimb
	case "pid_controller":
		return AlgorithmTypePIDController
	case "bayes_opt":
		return AlgorithmTypeBayesOpt
	case "cem":
		return AlgorithmTypeCEM
	case "simulated_annealing":
		return AlgorithmTypeSimAnnealing
	case "lin_ucb":
		return AlgorithmTypeLinUCB
	case "contextual_thompson":
		return AlgorithmTypeContextualThompson
	case "contextual_epsilon":
		return AlgorithmTypeContextualEpsilon
	default:
		return AlgorithmTypeUnknown
	}
}

func (algType AlgorithmType) Slug() string {
	return string(algType)
}

func (algType AlgorithmType) Kind() AlgorithmKind {
	switch algType {
	case AlgorithmTypeEpsilonGreedy, AlgorithmTypeThompsonSampling, AlgorithmTypeUCB:
		return AlgorithmKindBandit
	case AlgorithmTypeHillClimb, AlgorithmTypePIDController, AlgorithmTypeBayesOpt,
		AlgorithmTypeCEM, AlgorithmTypeSimAnnealing:
		return AlgorithmKindOptimizer
	case AlgorithmTypeLinUCB, AlgorithmTypeContextualThompson, AlgorithmTypeContextualEpsilon:
		return AlgorithmKindContextualBandit
	default:
		return AlgorithmKindUnknown
	}
}
