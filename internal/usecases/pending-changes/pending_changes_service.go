package pending_changes

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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

	// Apply changes using reflection
	if err := ApplyChangesToEntity(&current, entity.Changes); err != nil {
		return fmt.Errorf("apply changes to feature_params: %w", err)
	}

	// We need project ID for audit when updating feature_params
	feature, ferr := s.featuresRepo.GetByID(ctx, featureID)
	if ferr != nil {
		return fmt.Errorf("get feature for project id: %w", ferr)
	}

	// Persist
	_, err = s.featureParamsRepo.Update(ctx, feature.ProjectID, current)
	if err != nil {
		// If update failed because not found (when created baseline), try Create
		if errors.Is(err, domain.ErrEntityNotFound) {
			_, cerr := s.featureParamsRepo.Create(ctx, feature.ProjectID, current)
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

		// Apply changes using reflection
		if err := ApplyChangesToEntity(&current, entity.Changes); err != nil {
			return fmt.Errorf("apply changes to rule: %w", err)
		}

		if _, err := s.rulesRepo.Update(ctx, current); err != nil {
			return fmt.Errorf("update rule: %w", err)
		}

		return nil
	case domain.EntityActionInsert:
		// Create new rule from changes using reflection
		rule, err := CreateEntityFromChanges(reflect.TypeOf(domain.Rule{}), entity.Changes)
		if err != nil {
			return fmt.Errorf("create rule from changes: %w", err)
		}

		_, err = s.rulesRepo.Create(ctx, *rule.(*domain.Rule))
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

		// Apply changes using reflection
		if err := ApplyChangesToEntity(&current, entity.Changes); err != nil {
			return fmt.Errorf("apply changes to flag variant: %w", err)
		}

		if _, err := s.flagVariantsRepo.Update(ctx, current); err != nil {
			return fmt.Errorf("update flag variant: %w", err)
		}

		return nil
	case domain.EntityActionInsert:
		// Create new flag variant from changes using reflection
		variant, err := CreateEntityFromChanges(reflect.TypeOf(domain.FlagVariant{}), entity.Changes)
		if err != nil {
			return fmt.Errorf("create flag variant from changes: %w", err)
		}

		_, err = s.flagVariantsRepo.Create(ctx, *variant.(*domain.FlagVariant))
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

		// Apply changes using reflection
		if err := ApplyChangesToEntity(&current, entity.Changes); err != nil {
			return fmt.Errorf("apply changes to feature schedule: %w", err)
		}

		if _, err := s.schedulesRepo.Update(ctx, current); err != nil {
			return fmt.Errorf("update feature schedule: %w", err)
		}

		return nil
	case domain.EntityActionInsert:
		// Create new feature schedule from changes using reflection
		schedule, err := CreateEntityFromChanges(reflect.TypeOf(domain.FeatureSchedule{}), entity.Changes)
		if err != nil {
			return fmt.Errorf("create feature schedule from changes: %w", err)
		}

		_, err = s.schedulesRepo.Create(ctx, *schedule.(*domain.FeatureSchedule))
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

		// Apply changes using reflection
		if err := ApplyChangesToEntity(&currentFeature, entity.Changes); err != nil {
			return fmt.Errorf("apply changes to feature: %w", err)
		}

		// Update feature
		_, err = s.featuresRepo.Update(ctx, envID, currentFeature)
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
