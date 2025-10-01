package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type RolesRepository interface {
	GetByKey(ctx context.Context, key string) (domain.Role, error)
	GetByID(ctx context.Context, id domain.RoleID) (domain.Role, error)
}

type PermissionsRepository interface {
	RoleHasPermission(ctx context.Context, roleID string, key domain.PermKey) (bool, error)
}

type MembershipsRepository interface {
	GetForUserProject(ctx context.Context, userID int, projectID domain.ProjectID) (roleID string, err error)
}
