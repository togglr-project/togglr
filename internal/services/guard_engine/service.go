package guard_engine

import (
	"context"
	"fmt"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// Service is a default implementation of contract.GuardEngine.
// It encapsulates guard checks and pending change creation logic.
type Service struct {
	guardService          contract.GuardService
	pendingChangesUseCase contract.PendingChangesUseCase
}

func New(
	guardService contract.GuardService,
	pendingChangesUseCase contract.PendingChangesUseCase,
) *Service {
	return &Service{
		guardService:          guardService,
		pendingChangesUseCase: pendingChangesUseCase,
	}
}

var _ contract.GuardEngine = (*Service)(nil)

func (s *Service) CheckAndMaybeCreatePending(
	ctx context.Context,
	in contract.GuardEngineInput,
) (*domain.PendingChange, bool, bool, error) {
	// If feature is not guarded, proceed normally
	isGuarded, err := s.guardService.IsFeatureGuarded(ctx, in.FeatureID)
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
	hasConflict, err := s.pendingChangesUseCase.CheckEntityConflict(ctx, entities)
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

	pc, err := s.pendingChangesUseCase.Create(
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

	// Set SingleUserProject meta-flag analogous to features_service implementation
	activeUserCount, err := s.pendingChangesUseCase.GetProjectActiveUserCount(ctx, in.ProjectID)
	if err != nil {
		return nil, false, false, fmt.Errorf("get project active user count: %w", err)
	}

	if activeUserCount == 1 {
		pc.Change.Meta.SingleUserProject = true
	}

	return &pc, false, false, nil
}
