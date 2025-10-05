package apibackend

import (
	"context"
	"errors"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:nilerr // it's ok here
func (r *RestAPI) DeleteProjectMembership(
	ctx context.Context,
	params generatedapi.DeleteProjectMembershipParams,
) (generatedapi.DeleteProjectMembershipRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	if err := r.permissionsService.CanManageMembership(ctx, projectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied")}}, nil
	}
	if err := r.membershipsUseCase.DeleteProjectMembership(
		ctx,
		projectID,
		domain.MembershipID(params.MembershipID.String()),
	); err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("not found")}}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteProjectMembershipNoContent{}, nil
}
