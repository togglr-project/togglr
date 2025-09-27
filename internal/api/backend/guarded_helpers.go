package apibackend

import (
	"context"
	"fmt"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// GuardPendingInput groups parameters for the guarded pending change helper to avoid long arg lists.
type GuardPendingInput struct {
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

// guardCheckAndMaybeCreatePending encapsulates guarded logic for feature-scoped mutations.
// It checks whether the feature is guarded, detects conflicts on the feature entity,
// and if needed creates a pending change and returns its response.
// Returns:
//   - pendingResp (non-nil) when a pending change was created and 202 should be returned
//   - conflict = true when there is an active conflicting pending change and 409 should be returned
//   - proceed = true when operation can be applied immediately (not guarded)
//   - err on unexpected failures
func (r *RestAPI) guardCheckAndMaybeCreatePending(
	ctx context.Context,
	in GuardPendingInput,
) (pendingResp *generatedapi.PendingChangeResponse, conflict bool, proceed bool, err error) {
	// If feature is not guarded, proceed normally
	isGuarded, err := r.guardService.IsFeatureGuarded(ctx, in.FeatureID)
	if err != nil {
		return nil, false, false, fmt.Errorf("check feature guarded: %w", err)
	}
	if !isGuarded {
		return nil, false, true, nil
	}

	// Build entities: always include the feature entity to serialize/lock, and optionally the primary entity
	entities := []domain.EntityChange{
		{
			Entity:   string(domain.EntityFeature),
			EntityID: in.FeatureID.String(),
			Action:   domain.EntityActionUpdate,
			Changes:  map[string]domain.ChangeValue{},
		},
	}
	if in.PrimaryEntity != "" {
		changes := map[string]domain.ChangeValue{}
		for k, v := range in.ExtraChanges {
			changes[k] = v
		}
		entities = append(entities, domain.EntityChange{
			Entity:   in.PrimaryEntity,
			EntityID: in.PrimaryEntityID,
			Action:   in.Action,
			Changes:  changes,
		})
	}

	// Check conflicts for the entities
	hasConflict, err := r.pendingChangesUseCase.CheckEntityConflict(ctx, entities)
	if err != nil {
		return nil, false, false, fmt.Errorf("check entity conflict: %w", err)
	}
	if hasConflict {
		return nil, true, false, nil
	}

	// Create pending change
	requestedBy := appcontext.Username(ctx)
	requestUserID := appcontext.UserID(ctx)
	var requestUserIDPtr *int
	if requestUserID != 0 {
		v := int(requestUserID)
		requestUserIDPtr = &v
	}

	payload := domain.PendingChangePayload{
		Entities: entities,
		Meta: domain.PendingChangeMeta{
			Reason: in.Reason,
			Client: "ui",
			Origin: in.Origin,
		},
	}

	pc, err := r.pendingChangesUseCase.Create(
		ctx,
		in.ProjectID,
		in.EnvironmentID,
		requestedBy,
		requestUserIDPtr,
		payload,
	)
	if err != nil {
		return nil, false, false, fmt.Errorf("create pending change: %w", err)
	}

	resp := convertPendingChangeToResponse(&pc)

	return &resp, false, false, nil
}
