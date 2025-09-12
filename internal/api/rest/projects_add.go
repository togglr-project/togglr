package rest

import (
	"context"
	"log/slog"

	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) AddProject(
	ctx context.Context,
	req *generatedapi.AddProjectRequest,
) (generatedapi.AddProjectRes, error) {
	_, err := r.projectsUseCase.CreateProject(ctx, req.Name, req.Description)
	if err != nil {
		slog.Error("add project failed", "error", err)

		return nil, err
	}

	return &generatedapi.AddProjectCreated{}, nil
}
