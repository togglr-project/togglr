//nolint:nestif // This is a complex service implementation.
package guard_engine

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// EntityChanges represents changes for a specific entity type.
type EntityChanges struct {
	EntityType string
	EntityID   string
	Changes    map[string]domain.ChangeValue
}

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

// CheckGuardedOperation is a high-level method that automatically determines
// entity types and computes changes by comparing old and new entities.
func (s *Service) CheckGuardedOperation(
	ctx context.Context,
	req contract.GuardRequest,
) (*domain.PendingChange, bool, bool, error) {
	// Determine entity types and compute changes
	entityChangesList, err := s.determineEntityTypeAndChanges(req.OldEntity, req.NewEntity, req.Action)
	if err != nil {
		return nil, false, false, fmt.Errorf("determine entity type and changes: %w", err)
	}

	// If no changes detected, proceed normally
	if len(entityChangesList) == 0 {
		return nil, false, true, nil
	}

	// Use the existing guard engine with the computed changes
	return s.checkAndMaybeCreatePending(
		ctx,
		GuardEngineInput{
			ProjectID:         req.ProjectID,
			EnvironmentID:     req.EnvironmentID,
			FeatureID:         req.FeatureID,
			Reason:            req.Reason,
			Origin:            req.Origin,
			EntityChangesList: entityChangesList,
			Action:            req.Action,
		},
	)
}

// BuildChangeDiff computes changes between two entities using reflection and editable tags.
func (s *Service) BuildChangeDiff(oldEntity, newEntity any) map[string]domain.ChangeValue {
	return BuildChangeDiff(oldEntity, newEntity)
}

// BuildInsertChanges identifies fields that should be included in insert operations.
func (s *Service) BuildInsertChanges(entity any) map[string]domain.ChangeValue {
	return BuildInsertChanges(entity)
}

func (s *Service) checkAndMaybeCreatePending(
	ctx context.Context,
	in GuardEngineInput,
) (*domain.PendingChange, bool, bool, error) {
	// If feature is not guarded, proceed normally
	isGuarded, err := s.guardService.IsFeatureGuarded(ctx, in.FeatureID)
	if err != nil {
		return nil, false, false, fmt.Errorf("check feature guarded: %w", err)
	}
	if !isGuarded {
		return nil, false, true, nil
	}

	// Build entities from the changes list
	entities := []domain.EntityChange{
		{
			Entity:   string(domain.EntityFeature),
			EntityID: in.FeatureID.String(),
			Action:   domain.EntityActionUpdate,
			Changes:  map[string]domain.ChangeValue{},
		},
	}

	// Add entities from the changes list
	for _, entityChanges := range in.EntityChangesList {
		entities = append(entities, domain.EntityChange{
			Entity:   entityChanges.EntityType,
			EntityID: entityChanges.EntityID,
			Action:   in.Action,
			Changes:  entityChanges.Changes,
		})
	}

	// Get project active user count before creating pending change
	activeUserCount, err := s.pendingChangesUseCase.GetProjectActiveUserCount(ctx, in.ProjectID)
	if err != nil {
		return nil, false, false, fmt.Errorf("get project active user count: %w", err)
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
	// If exactly 1 active user, treat as a single-user project (enables auto-approve)
	if activeUserCount == 1 {
		pc.Change.Meta.SingleUserProject = true
	}

	return &pc, false, false, nil
}

// determineEntityTypeAndChanges determines entity types and computes changes
// by comparing old and new entities.
//
//nolint:gocritic // not critical
func (s *Service) determineEntityTypeAndChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) ([]EntityChanges, error) {
	// Determine an entity type from the old entity (or new entity if old is nil)
	var entity any
	if oldEntity != nil {
		entity = oldEntity
	} else if newEntity != nil {
		entity = newEntity
	} else {
		return nil, errors.New("both old and new entities are nil")
	}

	// Get entity type and ID
	entityType, entityID, err := s.getEntityTypeAndID(entity)
	if err != nil {
		return nil, fmt.Errorf("get entity type and ID: %w", err)
	}

	// Compute changes based on an entity type
	switch entityType {
	case string(domain.EntityFeature):
		return s.computeFeatureChanges(oldEntity, newEntity, action)
	case string(domain.EntityFeatureParams):
		changes, err := s.computeFeatureParamsChanges(oldEntity, newEntity, action)
		if err != nil {
			return nil, err
		}

		return []EntityChanges{{
			EntityType: entityType,
			EntityID:   entityID,
			Changes:    changes,
		}}, nil
	case string(domain.EntityRule):
		return s.computeRuleChanges(oldEntity, newEntity, action)
	case string(domain.EntityFlagVariant):
		return s.computeFlagVariantChanges(oldEntity, newEntity, action)
	case string(domain.EntityFeatureSchedule):
		return s.computeFeatureScheduleChanges(oldEntity, newEntity, action)
	case string(domain.EntityFeatureTag):
		return s.computeFeatureTagChanges(oldEntity, newEntity, action)
	default:
		return nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}
}

