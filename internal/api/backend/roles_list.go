package apibackend

import (
	"context"

	"github.com/google/uuid"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListRoles(ctx context.Context) (generatedapi.ListRolesRes, error) {
	roles, err := r.membershipsUseCase.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	resp := make(generatedapi.ListRolesOKApplicationJSON, 0, len(roles))
	for _, role := range roles {
		rid, err := uuid.Parse(string(role.ID))
		if err != nil {
			return nil, err
		}
		resp = append(resp, generatedapi.Role{ID: rid, Key: role.Key, Name: role.Name, Description: role.Description})
	}

	return &resp, nil
}
