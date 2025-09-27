package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListProjectEnvironments(
	ctx context.Context,
	params generatedapi.ListProjectEnvironmentsParams,
) (generatedapi.ListProjectEnvironmentsRes, error) {
	userID := appcontext.UserID(ctx)

	projectID := domain.ProjectID(params.ProjectID.String())

	// Check if user can access the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", params.ProjectID)

		return &generatedapi.ErrorPermissionDenied{
			Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			},
		}, nil
	}

	// Get environments for the project
	environments, err := r.environmentsUseCase.ListByProjectID(ctx, projectID)
	if err != nil {
		slog.Error("list environments failed", "error", err, "user_id", userID, "project_id", params.ProjectID)

		return nil, err
	}

	items := dto.DomainEnvironmentsToAPI(environments)

	resp := &generatedapi.ListEnvironmentsResponse{
		Items: items,
	}

	return resp, nil
}