// getEntityTypeAndID extracts entity type and ID from an entity.
func (s *Service) getEntityTypeAndID(entity any) (entityType, entityID string, err error) {
	switch e := entity.(type) {
	case *domain.Feature:
		return string(domain.EntityFeature), string(e.ID), nil
	case domain.Feature:
		return string(domain.EntityFeature), string(e.ID), nil
	case *domain.Rule:
		return string(domain.EntityRule), string(e.ID), nil
	case domain.Rule:
		return string(domain.EntityRule), string(e.ID), nil
	case *domain.FlagVariant:
		return string(domain.EntityFlagVariant), string(e.ID), nil
	case domain.FlagVariant:
		return string(domain.EntityFlagVariant), string(e.ID), nil
	case *domain.FeatureSchedule:
		return string(domain.EntityFeatureSchedule), string(e.ID), nil
	case domain.FeatureSchedule:
		return string(domain.EntityFeatureSchedule), string(e.ID), nil
	case *domain.FeatureParams:
		return string(domain.EntityFeatureParams), string(e.FeatureID), nil
	case domain.FeatureParams:
		return string(domain.EntityFeatureParams), string(e.FeatureID), nil
	// FeatureTag is not a separate entity type, it's a relationship
	// We'll handle it differently in computeFeatureTagChanges
	default:
		// Check if it's a FeatureTag relationship struct
		if entityType, entityID, err := s.handleFeatureTagStruct(entity); err == nil {
			return entityType, entityID, nil
		}

		return "", "", fmt.Errorf("unknown entity type: %T", entity)
	}
}

// handleFeatureTagStruct handles anonymous structs used for FeatureTag relationships.
// It checks if the struct has FeatureID and TagID fields and returns appropriate entity type and ID.
func (s *Service) handleFeatureTagStruct(entity any) (entityType, entityID string, err error) {
	// Use reflection to check if this is a FeatureTag relationship struct
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", "", errors.New("not a struct")
	}

	// Check if it has FeatureID and TagID fields
	featureIDField := val.FieldByName("FeatureID")
	tagIDField := val.FieldByName("TagID")

	if !featureIDField.IsValid() || !tagIDField.IsValid() {
		return "", "", errors.New("not a FeatureTag struct")
	}

	// Get the values
	featureID := featureIDField.String()
	tagID := tagIDField.String()

	if featureID == "" || tagID == "" {
		return "", "", errors.New("empty FeatureID or TagID")
	}

	// For FeatureTag relationships, we use the TagID as the entity ID
	// and return EntityFeatureTag as the entity type
	return string(domain.EntityFeatureTag), tagID, nil
}

// computeFeatureChanges computes changes for feature entities.
func (s *Service) computeFeatureChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) ([]EntityChanges, error) {
	switch action {
	case domain.EntityActionUpdate:
		oldFeature, err := s.convertToFeaturePtr(oldEntity)
		if err != nil {
			return nil, fmt.Errorf("old entity: %w", err)
		}
		newFeature, err := s.convertToFeaturePtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		var result []EntityChanges

		// Changes for the base features table (BasicFeature)
		basicFeatureChanges := BuildChangeDiff(oldFeature, newFeature)
		if len(basicFeatureChanges) > 0 {
			result = append(result, EntityChanges{
				EntityType: string(domain.EntityFeature),
				EntityID:   string(newFeature.ID),
				Changes:    basicFeatureChanges,
			})
		}

		// Changes for the feature_params table (enabled, default_value)
		// Use ConvertToFeatureParams method for automatic conversion
		oldParams := oldFeature.ConvertToFeatureParams()
		newParams := newFeature.ConvertToFeatureParams()

		paramsChanges := BuildChangeDiff(&oldParams, &newParams)
		if len(paramsChanges) > 0 {
			result = append(result, EntityChanges{
				EntityType: string(domain.EntityFeatureParams),
				EntityID:   string(newFeature.ID),
				Changes:    paramsChanges,
			})
		}

		return result, nil
	case domain.EntityActionInsert:
		// Insert action is not applicable for feature changes computation
		return nil, errors.New("insert action not supported for feature changes")
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for feature: %s", action)
	}
}

