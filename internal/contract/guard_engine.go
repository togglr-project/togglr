package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

// GuardEngineInput describes a request to create a pending change for a guarded operation.
// It is intentionally small and HTTP-agnostic to keep it usable from different layers.
type GuardEngineInput struct {
	ProjectID       domain.ProjectID
	EnvironmentID   domain.EnvironmentID
	FeatureID       domain.FeatureID
	Reason          string
	Origin          string
	PrimaryEntity   string
	PrimaryEntityID string
	Action          domain.EntityAction
	ExtraChanges    map[string]domain.ChangeValue
}

// GuardEngine encapsulates guard checks and pending change creation.
// It centralizes the guarded workflow so that outer layers (API) don't need to
// know how to assemble payloads or handle conflicts.
type GuardEngine interface {
	// CheckAndMaybeCreatePending checks guarded state and potential conflicts
	// for the involved entities and, when guarded and no conflict exists,
	// creates a pending change.
	//
	// Returns:
	//   - pendingChange (non-nil) when a pending change was created and 202 should be returned
	//   - conflict = true when there is an active conflicting pending change and 409 should be returned
	//   - proceed = true when operation can be applied immediately (not guarded)
	//   - err on unexpected failures
	CheckAndMaybeCreatePending(
		ctx context.Context,
		in GuardEngineInput,
	) (pendingChange *domain.PendingChange, conflict bool, proceed bool, err error)
}
