package apibackend

import (
	"context"

	"github.com/google/uuid"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListPermissions(ctx context.Context) (generatedapi.ListPermissionsRes, error) {
	perms, err := r.membershipsUseCase.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	resp := make(generatedapi.ListPermissionsOKApplicationJSON, 0, len(perms))
	for _, perm := range perms {
		pid, err := uuid.Parse(string(perm.ID))
		if err != nil {
			return nil, err
		}
		resp = append(resp, generatedapi.Permission{ID: pid, Key: string(perm.Key), Name: perm.Name})
	}

	return &resp, nil
}