// computeFeatureParamsChanges computes changes for feature params entities.
func (s *Service) computeFeatureParamsChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) (map[string]domain.ChangeValue, error) {
	switch action {
	case domain.EntityActionUpdate:
		oldParams, err := s.convertToFeatureParamsPtr(oldEntity)
		if err != nil {
			return nil, fmt.Errorf("old entity: %w", err)
		}
		newParams, err := s.convertToFeatureParamsPtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		return BuildChangeDiff(oldParams, newParams), nil
	case domain.EntityActionInsert:
		newParams, err := s.convertToFeatureParamsPtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		return BuildInsertChanges(newParams), nil
	case domain.EntityActionDelete:
		return map[string]domain.ChangeValue{}, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for feature params: %s", action)
	}
}

// computeRuleChanges computes changes for rule entities.
func (s *Service) computeRuleChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) ([]EntityChanges, error) {
	switch action {
	case domain.EntityActionUpdate:
		oldRule, err := s.convertToRulePtr(oldEntity)
		if err != nil {
			return nil, fmt.Errorf("old entity: %w", err)
		}
		newRule, err := s.convertToRulePtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		changes := BuildChangeDiff(oldRule, newRule)
		if len(changes) == 0 {
			return nil, nil
		}

		return []EntityChanges{{
			EntityType: string(domain.EntityRule),
			EntityID:   string(newRule.ID),
			Changes:    changes,
		}}, nil
	case domain.EntityActionInsert:
		newRule, err := s.convertToRulePtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		changes := BuildInsertChanges(newRule)

		return []EntityChanges{{
			EntityType: string(domain.EntityRule),
			EntityID:   string(newRule.ID),
			Changes:    changes,
		}}, nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for rule: %s", action)
	}
}

// computeFlagVariantChanges computes changes for flag variant entities.
func (s *Service) computeFlagVariantChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) ([]EntityChanges, error) {
	switch action {
	case domain.EntityActionUpdate:
		oldVariant, err := s.convertToFlagVariantPtr(oldEntity)
		if err != nil {
			return nil, fmt.Errorf("old entity: %w", err)
		}
		newVariant, err := s.convertToFlagVariantPtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		changes := BuildChangeDiff(oldVariant, newVariant)
		if len(changes) == 0 {
			return nil, nil
		}

		return []EntityChanges{{
			EntityType: string(domain.EntityFlagVariant),
			EntityID:   string(newVariant.ID),
			Changes:    changes,
		}}, nil
	case domain.EntityActionInsert:
		newVariant, err := s.convertToFlagVariantPtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		changes := BuildInsertChanges(newVariant)

		return []EntityChanges{{
			EntityType: string(domain.EntityFlagVariant),
			EntityID:   string(newVariant.ID),
			Changes:    changes,
		}}, nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for flag variant: %s", action)
	}
}

// computeFeatureScheduleChanges computes changes for feature schedule entities.
func (s *Service) computeFeatureScheduleChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) ([]EntityChanges, error) {
	switch action {
	case domain.EntityActionUpdate:
		oldSchedule, err := s.convertToFeatureSchedulePtr(oldEntity)
		if err != nil {
			return nil, fmt.Errorf("old entity: %w", err)
		}
		newSchedule, err := s.convertToFeatureSchedulePtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		changes := BuildChangeDiff(oldSchedule, newSchedule)
		if len(changes) == 0 {
			return nil, nil
		}

		return []EntityChanges{{
			EntityType: string(domain.EntityFeatureSchedule),
			EntityID:   string(newSchedule.ID),
			Changes:    changes,
		}}, nil
	case domain.EntityActionInsert:
		newSchedule, err := s.convertToFeatureSchedulePtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		changes := BuildInsertChanges(newSchedule)

		return []EntityChanges{{
			EntityType: string(domain.EntityFeatureSchedule),
			EntityID:   string(newSchedule.ID),
			Changes:    changes,
		}}, nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for feature schedule: %s", action)
	}
}

