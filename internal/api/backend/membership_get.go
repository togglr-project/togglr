package apibackend

import (
	"context"
	"errors"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:nilerr // it's ok here
func (r *RestAPI) GetProjectMembership(
	ctx context.Context,
	params generatedapi.GetProjectMembershipParams,
) (generatedapi.GetProjectMembershipRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied")}}, nil
	}
	membership, err := r.membershipsUseCase.GetProjectMembership(
		ctx,
		projectID,
		domain.MembershipID(params.MembershipID.String()),
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("not found")}}, nil
		}

		return nil, err
	}
	membershipResponse, err := dto.DomainMembershipToAPI(membership)
	if err != nil {
		return nil, err
	}

	return &membershipResponse, nil
}
