package contract

import (
	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type AlgorithmsProcessor interface {
	EvaluateFeature(
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) (string, bool)
	HandleTrackEvent(
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
		variantKey string,
		eventType domain.FeedbackEventType,
		metric decimal.Decimal,
	)
}
