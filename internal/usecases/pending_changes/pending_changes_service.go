package pending_changes

import (
	"context"
	"errors"
	"fmt"
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
	auditLogRepo         contract.AuditLogRepository
	usersUseCase         contract.UsersUseCase
	permissionsService   contract.PermissionsService
}

func New(
	txManager db.TxManager,
	pendingChangesRepo contract.PendingChangesRepository,
	projectApproversRepo contract.ProjectApproversRepository,
	projectSettingsRepo contract.ProjectSettingsRepository,
	guardService contract.GuardService,
	featuresRepo contract.FeaturesRepository,
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
		auditLogRepo:         auditLogRepo,
		usersUseCase:         usersUseCase,
		permissionsService:   permissionsService,
	}
}

// Create creates a new pending change.
func (s *Service) Create(
	ctx context.Context,
	projectID domain.ProjectID,
	requestedBy string,
	requestUserID *int,
	change domain.PendingChangePayload,
) (domain.PendingChange, error) {
	var created domain.PendingChange

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.pendingChangesRepo.Create(ctx, projectID, requestedBy, requestUserID, change)
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
		case "feature":
			if err := s.applyFeatureChange(ctx, entity); err != nil {
				return fmt.Errorf("apply feature change: %w", err)
			}
		// TODO: Add other entity types
		default:
			return fmt.Errorf("unsupported entity type: %s", entity.Entity)
		}
	}

	return nil
}

// applyFeatureChange applies a change to a feature.
func (s *Service) applyFeatureChange(
	ctx context.Context,
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
			case "enabled":
				if newValue, ok := change.New.(bool); ok {
					updatedFeature.Enabled = newValue
				}
			case "default_variant":
				if newValue, ok := change.New.(string); ok {
					updatedFeature.DefaultVariant = newValue
				}
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
		_, err = s.featuresRepo.Update(ctx, updatedFeature)
		if err != nil {
			return fmt.Errorf("update feature: %w", err)
		}

	case domain.EntityActionDelete:
		// Delete feature
		err := s.featuresRepo.Delete(ctx, featureID)
		if err != nil {
			return fmt.Errorf("delete feature: %w", err)
		}

	default:
		return fmt.Errorf("unsupported action: %s", entity.Action)
	}

	return nil
}
