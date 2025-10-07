package apibackend

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:nilerr // it's ok here
func (r *RestAPI) CreateProjectMembership(
	ctx context.Context,
	req *generatedapi.CreateMembershipRequest,
	params generatedapi.CreateProjectMembershipParams,
) (generatedapi.CreateProjectMembershipRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	if err := r.permissionsService.CanManageProject(ctx, projectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied")}}, nil
	}

	created, err := r.membershipsUseCase.CreateProjectMembership(
		ctx,
		projectID,
		domain.UserID(req.UserID),
		domain.RoleID(req.RoleID.String()),
	)
	if err != nil {
		return nil, err
	}
	resp, err := dto.DomainMembershipToAPI(created)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
