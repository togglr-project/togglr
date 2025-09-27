package pending_changes

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Service struct {
	txManager            db.TxManager
	pendingChangesRepo   contract.PendingChangesRepository
	projectApproversRepo contract.ProjectApproversRepository
	projectSettingsRepo  contract.ProjectSettingsRepository
	guardService         contract.GuardService
	featuresRepo         contract.FeaturesRepository
	featureParamsRepo    contract.FeatureParamsRepository
	// Added repositories to apply changes for child entities
	rulesRepo          contract.RulesRepository
	flagVariantsRepo   contract.FlagVariantsRepository
	schedulesRepo      contract.FeatureSchedulesRepository
	featureTagsRepo    contract.FeatureTagsRepository
	auditLogRepo       contract.AuditLogRepository
	usersUseCase       contract.UsersUseCase
	permissionsService contract.PermissionsService
}

func New(
	txManager db.TxManager,
	pendingChangesRepo contract.PendingChangesRepository,
	projectApproversRepo contract.ProjectApproversRepository,
	projectSettingsRepo contract.ProjectSettingsRepository,
	guardService contract.GuardService,
	featuresRepo contract.FeaturesRepository,
	featureParamsRepo contract.FeatureParamsRepository,
	rulesRepo contract.RulesRepository,
	flagVariantsRepo contract.FlagVariantsRepository,
	schedulesRepo contract.FeatureSchedulesRepository,
	featureTagsRepo contract.FeatureTagsRepository,
	auditLogRepo contract.AuditLogRepository,
	usersUseCase contract.UsersUseCase,
	permissionsService contract.PermissionsService,
) *Service {
	return &Service{
		txManager:            txManager,
		pendingChangesRepo:   pendingChangesRepo,
		projectApproversRepo: projectApproversRepo,
		projectSettingsRepo:  projectSettingsRepo,
		guardService:         guardService,
		featuresRepo:         featuresRepo,
		featureParamsRepo:    featureParamsRepo,
		rulesRepo:            rulesRepo,
		flagVariantsRepo:     flagVariantsRepo,
		schedulesRepo:        schedulesRepo,
		featureTagsRepo:      featureTagsRepo,
		auditLogRepo:         auditLogRepo,
		usersUseCase:         usersUseCase,
		permissionsService:   permissionsService,
	}
}

// Create creates a new pending change.
func (s *Service) Create(
	ctx context.Context,
	projectID domain.ProjectID,
	environmentID domain.EnvironmentID,
	requestedBy string,
	requestUserID *int,
	change domain.PendingChangePayload,
) (domain.PendingChange, error) {
	var created domain.PendingChange

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.pendingChangesRepo.Create(ctx, projectID, environmentID, requestedBy, requestUserID, change)
		if err != nil {
			return fmt.Errorf("create pending change: %w", err)
		}

		return nil
	}); err != nil {
		return domain.PendingChange{}, fmt.Errorf("tx create pending change: %w", err)
	}

	return created, nil
}

// GetByID retrieves a pending change by ID.
func (s *Service) GetByID(ctx context.Context, id domain.PendingChangeID) (domain.PendingChange, error) {
	return s.pendingChangesRepo.GetByID(ctx, id)
}

// List retrieves pending changes with filtering.
func (s *Service) List(
	ctx context.Context,
	filter contract.PendingChangesListFilter,
) ([]domain.PendingChange, int, error) {
	return s.pendingChangesRepo.List(ctx, filter)
}

// InitiateTOTPApproval creates a 2FA session for TOTP approval.
func (s *Service) InitiateTOTPApproval(
	ctx context.Context,
	id domain.PendingChangeID,
	approverUserID int,
) (string, error) {
	// Get the pending change
	pendingChange, err := s.pendingChangesRepo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("get pending change: %w", err)
	}

	// Check if already processed
	if pendingChange.Status != domain.PendingChangeStatusPending {
		return "", fmt.Errorf("pending change %s is not in pending status", id)
	}

	// Check if user is approver
	isApprover, err := s.IsUserApprover(ctx, pendingChange.ProjectID, approverUserID)
	if err != nil {
		return "", fmt.Errorf("check user approver: %w", err)
	}

	if !isApprover {
		return "", domain.ErrPermissionDenied
	}

	// Create 2FA session for TOTP approval
	sessionID, err := s.usersUseCase.InitiateTOTPApproval(ctx, domain.UserID(approverUserID))
	if err != nil {
		return "", fmt.Errorf("initiate TOTP approval: %w", err)
	}

	return sessionID, nil
}

