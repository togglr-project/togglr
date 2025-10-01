package apibackend

import (
	"context"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetRolePermissions(
	ctx context.Context,
	params generatedapi.GetRolePermissionsParams,
) (generatedapi.GetRolePermissionsRes, error) {
	perms, err := r.membershipsUseCase.GetRolePermissions(ctx, domain.RoleID(params.RoleID.String()))
	if err != nil {
		return nil, err
	}
	resp := make(generatedapi.GetRolePermissionsOKApplicationJSON, 0, len(perms))
	for _, perm := range perms {
		pid, err := uuid.Parse(string(perm.ID))
		if err != nil {
			return nil, err
		}
		resp = append(resp, generatedapi.Permission{ID: pid, Key: string(perm.Key), Name: perm.Name})
	}

	return &resp, nil
}
