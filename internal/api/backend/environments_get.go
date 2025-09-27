package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetEnvironment(
	ctx context.Context,
	params generatedapi.GetEnvironmentParams,
) (generatedapi.GetEnvironmentRes, error) {
	userID := appcontext.UserID(ctx)

	// Get environment
	envID := domain.EnvironmentID(params.EnvironmentID)
	environment, err := r.environmentsUseCase.GetByID(ctx, envID)
	if err != nil {
		slog.Error("get environment failed", "error", err, "user_id", userID, "environment_id", params.EnvironmentID)

		return &generatedapi.ErrorNotFound{
			Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("environment not found"),
			},
		}, nil
	}

	// Check if user can access the project
	if err := r.permissionsService.CanAccessProject(ctx, environment.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", environment.ProjectID)

		return &generatedapi.ErrorPermissionDenied{
			Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			},
		}, nil
	}

	envResp := dto.DomainEnvironmentToAPI(environment)

	resp := &generatedapi.EnvironmentResponse{
		Environment: generatedapi.NewOptEnvironment(envResp),
	}

	return resp, nil
}
