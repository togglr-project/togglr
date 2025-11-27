package contract

import (
	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type AlgorithmsProcessor interface {
	HasAlgorithm(featureKey, envKey string) bool
	GetAlgorithmKind(featureKey, envKey string) (domain.AlgorithmKind, bool)
	IsCustomAlgorithm(featureKey, envKey string) bool

	// EvaluateFeature for multi-variant bandits (non-contextual)
	EvaluateFeature(featureKey, envKey string) (string, bool)

	// EvaluateContextual for contextual bandits (uses user context)
	EvaluateContextual(featureKey, envKey string, ctx map[string]any) (string, bool)

	// EvaluateOptimizer for single-variant optimizers
	EvaluateOptimizer(featureKey, envKey string) (decimal.Decimal, bool)

	// EvaluateCustom for custom WASM bandit algorithms
	EvaluateCustom(featureKey, envKey string, ctx map[string]any) (string, bool)

	// EvaluateCustomOptimizer for custom WASM optimizer algorithms
	EvaluateCustomOptimizer(featureKey, envKey string) (decimal.Decimal, bool)

	HandleTrackEvent(
		featureKey string,
		envKey string,
		variantKey string,
		eventType domain.FeedbackEventType,
		metric decimal.Decimal,
	)

	// HandleContextualTrackEvent for contextual bandits with context
	HandleContextualTrackEvent(
		featureKey string,
		envKey string,
		variantKey string,
		eventType domain.FeedbackEventType,
		metric decimal.Decimal,
		ctx map[string]any,
	)

	// HandleCustomTrackEvent for custom WASM algorithms
	HandleCustomTrackEvent(
		featureKey string,
		envKey string,
		variantKey string,
		eventType domain.FeedbackEventType,
		metric decimal.Decimal,
		ctx map[string]any,
	)
}