// Approve approves a pending change and applies the changes.
func (s *Service) Approve(
	ctx context.Context,
	id domain.PendingChangeID,
	approverUserID int,
	approverName string,
	authMethod string,
	credential string,
	sessionID string, // Optional sessionID for TOTP approval
) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Get the pending change
		pendingChange, err := s.pendingChangesRepo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("get pending change: %w", err)
		}

		// Check if already processed
		if pendingChange.Status != domain.PendingChangeStatusPending {
			return fmt.Errorf("pending change %s is not in pending status", id)
		}

		// Check if user is approver
		isApprover, err := s.IsUserApprover(ctx, pendingChange.ProjectID, approverUserID)
		if err != nil {
			return fmt.Errorf("check user approver: %w", err)
		}

		if !isApprover {
			return domain.ErrPermissionDenied
		}

		// Verify credentials based on auth method
		switch authMethod {
		case "password":
			// Verify password using users usecase
			if err := s.usersUseCase.VerifyPassword(ctx, domain.UserID(approverUserID), credential); err != nil {
				return fmt.Errorf("password verification failed: %w", err)
			}
		case "totp":
			// Verify TOTP using 2FA session
			if sessionID == "" {
				return errors.New("sessionID is required for TOTP approval")
			}

			_, _, _, err := s.usersUseCase.Verify2FA(ctx, credential, sessionID)
			if err != nil {
				return fmt.Errorf("TOTP verification failed: %w", err)
			}
		default:
			return fmt.Errorf("unsupported auth method: %s", authMethod)
		}

		// Apply the changes
		if err := s.applyChanges(ctx, pendingChange); err != nil {
			return fmt.Errorf("apply changes: %w", err)
		}

		// Update status to approved
		now := time.Now().Format(time.RFC3339)
		if err := s.pendingChangesRepo.UpdateStatus(
			ctx,
			id,
			domain.PendingChangeStatusApproved,
			&approverName,
			&approverUserID,
			&now,
			nil,
			nil,
			nil,
		); err != nil {
			return fmt.Errorf("update pending change status: %w", err)
		}

		return nil
	})
}

// Reject rejects a pending change.
func (s *Service) Reject(
	ctx context.Context,
	id domain.PendingChangeID,
	rejectedBy string,
	reason string,
) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Get the pending change
		pendingChange, err := s.pendingChangesRepo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("get pending change: %w", err)
		}

		// Check if already processed
		if pendingChange.Status != domain.PendingChangeStatusPending {
			return fmt.Errorf("pending change %s is not in pending status", id)
		}

		// Update status to rejected
		now := time.Now().Format(time.RFC3339)
		if err := s.pendingChangesRepo.UpdateStatus(
			ctx,
			id,
			domain.PendingChangeStatusRejected,
			nil,
			nil,
			nil,
			&rejectedBy,
			&now,
			&reason,
		); err != nil {
			return fmt.Errorf("update pending change status: %w", err)
		}

		return nil
	})
}

// Cancel cancels a pending change.
func (s *Service) Cancel(
	ctx context.Context,
	id domain.PendingChangeID,
	cancelledBy string,
) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Get the pending change
		pendingChange, err := s.pendingChangesRepo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("get pending change: %w", err)
		}

		// Check if already processed
		if pendingChange.Status != domain.PendingChangeStatusPending {
			return fmt.Errorf("pending change %s is not in pending status", id)
		}

		// Update status to cancelled
		if err := s.pendingChangesRepo.UpdateStatus(
			ctx,
			id,
			domain.PendingChangeStatusCancelled,
			nil,
			nil,
			nil,
			&cancelledBy,
			nil,
			nil,
		); err != nil {
			return fmt.Errorf("update pending change status: %w", err)
		}

		return nil
	})
}

