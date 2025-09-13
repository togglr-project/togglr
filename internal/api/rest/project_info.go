package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) GetProject(
	ctx context.Context,
	params generatedapi.GetProjectParams,
) (generatedapi.GetProjectRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check if the user has access to the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", projectID)

		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			}}, nil
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		}

		return nil, err
	}

	project, err := r.projectsUseCase.GetProject(ctx, projectID)
	if err != nil {
		slog.Error("get project failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.ProjectResponse{
		Project: generatedapi.Project{
			ID:          project.ID.String(),
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   project.CreatedAt,
		},
	}, nil
}
