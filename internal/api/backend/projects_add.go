package apibackend

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) AddProject(
	ctx context.Context,
	req *generatedapi.AddProjectRequest,
) (generatedapi.AddProjectRes, error) {
	// Require authenticated user with project.create global permission (or superuser)
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("unauthorized"),
		}}, nil
	}

	allowed, err := r.permissionsService.HasGlobalPermission(ctx, domain.PermProjectCreate)
	if err != nil || !allowed {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("unauthorized"),
		}}, nil
	}

	_, err = r.projectsUseCase.CreateProject(ctx, req.Name, req.Description)
	if err != nil {
		slog.Error("add project failed", "error", err)
		return nil, err
	}

	return &generatedapi.AddProjectCreated{}, nil
}