// CheckEntityConflict checks if there are any pending changes for the given entities.
func (s *Service) CheckEntityConflict(
	ctx context.Context,
	entities []domain.EntityChange,
) (bool, error) {
	return s.pendingChangesRepo.CheckEntityConflict(ctx, entities)
}

// GetProjectApprovers returns a list of users who can approve changes for a project.
func (s *Service) GetProjectApprovers(
	ctx context.Context,
	projectID domain.ProjectID,
) ([]domain.ProjectApprover, error) {
	return s.projectApproversRepo.GetByProjectID(ctx, projectID)
}

// IsUserApprover checks if a user can approve changes for a project.
func (s *Service) IsUserApprover(ctx context.Context, projectID domain.ProjectID, userID int) (bool, error) {
	// First, check explicit approvers
	isExplicitApprover, err := s.projectApproversRepo.IsUserApprover(ctx, projectID, userID)
	if err != nil {
		return false, fmt.Errorf("check explicit approver: %w", err)
	}

	if isExplicitApprover {
		return true, nil
	}

	// Check if user is superuser
	if appcontext.IsSuper(ctx) {
		return true, nil
	}

	// Check if the user can manage the project (project owner/manager)
	if err := s.permissionsService.CanManageProject(ctx, projectID); err == nil {
		return true, nil
	}

	return false, nil
}

// GetProjectActiveUserCount returns the number of active users in a project.
func (s *Service) GetProjectActiveUserCount(ctx context.Context, projectID domain.ProjectID) (int, error) {
	return s.guardService.GetProjectActiveUserCount(ctx, projectID)
}

// applyChanges applies the changes from a pending change.
func (s *Service) applyChanges(ctx context.Context, pendingChange domain.PendingChange) error {
	for _, entity := range pendingChange.Change.Entities {
		switch entity.Entity {
		case string(domain.EntityFeature):
			if err := s.applyFeatureChange(ctx, pendingChange.EnvironmentID, entity); err != nil {
				return fmt.Errorf("apply feature change: %w", err)
			}
		case string(domain.EntityFeatureParams):
			if err := s.applyFeatureParamsChange(ctx, pendingChange.EnvironmentID, entity); err != nil {
				return fmt.Errorf("apply feature_params change: %w", err)
			}
		case string(domain.EntityRule):
			if err := s.applyRuleChange(ctx, entity); err != nil {
				return fmt.Errorf("apply rule change: %w", err)
			}
		case string(domain.EntityFlagVariant):
			if err := s.applyFlagVariantChange(ctx, entity); err != nil {
				return fmt.Errorf("apply flag_variant change: %w", err)
			}
		case string(domain.EntityFeatureSchedule):
			if err := s.applyFeatureScheduleChange(ctx, entity); err != nil {
				return fmt.Errorf("apply feature_schedule change: %w", err)
			}
		case string(domain.EntityFeatureTag):
			if err := s.applyFeatureTagChange(ctx, entity); err != nil {
				return fmt.Errorf("apply feature_tag change: %w", err)
			}
		default:
			return fmt.Errorf("unsupported entity type: %s", entity.Entity)
		}
	}

	return nil
}

// applyFeatureParamsChange applies a change to feature_params (enabled, default_value) for the given env.
func (s *Service) applyFeatureParamsChange(
	ctx context.Context,
	envID domain.EnvironmentID,
	entity domain.EntityChange,
) error {
	featureID := domain.FeatureID(entity.EntityID)

	// Get current params
	current, err := s.featureParamsRepo.GetByFeatureWithEnv(ctx, featureID, envID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			// If not found, create baseline params with defaults before applying updates
			current = domain.FeatureParams{
				FeatureID:     featureID,
				EnvironmentID: envID,
			}
		} else {
			return fmt.Errorf("get current feature_params: %w", err)
		}
	}

	updated := current
	for field, change := range entity.Changes {
		switch field {
		case "enabled":
			if newVal, ok := change.New.(bool); ok {
				updated.Enabled = newVal
			}
		case "default_value":
			if newVal, ok := change.New.(string); ok {
				updated.DefaultValue = newVal
			}
		}
	}

	// We need project ID for audit when updating feature_params
	feature, ferr := s.featuresRepo.GetByID(ctx, featureID)
	if ferr != nil {
		return fmt.Errorf("get feature for project id: %w", ferr)
	}

	// Persist
	_, err = s.featureParamsRepo.Update(ctx, feature.ProjectID, updated)
	if err != nil {
		// If update failed because not found (when created baseline), try Create
		if errors.Is(err, domain.ErrEntityNotFound) {
			_, cerr := s.featureParamsRepo.Create(ctx, feature.ProjectID, updated)
			if cerr != nil {
				return fmt.Errorf("create feature_params: %w", cerr)
			}

			return nil
		}

		return fmt.Errorf("update feature_params: %w", err)
	}

	return nil
}

