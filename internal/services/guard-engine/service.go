package guard_engine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

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

// CheckGuardedOperation is a high-level method that automatically determines
// an entity type and computes changes by comparing old and new entities.
func (s *Service) CheckGuardedOperation(
	ctx context.Context,
	req contract.GuardRequest,
) (*domain.PendingChange, bool, bool, error) {
	// Determine an entity type and compute changes
	entityType, entityID, changes, err := s.determineEntityTypeAndChanges(req.OldEntity, req.NewEntity, req.Action)
	if err != nil {
		return nil, false, false, fmt.Errorf("determine entity type and changes: %w", err)
	}

	// If no changes detected, proceed normally
	if len(changes) == 0 {
		return nil, false, true, nil
	}

	// Use the existing guard engine with the computed changes
	return s.checkAndMaybeCreatePending(
		ctx,
		GuardEngineInput{
			ProjectID:       req.ProjectID,
			EnvironmentID:   req.EnvironmentID,
			FeatureID:       req.FeatureID,
			Reason:          req.Reason,
			Origin:          req.Origin,
			PrimaryEntity:   entityType,
			PrimaryEntityID: entityID,
			Action:          req.Action,
			ExtraChanges:    changes,
		},
	)
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

// determineEntityTypeAndChanges determines the entity type and computes changes
// by comparing old and new entities.
func (s *Service) determineEntityTypeAndChanges(
	oldEntity, newEntity any,
	action domain.EntityAction,
) (entityType, entityID string, changes map[string]domain.ChangeValue, err error) {
	// Determine entity type from the old entity (or new entity if old is nil)
	var entity any
	if oldEntity != nil {
		entity = oldEntity
	} else if newEntity != nil {
		entity = newEntity
	} else {
		return "", "", nil, errors.New("both old and new entities are nil")
	}

	// Get entity type and ID
	entityType, entityID, err = s.getEntityTypeAndID(entity)
	if err != nil {
		return "", "", nil, fmt.Errorf("get entity type and ID: %w", err)
	}

	// Compute changes based on entity type
	switch entityType {
	case string(domain.EntityFeature):
		changes, err = s.computeFeatureChanges(oldEntity, newEntity, action)
	case string(domain.EntityFeatureParams):
		changes, err = s.computeFeatureParamsChanges(oldEntity, newEntity, action)
	case string(domain.EntityRule):
		changes, err = s.computeRuleChanges(oldEntity, newEntity, action)
	case string(domain.EntityFlagVariant):
		changes, err = s.computeFlagVariantChanges(oldEntity, newEntity, action)
	case string(domain.EntityFeatureSchedule):
		changes, err = s.computeFeatureScheduleChanges(oldEntity, newEntity, action)
	case string(domain.EntityFeatureTag):
		changes, err = s.computeFeatureTagChanges(oldEntity, newEntity, action)
	default:
		return "", "", nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}

	return entityType, entityID, changes, err
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
		return "", "", fmt.Errorf("unknown entity type: %T", entity)
	}
}

// computeFeatureChanges computes changes for feature entities.
func (s *Service) computeFeatureChanges(oldEntity, newEntity any, action domain.EntityAction) (map[string]domain.ChangeValue, error) {
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

		return s.buildFeatureChangeDiff(oldFeature, newFeature), nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for feature: %s", action)
	}
}

// computeFeatureParamsChanges computes changes for feature params entities.
func (s *Service) computeFeatureParamsChanges(oldEntity, newEntity any, action domain.EntityAction) (map[string]domain.ChangeValue, error) {
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

		return s.buildFeatureParamsChangeDiff(oldParams, newParams), nil
	case domain.EntityActionInsert:
		newParams, err := s.convertToFeatureParamsPtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		return s.buildFeatureParamsInsertChanges(newParams), nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for feature params: %s", action)
	}
}

