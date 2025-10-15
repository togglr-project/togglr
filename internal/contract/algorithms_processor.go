package contract

import (
	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type AlgorithmsProcessor interface {
	HasAlgorithm(featureKey, envKey string) bool
	EvaluateFeature(featureKy, envKey string) (string, bool)
	HandleTrackEvent(
		featureKey string,
		envKey string,
		variantKey string,
		eventType domain.FeedbackEventType,
		metric decimal.Decimal,
	)
}
