package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:nilerr // it's ok here
func (r *RestAPI) AddProject(
	ctx context.Context,
	req *generatedapi.AddProjectRequest,
) (generatedapi.AddProjectRes, error) {
	// Require authenticated user with project.create global permission (or superuser)
	userID := appcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("unauthorized"),
		}}, nil
	}

	allowed, err := r.permissionsService.HasGlobalPermission(ctx, domain.PermProjectCreate)
	if err != nil || !allowed {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	_, err = r.projectsUseCase.CreateProject(ctx, req.Name, req.Description)
	if err != nil {
		slog.Error("add project failed", "error", err)

		if errors.Is(err, domain.ErrEntityAlreadyExists) {
			return &generatedapi.Error{
				Error: generatedapi.ErrorError{Message: generatedapi.NewOptString("project already exists")},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.AddProjectCreated{}, nil
}