// computeRuleChanges computes changes for rule entities.
func (s *Service) computeRuleChanges(oldEntity, newEntity any, action domain.EntityAction) (map[string]domain.ChangeValue, error) {
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

		return s.buildRuleChangeDiff(oldRule, newRule), nil
	case domain.EntityActionInsert:
		newRule, err := s.convertToRulePtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		return s.buildRuleInsertChanges(newRule), nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for rule: %s", action)
	}
}

// computeFlagVariantChanges computes changes for flag variant entities.
func (s *Service) computeFlagVariantChanges(oldEntity, newEntity any, action domain.EntityAction) (map[string]domain.ChangeValue, error) {
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

		return s.buildFlagVariantChangeDiff(oldVariant, newVariant), nil
	case domain.EntityActionInsert:
		newVariant, err := s.convertToFlagVariantPtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		return s.buildFlagVariantInsertChanges(newVariant), nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for flag variant: %s", action)
	}
}

// computeFeatureScheduleChanges computes changes for feature schedule entities.
func (s *Service) computeFeatureScheduleChanges(oldEntity, newEntity any, action domain.EntityAction) (map[string]domain.ChangeValue, error) {
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

		return s.buildFeatureScheduleChangeDiff(oldSchedule, newSchedule), nil
	case domain.EntityActionInsert:
		newSchedule, err := s.convertToFeatureSchedulePtr(newEntity)
		if err != nil {
			return nil, fmt.Errorf("new entity: %w", err)
		}

		return s.buildFeatureScheduleInsertChanges(newSchedule), nil
	case domain.EntityActionDelete:
		return nil, nil // No changes needed for delete
	default:
		return nil, fmt.Errorf("unsupported action for feature schedule: %s", action)
	}
}

// computeFeatureTagChanges computes changes for feature tag entities.
// For feature tags, we expect a struct with FeatureID and TagID fields.
func (s *Service) computeFeatureTagChanges(oldEntity, newEntity any, action domain.EntityAction) (map[string]domain.ChangeValue, error) {
	switch action {
	case domain.EntityActionInsert, domain.EntityActionDelete:
		// For feature tags, we need to extract feature_id and tag_id
		var featureID, tagID string

		// Try to extract from old entity first, then new entity
		if oldEntity != nil {
			if featureID, tagID = s.extractFeatureTagIDs(oldEntity); featureID != "" && tagID != "" {
				// Found in old entity
			}
		}
		if newEntity != nil && (featureID == "" || tagID == "") {
			if featureID, tagID = s.extractFeatureTagIDs(newEntity); featureID != "" && tagID != "" {
				// Found in new entity
			}
		}

		if featureID == "" || tagID == "" {
			return nil, errors.New("feature_id and tag_id are required for feature tag changes")
		}

		return map[string]domain.ChangeValue{
			"feature_id": {New: featureID},
			"tag_id":     {New: tagID},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported action for feature tag: %s", action)
	}
}

// extractFeatureTagIDs extracts FeatureID and TagID from an entity.
// This handles both struct types and map types.
func (s *Service) extractFeatureTagIDs(entity any) (featureID, tagID string) {
	// Try struct type first
	if v := reflect.ValueOf(entity); v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v := reflect.ValueOf(entity); v.Kind() == reflect.Struct {
		// Look for FeatureID field
		if featureField := v.FieldByName("FeatureID"); featureField.IsValid() {
			featureID = featureField.String()
		}
		// Look for TagID field
		if tagField := v.FieldByName("TagID"); tagField.IsValid() {
			tagID = tagField.String()
		}
	}

	// Try map type
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

// Helper methods for building change diffs

func (s *Service) buildFeatureChangeDiff(old, new *domain.Feature) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if old.Name != new.Name {
		changes["name"] = domain.ChangeValue{Old: old.Name, New: new.Name}
	}
	if old.Description != new.Description {
		changes["description"] = domain.ChangeValue{Old: old.Description, New: new.Description}
	}
	if old.Enabled != new.Enabled {
		changes["enabled"] = domain.ChangeValue{Old: old.Enabled, New: new.Enabled}
	}
	if old.DefaultValue != new.DefaultValue {
		changes["default_value"] = domain.ChangeValue{Old: old.DefaultValue, New: new.DefaultValue}
	}

	return changes
}

func (s *Service) buildFeatureParamsChangeDiff(old, new *domain.FeatureParams) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if old.Enabled != new.Enabled {
		changes["enabled"] = domain.ChangeValue{Old: old.Enabled, New: new.Enabled}
	}
	if old.DefaultValue != new.DefaultValue {
		changes["default_value"] = domain.ChangeValue{Old: old.DefaultValue, New: new.DefaultValue}
	}

	return changes
}

