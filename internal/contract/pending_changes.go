package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type PendingChangesListFilter struct {
	ProjectID *domain.ProjectID
	Status    *domain.PendingChangeStatus
	UserID    *int
	SortBy    string // created_at, status, requested_by
	SortDesc  bool
	Page      uint
	PerPage   uint
}

type PendingChangesUseCase interface {
	// Create creates a new pending change
	Create(
		ctx context.Context,
		projectID domain.ProjectID,
		requestedBy string,
		requestUserID *int,
		change domain.PendingChangePayload,
	) (domain.PendingChange, error)

	// GetByID retrieves a pending change by ID
	GetByID(ctx context.Context, id domain.PendingChangeID) (domain.PendingChange, error)

	// List retrieves pending changes with filtering
	List(
		ctx context.Context,
		filter PendingChangesListFilter,
	) ([]domain.PendingChange, int, error)

	// Approve approves a pending change and applies the changes
	Approve(
		ctx context.Context,
		id domain.PendingChangeID,
		approverUserID int,
		approverName string,
		authMethod string, // "password" or "totp"
		credential string, // password or TOTP code
	) error

	// Reject rejects a pending change
	Reject(
		ctx context.Context,
		id domain.PendingChangeID,
		rejectedBy string,
		reason string,
	) error

	// Cancel cancels a pending change
	Cancel(
		ctx context.Context,
		id domain.PendingChangeID,
		cancelledBy string,
	) error

	// CheckEntityConflict checks if there are any pending changes for the given entities
	CheckEntityConflict(
		ctx context.Context,
		entities []domain.EntityChange,
	) (bool, error)

	// GetProjectApprovers returns list of users who can approve changes for a project
	GetProjectApprovers(ctx context.Context, projectID domain.ProjectID) ([]domain.ProjectApprover, error)

	// IsUserApprover checks if a user can approve changes for a project
	IsUserApprover(ctx context.Context, projectID domain.ProjectID, userID int) (bool, error)

	// GetProjectActiveUserCount returns the number of active users in a project
	GetProjectActiveUserCount(ctx context.Context, projectID domain.ProjectID) (int, error)
}

type PendingChangesRepository interface {
	// Create creates a new pending change and its entities
	Create(
		ctx context.Context,
		projectID domain.ProjectID,
		requestedBy string,
		requestUserID *int,
		change domain.PendingChangePayload,
	) (domain.PendingChange, error)

	// GetByID retrieves a pending change by ID
	GetByID(ctx context.Context, id domain.PendingChangeID) (domain.PendingChange, error)

	// List retrieves pending changes with filtering
	List(
		ctx context.Context,
		filter PendingChangesListFilter,
	) ([]domain.PendingChange, int, error)

	// UpdateStatus updates the status of a pending change
	UpdateStatus(
		ctx context.Context,
		id domain.PendingChangeID,
		status domain.PendingChangeStatus,
		approvedBy *string,
		approvedUserID *int,
		approvedAt *string,
		rejectedBy *string,
		rejectedAt *string,
		rejectionReason *string,
	) error

	// CheckEntityConflict checks if there are any pending changes for the given entities
	CheckEntityConflict(
		ctx context.Context,
		entities []domain.EntityChange,
	) (bool, error)

	// GetEntitiesByPendingChangeID retrieves all entities for a pending change
	GetEntitiesByPendingChangeID(
		ctx context.Context,
		pendingChangeID domain.PendingChangeID,
	) ([]domain.PendingChangeEntity, error)
}

type ProjectApproversRepository interface {
	// Create adds a new approver to a project
	Create(ctx context.Context, approver domain.ProjectApprover) error

	// GetByProjectID retrieves all approvers for a project
	GetByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.ProjectApprover, error)

	// Delete removes an approver from a project
	Delete(ctx context.Context, projectID domain.ProjectID, userID int) error

	// IsUserApprover checks if a user is an approver for a project
	IsUserApprover(ctx context.Context, projectID domain.ProjectID, userID int) (bool, error)
}

type ProjectSettingsRepository interface {
	// Set sets a project setting
	Set(
		ctx context.Context,
		projectID domain.ProjectID,
		name string,
		value interface{},
	) error

	// Get retrieves a project setting
	Get(
		ctx context.Context,
		projectID domain.ProjectID,
		name string,
	) (domain.ProjectSetting, error)

	// GetAll retrieves all settings for a project
	GetAll(ctx context.Context, projectID domain.ProjectID) ([]domain.ProjectSetting, error)

	// Delete removes a project setting
	Delete(ctx context.Context, projectID domain.ProjectID, name string) error
}

// GuardService provides methods to check if entities are guarded
type GuardService interface {
	// IsFeatureGuarded checks if a feature has the guarded tag
	IsFeatureGuarded(ctx context.Context, featureID domain.FeatureID) (bool, error)

	// IsEntityGuarded checks if any entity in the list is guarded
	IsEntityGuarded(ctx context.Context, entities []domain.EntityChange) (bool, error)

	// GetProjectActiveUserCount returns the number of active users in a project
	GetProjectActiveUserCount(ctx context.Context, projectID domain.ProjectID) (int, error)
}
