package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateEnvironment(
	ctx context.Context,
	req *generatedapi.CreateEnvironmentRequest,
	params generatedapi.CreateEnvironmentParams,
) (generatedapi.CreateEnvironmentRes, error) {
	userID := appcontext.UserID(ctx)

	projectID := domain.ProjectID(params.ProjectID.String())

	// Check if user can manage the project
	if err := r.permissionsService.CanManageProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", params.ProjectID)

		return &generatedapi.ErrorPermissionDenied{
			Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			},
		}, nil
	}

	environment, err := r.environmentsUseCase.Create(ctx, projectID, req.Key, req.Name)
	if err != nil {
		slog.Error("create environment failed", "error", err, "user_id", userID, "project_id", params.ProjectID)

		return nil, err
	}

	envResp := dto.DomainEnvironmentToAPI(environment)

	resp := &generatedapi.EnvironmentResponse{
		Environment: generatedapi.NewOptEnvironment(envResp),
	}

	return resp, nil
}
