package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListUsers(ctx context.Context) (generatedapi.ListUsersRes, error) {
	users, err := r.usersUseCase.List(ctx)
	if err != nil {
		slog.Error("list users failed", "error", err)

		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("Only superusers can list users"),
				},
			}, nil
		}

		return nil, err
	}

	resp := dto.DomainUsersToAPI(users)
	listResp := generatedapi.ListUsersResponse(resp)

	return &listResp, nil
}
