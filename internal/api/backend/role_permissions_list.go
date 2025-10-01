package apibackend

import (
	"context"

	"github.com/google/uuid"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListRolePermissions(ctx context.Context) (generatedapi.ListRolePermissionsRes, error) {
	permissions, err := r.membershipsUseCase.ListRolePermissions(ctx)
	if err != nil {
		return nil, err
	}
	resp := make(generatedapi.ListRolePermissionsOKApplicationJSON, 0, len(permissions))
	for role, perms := range permissions {
		rid, err := uuid.Parse(string(role.ID))
		if err != nil {
			return nil, err
		}
		apiRole := generatedapi.Role{ID: rid, Key: role.Key, Name: role.Name, Description: role.Description}
		apiPerms := make([]generatedapi.Permission, 0, len(perms))
		for _, perm := range perms {
			pid, err := uuid.Parse(string(perm.ID))
			if err != nil {
				return nil, err
			}
			apiPerms = append(apiPerms, generatedapi.Permission{ID: pid, Key: string(perm.Key), Name: perm.Name})
		}
		item := generatedapi.ListRolePermissionsOKItem{Role: generatedapi.NewOptRole(apiRole), Permissions: apiPerms}
		resp = append(resp, item)
	}

	return &resp, nil
}