// applyRuleChange applies a change to a rule entity.
func (s *Service) applyRuleChange(
	ctx context.Context,
	entity domain.EntityChange,
) error {
	ruleID := domain.RuleID(entity.EntityID)

	switch entity.Action {
	case domain.EntityActionDelete:
		if err := s.rulesRepo.Delete(ctx, ruleID); err != nil {
			return fmt.Errorf("delete rule: %w", err)
		}

		return nil
	case domain.EntityActionUpdate:
		current, err := s.rulesRepo.GetByID(ctx, ruleID)
		if err != nil {
			return fmt.Errorf("get current rule: %w", err)
		}

		updated := current
		for field, change := range entity.Changes {
			switch field {
			case "is_customized":
				if v, ok := change.New.(bool); ok {
					updated.IsCustomized = v
				}
			case "action":
				if v, ok := change.New.(string); ok {
					updated.Action = domain.RuleAction(v)
				}
			case "priority":
				// JSON numbers come as float64
				if num, ok := change.New.(float64); ok {
					updated.Priority = uint8(num)
				}
			case "flag_variant_id":
				if change.New == nil {
					updated.FlagVariantID = nil
				} else if v, ok := change.New.(string); ok {
					id := domain.FlagVariantID(v)
					updated.FlagVariantID = &id
				}
			case "segment_id":
				if change.New == nil {
					updated.SegmentID = nil
				} else if v, ok := change.New.(string); ok {
					id := domain.SegmentID(v)
					updated.SegmentID = &id
				}
				// Note: conditions update is non-trivial (AST). If provided, keep existing for now.
			}
		}

		if _, err := s.rulesRepo.Update(ctx, updated); err != nil {
			return fmt.Errorf("update rule: %w", err)
		}

		return nil
	case domain.EntityActionInsert:
		// For insert we expect minimal required fields in changes: feature_id, project_id etc.
		var rule domain.Rule
		for field, change := range entity.Changes {
			switch field {
			case "project_id":
				if v, ok := change.New.(string); ok {
					rule.ProjectID = domain.ProjectID(v)
				}
			case "feature_id":
				if v, ok := change.New.(string); ok {
					rule.FeatureID = domain.FeatureID(v)
				}
			case "is_customized":
				if v, ok := change.New.(bool); ok {
					rule.IsCustomized = v
				}
			case "action":
				if v, ok := change.New.(string); ok {
					rule.Action = domain.RuleAction(v)
				}
			case "priority":
				if num, ok := change.New.(float64); ok {
					rule.Priority = uint8(num)
				}
			case "flag_variant_id":
				if change.New != nil {
					if v, ok := change.New.(string); ok {
						id := domain.FlagVariantID(v)
						rule.FlagVariantID = &id
					}
				}
			case "segment_id":
				if change.New != nil {
					if v, ok := change.New.(string); ok {
						id := domain.SegmentID(v)
						rule.SegmentID = &id
					}
				}
			}
		}
		_, err := s.rulesRepo.Create(ctx, rule)
		if err != nil {
			return fmt.Errorf("create rule: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unsupported action: %s", entity.Action)
	}
}

// applyFlagVariantChange applies a change to a flag variant entity.
func (s *Service) applyFlagVariantChange(
	ctx context.Context,
	entity domain.EntityChange,
) error {
	id := domain.FlagVariantID(entity.EntityID)

	switch entity.Action {
	case domain.EntityActionDelete:
		if err := s.flagVariantsRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("delete flag variant: %w", err)
		}

		return nil
	case domain.EntityActionUpdate:
		current, err := s.flagVariantsRepo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("get current flag variant: %w", err)
		}
		updated := current
		for field, change := range entity.Changes {
			switch field {
			case "name":
				if v, ok := change.New.(string); ok {
					updated.Name = v
				}
			case "rollout_percent":
				if num, ok := change.New.(float64); ok {
					updated.RolloutPercent = uint8(num)
				}
			}
		}
		if _, err := s.flagVariantsRepo.Update(ctx, updated); err != nil {
			return fmt.Errorf("update flag variant: %w", err)
		}

		return nil
	case domain.EntityActionInsert:
		var v domain.FlagVariant
		for field, change := range entity.Changes {
			switch field {
			case "project_id":
				if s, ok := change.New.(string); ok {
					v.ProjectID = domain.ProjectID(s)
				}
			case "feature_id":
				if s2, ok := change.New.(string); ok {
					v.FeatureID = domain.FeatureID(s2)
				}
			case "name":
				if s3, ok := change.New.(string); ok {
					v.Name = s3
				}
			case "rollout_percent":
				if num, ok := change.New.(float64); ok {
					v.RolloutPercent = uint8(num)
				}
			}
		}
		_, err := s.flagVariantsRepo.Create(ctx, v)
		if err != nil {
			return fmt.Errorf("create flag variant: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unsupported action: %s", entity.Action)
	}
}

