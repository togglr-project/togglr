package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteEnvironment(
	ctx context.Context,
	params generatedapi.DeleteEnvironmentParams,
) (generatedapi.DeleteEnvironmentRes, error) {
	userID := appcontext.UserID(ctx)

	// Get environment to check project access
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

	// Check if user can manage the project
	if err := r.permissionsService.CanManageProject(ctx, environment.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", environment.ProjectID)

		return &generatedapi.ErrorPermissionDenied{
			Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			},
		}, nil
	}

	// Delete environment
	err = r.environmentsUseCase.Delete(ctx, envID)
	if err != nil {
		slog.Error("delete environment failed", "error", err, "user_id", userID, "environment_id", params.EnvironmentID)

		return nil, err
	}

	return &generatedapi.DeleteEnvironmentNoContent{}, nil
}