func (s *Service) buildFeatureParamsInsertChanges(params *domain.FeatureParams) map[string]domain.ChangeValue {
	return map[string]domain.ChangeValue{
		"feature_id":     {New: string(params.FeatureID)},
		"environment_id": {New: params.EnvironmentID},
		"enabled":        {New: params.Enabled},
		"default_value":  {New: params.DefaultValue},
	}
}

func (s *Service) buildRuleChangeDiff(old, new *domain.Rule) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if old.IsCustomized != new.IsCustomized {
		changes["is_customized"] = domain.ChangeValue{Old: old.IsCustomized, New: new.IsCustomized}
	}
	if old.Action != new.Action {
		changes["action"] = domain.ChangeValue{Old: old.Action.String(), New: new.Action.String()}
	}
	if old.Priority != new.Priority {
		changes["priority"] = domain.ChangeValue{Old: old.Priority, New: new.Priority}
	}

	// Handle optional fields
	if s.ruleFlagVariantIDChanged(old.FlagVariantID, new.FlagVariantID) {
		changes["flag_variant_id"] = domain.ChangeValue{
			Old: s.ruleFlagVariantIDString(old.FlagVariantID),
			New: s.ruleFlagVariantIDString(new.FlagVariantID),
		}
	}
	if s.ruleSegmentIDChanged(old.SegmentID, new.SegmentID) {
		changes["segment_id"] = domain.ChangeValue{
			Old: s.ruleSegmentIDString(old.SegmentID),
			New: s.ruleSegmentIDString(new.SegmentID),
		}
	}

	return changes
}

func (s *Service) buildRuleInsertChanges(rule *domain.Rule) map[string]domain.ChangeValue {
	changes := map[string]domain.ChangeValue{
		"project_id":    {New: string(rule.ProjectID)},
		"feature_id":    {New: string(rule.FeatureID)},
		"is_customized": {New: rule.IsCustomized},
		"action":        {New: rule.Action.String()},
		"priority":      {New: rule.Priority},
	}

	if rule.FlagVariantID != nil {
		changes["flag_variant_id"] = domain.ChangeValue{New: s.ruleFlagVariantIDString(rule.FlagVariantID)}
	}
	if rule.SegmentID != nil {
		changes["segment_id"] = domain.ChangeValue{New: s.ruleSegmentIDString(rule.SegmentID)}
	}

	return changes
}

func (s *Service) buildFlagVariantChangeDiff(old, new *domain.FlagVariant) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if old.Name != new.Name {
		changes["name"] = domain.ChangeValue{Old: old.Name, New: new.Name}
	}
	if old.RolloutPercent != new.RolloutPercent {
		changes["rollout_percent"] = domain.ChangeValue{Old: old.RolloutPercent, New: new.RolloutPercent}
	}

	return changes
}

func (s *Service) buildFlagVariantInsertChanges(variant *domain.FlagVariant) map[string]domain.ChangeValue {
	return map[string]domain.ChangeValue{
		"project_id":      {New: string(variant.ProjectID)},
		"feature_id":      {New: string(variant.FeatureID)},
		"name":            {New: variant.Name},
		"rollout_percent": {New: variant.RolloutPercent},
	}
}