// applyFeatureScheduleChange applies a change to a feature schedule entity.
func (s *Service) applyFeatureScheduleChange(
	ctx context.Context,
	entity domain.EntityChange,
) error {
	id := domain.FeatureScheduleID(entity.EntityID)

	parseTime := func(val interface{}) (*time.Time, bool) {
		if val == nil {
			return nil, true
		}
		if s, ok := val.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return &t, true
			}
		}

		return nil, false
	}

	parseDuration := func(val interface{}) (*time.Duration, bool) {
		if val == nil {
			return nil, true
		}
		switch v := val.(type) {
		case string:
			d, err := time.ParseDuration(v)
			if err == nil {
				return &d, true
			}
		case float64:
			// seconds
			d := time.Duration(v) * time.Second

			return &d, true
		}

		return nil, false
	}

	switch entity.Action {
	case domain.EntityActionDelete:
		if err := s.schedulesRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("delete feature schedule: %w", err)
		}

		return nil
	case domain.EntityActionUpdate:
		current, err := s.schedulesRepo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("get current feature schedule: %w", err)
		}
		updated := current
		for field, change := range entity.Changes {
			switch field {
			case "starts_at":
				if t, ok := parseTime(change.New); ok {
					updated.StartsAt = t
				}
			case "ends_at":
				if t, ok := parseTime(change.New); ok {
					updated.EndsAt = t
				}
			case "cron_expr":
				if change.New == nil {
					updated.CronExpr = nil
				} else if s, ok := change.New.(string); ok {
					updated.CronExpr = &s
				}
			case "cron_duration":
				if d, ok := parseDuration(change.New); ok {
					updated.CronDuration = d
				}
			case "timezone":
				if s, ok := change.New.(string); ok {
					updated.Timezone = s
				}
			case "action":
				if s, ok := change.New.(string); ok {
					updated.Action = domain.FeatureScheduleAction(s)
				}
			}
		}
		if _, err := s.schedulesRepo.Update(ctx, updated); err != nil {
			return fmt.Errorf("update feature schedule: %w", err)
		}

		return nil
	case domain.EntityActionInsert:
		var sch domain.FeatureSchedule
		for field, change := range entity.Changes {
			switch field {
			case "project_id":
				if s, ok := change.New.(string); ok {
					sch.ProjectID = domain.ProjectID(s)
				}
			case "feature_id":
				if s, ok := change.New.(string); ok {
					sch.FeatureID = domain.FeatureID(s)
				}
			case "environment_id":
				switch v := change.New.(type) {
				case float64:
					sch.EnvironmentID = domain.EnvironmentID(int64(v))
				case string:
					if n, err := strconv.ParseInt(v, 10, 64); err == nil {
						sch.EnvironmentID = domain.EnvironmentID(n)
					}
				}
			case "starts_at":
				if t, ok := parseTime(change.New); ok {
					sch.StartsAt = t
				}
			case "ends_at":
				if t, ok := parseTime(change.New); ok {
					sch.EndsAt = t
				}
			case "cron_expr":
				if change.New == nil {
					sch.CronExpr = nil
				} else if s, ok := change.New.(string); ok {
					sch.CronExpr = &s
				}
			case "cron_duration":
				if d, ok := parseDuration(change.New); ok {
					sch.CronDuration = d
				}
			case "timezone":
				if s, ok := change.New.(string); ok {
					sch.Timezone = s
				}
			case "action":
				if s, ok := change.New.(string); ok {
					sch.Action = domain.FeatureScheduleAction(s)
				}
			}
		}
		_, err := s.schedulesRepo.Create(ctx, sch)
		if err != nil {
			return fmt.Errorf("create feature schedule: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unsupported action: %s", entity.Action)
	}
}

