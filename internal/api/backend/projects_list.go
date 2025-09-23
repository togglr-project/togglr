package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListProjects(ctx context.Context) (generatedapi.ListProjectsRes, error) {
	userID := appcontext.UserID(ctx)

	// Get all projects
	allProjects, err := r.projectsUseCase.List(ctx)
	if err != nil {
		slog.Error("get all projects failed", "error", err)

		return nil, err
	}

	// Filter projects based on user permissions
	projects, err := r.permissionsService.GetAccessibleProjects(ctx, allProjects)
	if err != nil {
		slog.Error("filter projects failed", "error", err, "user_id", userID)

		return nil, err
	}

	items := dto.DomainProjectsToAPI(projects)

	resp := generatedapi.ListProjectsResponse(items)

	return &resp, nil
}
