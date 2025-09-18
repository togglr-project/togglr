package contract

import (
	"github.com/rom8726/etoggle/internal/domain"
)

type FeatureProcessor interface {
	Evaluate(
		projectID domain.ProjectID,
		featureKey string,
		reqCtx map[domain.RuleAttribute]any,
	) (value string, enabled bool, found bool)
}