// computeFeatureTagChanges computes changes for feature tag entities.
// For feature tags, we expect a struct with FeatureID and TagID fields.
//
//nolint:staticcheck // This is a workaround for a staticcheck bug
func (s *Service) computeFeatureTagChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) ([]EntityChanges, error) {
	switch action {
	case domain.EntityActionInsert, domain.EntityActionDelete:
		// For feature tags, we need to extract feature_id and tag_id
		var featureID, tagID string

		// Try to extract from an old entity first, then a new entity
		if oldEntity != nil {
			if featureID, tagID = s.extractFeatureTagIDs(oldEntity); featureID != "" && tagID != "" {
				// Found in old entity
			}
		}
		if newEntity != nil && (featureID == "" || tagID == "") {
			if featureID, tagID = s.extractFeatureTagIDs(newEntity); featureID != "" && tagID != "" {
				// Found in a new entity
			}
		}

		if featureID == "" || tagID == "" {
			return nil, errors.New("feature_id and tag_id are required for feature tag changes")
		}

		changes := map[string]domain.ChangeValue{
			"feature_id": {New: featureID},
			"tag_id":     {New: tagID},
		}

		return []EntityChanges{{
			EntityType: string(domain.EntityFeatureTag),
			EntityID:   featureID,
			Changes:    changes,
		}}, nil
	case domain.EntityActionUpdate:
		// For update, we can use the generic function
		oldTag, err := s.convertToFeatureTagPtr(oldEntity)
		if err != nil {
			return nil, fmt.Errorf("old entity: %w", err)
		}
		newTag, err := s.convertToFeatureTagPtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		changes := BuildChangeDiff(oldTag, newTag)
		if len(changes) == 0 {
			return nil, nil
		}

		return []EntityChanges{{
			EntityType: string(domain.EntityFeatureTag),
			EntityID:   string(newTag.FeatureID),
			Changes:    changes,
		}}, nil
	default:
		return nil, fmt.Errorf("unsupported action for feature tag: %s", action)
	}
}

// extractFeatureTagIDs extracts FeatureID and TagID from an entity.
// This handles both struct types and map types.
func (s *Service) extractFeatureTagIDs(entity any) (featureID, tagID string) {
	// Try struct type first
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		// Look for FeatureID field
		if featureField := v.FieldByName("FeatureID"); featureField.IsValid() {
			featureID = featureField.String()
		}
		// Look for TagID field
		if tagField := v.FieldByName("TagID"); tagField.IsValid() {
			tagID = tagField.String()
		}
	}

	// Try a map type
	if m, ok := entity.(map[string]any); ok {
		if f, exists := m["feature_id"]; exists {
			if fStr, ok := f.(string); ok {
				featureID = fStr
			}
		}
		if t, exists := m["tag_id"]; exists {
			if tStr, ok := t.(string); ok {
				tagID = tStr
			}
		}
	}

	return featureID, tagID
}

// Helper functions to convert any entity to a pointer type for consistent handling

func (s *Service) convertToFeaturePtr(entity any) (*domain.Feature, error) {
	switch e := entity.(type) {
	case *domain.Feature:
		return e, nil
	case domain.Feature:
		return &e, nil
	default:
		return nil, fmt.Errorf("expected *domain.Feature or domain.Feature, got %T", entity)
	}
}

func (s *Service) convertToFeatureParamsPtr(entity any) (*domain.FeatureParams, error) {
	switch e := entity.(type) {
	case *domain.FeatureParams:
		return e, nil
	case domain.FeatureParams:
		return &e, nil
	default:
		return nil, fmt.Errorf("expected *domain.FeatureParams or domain.FeatureParams, got %T", entity)
	}
}

func (s *Service) convertToRulePtr(entity any) (*domain.Rule, error) {
	switch e := entity.(type) {
	case *domain.Rule:
		return e, nil
	case domain.Rule:
		return &e, nil
	default:
		return nil, fmt.Errorf("expected *domain.Rule or domain.Rule, got %T", entity)
	}
}

func (s *Service) convertToFlagVariantPtr(entity any) (*domain.FlagVariant, error) {
	switch e := entity.(type) {
	case *domain.FlagVariant:
		return e, nil
	case domain.FlagVariant:
		return &e, nil
	default:
		return nil, fmt.Errorf("expected *domain.FlagVariant or domain.FlagVariant, got %T", entity)
	}
}

func (s *Service) convertToFeatureSchedulePtr(entity any) (*domain.FeatureSchedule, error) {
	switch e := entity.(type) {
	case *domain.FeatureSchedule:
		return e, nil
	case domain.FeatureSchedule:
		return &e, nil
	default:
		return nil, fmt.Errorf("expected *domain.FeatureSchedule or domain.FeatureSchedule, got %T", entity)
	}
}

func (s *Service) convertToFeatureTagPtr(entity any) (*domain.FeatureTags, error) {
	switch e := entity.(type) {
	case *domain.FeatureTags:
		return e, nil
	case domain.FeatureTags:
		return &e, nil
	default:
		return nil, fmt.Errorf("expected *domain.FeatureTags or domain.FeatureTags, got %T", entity)
	}
}
