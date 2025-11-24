package contract

import (
	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type AlgorithmsProcessor interface {
	HasAlgorithm(featureKey, envKey string) bool
	GetAlgorithmKind(featureKey, envKey string) (domain.AlgorithmKind, bool)

	// EvaluateFeature for multi-variant bandits
	EvaluateFeature(featureKy, envKey string) (string, bool)

	// EvaluateOptimizer for single-variant optimizers
	EvaluateOptimizer(featureKey, envKey string) (decimal.Decimal, bool)

	HandleTrackEvent(
		featureKey string,
		envKey string,
		variantKey string,
		eventType domain.FeedbackEventType,
		metric decimal.Decimal,
	)
}
