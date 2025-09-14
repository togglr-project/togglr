package contract

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type PermissionsService interface {
	CanAccessProject(ctx context.Context, projectID domain.ProjectID) error
	CanManageProject(ctx context.Context, projectID domain.ProjectID) error
	GetAccessibleProjects(
		ctx context.Context,
		projects []domain.Project,
	) ([]domain.Project, error)
	HasProjectPermission(ctx context.Context, projectID domain.ProjectID, permKey domain.PermKey) (bool, error)
	HasGlobalPermission(ctx context.Context, permKey domain.PermKey) (bool, error)
	// GetMyProjectPermissions returns a map of projectID to permission keys for projects where the current user has a membership.
	GetMyProjectPermissions(ctx context.Context) (map[domain.ProjectID][]domain.PermKey, error)
}
