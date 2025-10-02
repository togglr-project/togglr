package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type PermissionsService interface {
	CanAccessProject(ctx context.Context, projectID domain.ProjectID) error
	CanViewProject(ctx context.Context, projectID domain.ProjectID) error
	CanManageProject(ctx context.Context, projectID domain.ProjectID) error
	CanToggleFeature(ctx context.Context, projectID domain.ProjectID) error
	CanManageFeature(ctx context.Context, projectID domain.ProjectID) error
	CanManageSegment(ctx context.Context, projectID domain.ProjectID) error
	CanManageSchedule(ctx context.Context, projectID domain.ProjectID) error
	CanViewAudit(ctx context.Context, projectID domain.ProjectID) error
	CanManageMembership(ctx context.Context, projectID domain.ProjectID) error
	CanManageTags(ctx context.Context, projectID domain.ProjectID) error
	CanManageCategories(ctx context.Context) error
	GetAccessibleProjects(
		ctx context.Context,
		projects []domain.Project,
	) ([]domain.Project, error)
	HasProjectPermission(ctx context.Context, projectID domain.ProjectID, permKey domain.PermKey) (bool, error)
	HasGlobalPermission(ctx context.Context, permKey domain.PermKey) (bool, error)
	GetMyProjectPermissions(ctx context.Context) (map[domain.ProjectID][]domain.PermKey, error)
	GetMyProjectRoles(ctx context.Context) (map[domain.ProjectID]domain.Role, error)
}
