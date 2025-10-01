package apibackend

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListProjectMemberships(
	ctx context.Context,
	params generatedapi.ListProjectMembershipsParams,
) (generatedapi.ListProjectMembershipsRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{Message: generatedapi.NewOptString("permission denied")}}, nil
	}
	items, err := r.membershipsUseCase.ListProjectMemberships(ctx, projectID)
	if err != nil {
		return nil, err
	}
	resp := make(generatedapi.ListProjectMembershipsOKApplicationJSON, 0, len(items))
	for _, item := range items {
		membership, err := dto.DomainMembershipToAPI(item)
		if err != nil {
			return nil, err
		}
		resp = append(resp, membership)
	}

	return &resp, nil
}
