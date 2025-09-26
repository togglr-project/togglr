package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

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
