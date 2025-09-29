package guard_engine

import (
	"github.com/togglr-project/togglr/internal/domain"
)

// GuardEngineInput describes a request to create a pending change for a guarded operation.
// It is intentionally small and HTTP-agnostic to keep it usable from different layers.
type GuardEngineInput struct {
	ProjectID         domain.ProjectID
	EnvironmentID     domain.EnvironmentID
	FeatureID         domain.FeatureID
	Reason            string
	Origin            string
	EntityChangesList []EntityChanges
	Action            domain.EntityAction
}
