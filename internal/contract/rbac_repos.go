package contract

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type RolesRepository interface {
	GetByKey(ctx context.Context, key string) (id string, err error)
}

type PermissionsRepository interface {
	RoleHasPermission(ctx context.Context, roleID string, key domain.PermKey) (bool, error)
}

type MembershipsRepository interface {
	GetForUserProject(ctx context.Context, userID int, projectID domain.ProjectID) (roleID string, err error)
}