func (s *Service) buildFeatureScheduleChangeDiff(old, new *domain.FeatureSchedule) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if s.timePtrChanged(old.StartsAt, new.StartsAt) {
		changes["starts_at"] = domain.ChangeValue{
			Old: s.timePtrString(old.StartsAt),
			New: s.timePtrString(new.StartsAt),
		}
	}
	if s.timePtrChanged(old.EndsAt, new.EndsAt) {
		changes["ends_at"] = domain.ChangeValue{
			Old: s.timePtrString(old.EndsAt),
			New: s.timePtrString(new.EndsAt),
		}
	}
	if s.stringPtrChanged(old.CronExpr, new.CronExpr) {
		changes["cron_expr"] = domain.ChangeValue{
			Old: s.stringPtrString(old.CronExpr),
			New: s.stringPtrString(new.CronExpr),
		}
	}
	if s.durationPtrChanged(old.CronDuration, new.CronDuration) {
		changes["cron_duration"] = domain.ChangeValue{
			Old: s.durationPtrString(old.CronDuration),
			New: s.durationPtrString(new.CronDuration),
		}
	}
	if old.Timezone != new.Timezone {
		changes["timezone"] = domain.ChangeValue{Old: old.Timezone, New: new.Timezone}
	}
	if old.Action != new.Action {
		changes["action"] = domain.ChangeValue{Old: old.Action.String(), New: new.Action.String()}
	}

	return changes
}

func (s *Service) buildFeatureScheduleInsertChanges(schedule *domain.FeatureSchedule) map[string]domain.ChangeValue {
	changes := map[string]domain.ChangeValue{
		"project_id":     {New: string(schedule.ProjectID)},
		"feature_id":     {New: string(schedule.FeatureID)},
		"environment_id": {New: schedule.EnvironmentID},
		"timezone":       {New: schedule.Timezone},
		"action":         {New: schedule.Action.String()},
	}

	if schedule.StartsAt != nil {
		changes["starts_at"] = domain.ChangeValue{New: s.timePtrString(schedule.StartsAt)}
	}
	if schedule.EndsAt != nil {
		changes["ends_at"] = domain.ChangeValue{New: s.timePtrString(schedule.EndsAt)}
	}
	if schedule.CronExpr != nil {
		changes["cron_expr"] = domain.ChangeValue{New: s.stringPtrString(schedule.CronExpr)}
	}
	if schedule.CronDuration != nil {
		changes["cron_duration"] = domain.ChangeValue{New: s.durationPtrString(schedule.CronDuration)}
	}

	return changes
}

// Helper methods for comparing optional fields

func (s *Service) ruleFlagVariantIDChanged(old, new *domain.FlagVariantID) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}

	return *old != *new
}

func (s *Service) ruleFlagVariantIDString(id *domain.FlagVariantID) interface{} {
	if id == nil {
		return nil
	}

	return string(*id)
}

func (s *Service) ruleSegmentIDChanged(old, new *domain.SegmentID) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}

	return *old != *new
}

func (s *Service) ruleSegmentIDString(id *domain.SegmentID) interface{} {
	if id == nil {
		return nil
	}

	return string(*id)
}

func (s *Service) timePtrChanged(old, new *time.Time) bool {
	return s.timePtrString(old) != s.timePtrString(new)
}

func (s *Service) timePtrString(p *time.Time) interface{} {
	if p == nil {
		return nil
	}

	return p.Format(time.RFC3339)
}

func (s *Service) stringPtrChanged(old, new *string) bool {
	return s.stringPtrString(old) != s.stringPtrString(new)
}

func (s *Service) stringPtrString(p *string) interface{} {
	if p == nil {
		return nil
	}

	return *p
}

func (s *Service) durationPtrChanged(old, new *time.Duration) bool {
	return s.durationPtrString(old) != s.durationPtrString(new)
}

func (s *Service) durationPtrString(p *time.Duration) interface{} {
	if p == nil {
		return nil
	}

	return p.String()
}

// Helper functions to convert any entity to pointer type for consistent handling

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
