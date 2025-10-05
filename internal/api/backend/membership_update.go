package apibackend

import (
	"context"
	"errors"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:nilerr // it's ok here
func (r *RestAPI) UpdateProjectMembership(
	ctx context.Context,
	req *generatedapi.UpdateMembershipRequest,
	params generatedapi.UpdateProjectMembershipParams,
) (generatedapi.UpdateProjectMembershipRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	if err := r.permissionsService.CanManageMembership(ctx, projectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied")}}, nil
	}
	projectMembership, err := r.membershipsUseCase.UpdateProjectMembership(
		ctx,
		projectID,
		domain.MembershipID(params.MembershipID.String()),
		domain.RoleID(req.RoleID.String()),
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("not found")}}, nil
		}

		return nil, err
	}
	resp, err := dto.DomainMembershipToAPI(projectMembership)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
