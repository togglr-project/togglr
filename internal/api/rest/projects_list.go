package rest

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggl/internal/context"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) ListProjects(ctx context.Context) (generatedapi.ListProjectsRes, error) {
	userID := etogglcontext.UserID(ctx)

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

	items := make([]generatedapi.Project, 0, len(projects))
	for i := range projects {
		project := projects[i]
		items = append(items, generatedapi.Project{
			ID:          project.ID.String(),
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   project.CreatedAt,
		})
	}

	resp := generatedapi.ListProjectsResponse(items)

	return &resp, nil
}
