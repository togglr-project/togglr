package contract

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureProcessor interface {
	Evaluate(
		projectID domain.ProjectID,
		featureKey string,
		reqCtx map[domain.RuleAttribute]any,
	) (value string, enabled bool, found bool)
	IsFeatureActive(feature domain.FeatureExtended) bool
	NextState(feature domain.FeatureExtended) (enabled bool, timestamp time.Time)
	BuildFeatureTimeline(
		feature domain.FeatureExtended,
		from time.Time,
		to time.Time,
	) ([]domain.TimelineEvent, error)
}
