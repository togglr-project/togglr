package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/internal/dto"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) ListProjectChanges(
	ctx context.Context,
	params generatedapi.ListProjectChangesParams,
) (generatedapi.ListProjectChangesRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())

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

	// Convert API parameters to domain filter
	filter := dto.APIChangesFilterToDomain(projectID, params)

	// Get changes from a use case
	result, err := r.projectsUseCase.ListChanges(ctx, filter)
	if err != nil {
		slog.Error("list project changes failed", "error", err, "project_id", projectID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("project not found"),
			}}, nil
		}

		return nil, err
	}

	// Convert domain result to API response
	response := dto.DomainChangesToAPI(result)

	return &response, nil
}