// applyFeatureChange applies a change to a feature.
func (s *Service) applyFeatureChange(
	ctx context.Context,
	envID domain.EnvironmentID,
	entity domain.EntityChange,
) error {
	featureID := domain.FeatureID(entity.EntityID)

	switch entity.Action {
	case domain.EntityActionUpdate:
		// Get current feature
		currentFeature, err := s.featuresRepo.GetByID(ctx, featureID)
		if err != nil {
			return fmt.Errorf("get current feature: %w", err)
		}

		// Apply changes
		updatedFeature := currentFeature

		for field, change := range entity.Changes {
			switch field {
			case "name":
				if newValue, ok := change.New.(string); ok {
					updatedFeature.Name = newValue
				}
			case "description":
				if newValue, ok := change.New.(string); ok {
					updatedFeature.Description = newValue
				}
			}
		}

		// Update feature
		_, err = s.featuresRepo.Update(ctx, envID, updatedFeature)
		if err != nil {
			return fmt.Errorf("update feature: %w", err)
		}

	case domain.EntityActionDelete:
		// Delete feature
		err := s.featuresRepo.Delete(ctx, envID, featureID)
		if err != nil {
			return fmt.Errorf("delete feature: %w", err)
		}

	default:
		return fmt.Errorf("unsupported action: %s", entity.Action)
	}

	return nil
}

// applyFeatureTagChange applies a change to feature<->tag association.
func (s *Service) applyFeatureTagChange(
	ctx context.Context,
	entity domain.EntityChange,
) error {
	var featureID domain.FeatureID
	var tagID domain.TagID

	if ch, ok := entity.Changes["feature_id"]; ok {
		if v, ok2 := ch.New.(string); ok2 {
			featureID = domain.FeatureID(v)
		}
	}
	if ch, ok := entity.Changes["tag_id"]; ok {
		if v, ok2 := ch.New.(string); ok2 {
			tagID = domain.TagID(v)
		}
	}

	// Fallbacks if not provided in changes
	if featureID == "" {
		// try to infer from existing tag relation is not feasible here; require feature_id
		return errors.New("feature_id is required for feature_tag change")
	}
	if tagID == "" {
		// if not in changes, try entity.EntityID as tag id
		if entity.EntityID != "" {
			tagID = domain.TagID(entity.EntityID)
		} else {
			return errors.New("tag_id is required for feature_tag change")
		}
	}

	switch entity.Action {
	case domain.EntityActionInsert:
		if err := s.featureTagsRepo.AddFeatureTag(ctx, featureID, tagID); err != nil {
			return fmt.Errorf("add feature tag: %w", err)
		}

		return nil
	case domain.EntityActionDelete:
		if err := s.featureTagsRepo.RemoveFeatureTag(ctx, featureID, tagID); err != nil {
			return fmt.Errorf("remove feature tag: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unsupported action for feature_tag: %s", entity.Action)
	}
}
