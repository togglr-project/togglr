package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

// GuardRequest represents a high-level request for guard checking.
// The engine will automatically determine entity type and compute changes.
type GuardRequest struct {
	ProjectID     domain.ProjectID
	EnvironmentID domain.EnvironmentID
	FeatureID     domain.FeatureID
	Reason        string
	Origin        string
	Action        domain.EntityAction

	// Entity data - the engine will determine type and compute changes
	OldEntity any
	NewEntity any // nil for delete operations
}

// GuardEngine encapsulates guard checks and pending change creation.
// It centralizes the guarded workflow so that outer layers (API) don't need to
// know how to assemble payloads or handle conflicts.
type GuardEngine interface {
	// CheckGuardedOperation is a high-level method that automatically determines
	// entity type and computes changes by comparing old and new entities.
	//
	// The method will:
	// 1. Determine entity type from the provided entities
	// 2. Compute changes by comparing old and new entities
	// 3. Check if the feature is guarded
	// 4. Create pending change if needed
	//
	// Returns:
	//   - pendingChange (non-nil) when a pending change was created and 202 should be returned
	//   - conflict = true when there is an active conflicting pending change and 409 should be returned
	//   - proceed = true when operation can be applied immediately (not guarded)
	//   - err on unexpected failures
	CheckGuardedOperation(
		ctx context.Context,
		req GuardRequest,
	) (pendingChange *domain.PendingChange, conflict bool, proceed bool, err error)
}
