package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) UpdateEnvironment(
	ctx context.Context,
	req *generatedapi.UpdateEnvironmentRequest,
	params generatedapi.UpdateEnvironmentParams,
) (generatedapi.UpdateEnvironmentRes, error) {
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

	// Update environment
	updatedEnvironment, err := r.environmentsUseCase.Update(ctx, envID, req.Name)
	if err != nil {
		slog.Error("update environment failed", "error", err, "user_id", userID, "environment_id", params.EnvironmentID)

		return nil, err
	}

	envResp := dto.DomainEnvironmentToAPI(updatedEnvironment)

	resp := &generatedapi.EnvironmentResponse{
		Environment: generatedapi.NewOptEnvironment(envResp),
	}

	return resp, nil
}
