package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListProjectTags(
	ctx context.Context,
	params generatedapi.ListProjectTagsParams,
) (generatedapi.ListProjectTagsRes, error) {
	userID := appcontext.UserID(ctx)
	projectID := domain.ProjectID(params.ProjectID.String())

	// Check if user can manage the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", projectID)
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	var categoryID *domain.CategoryID
	if params.CategoryID.Set {
		str := params.CategoryID.Value.String()
		categoryID = (*domain.CategoryID)(&str)
	}

	// Get project tags
	tags, err := r.tagsUseCase.ListProjectTags(ctx, projectID, categoryID)
	if err != nil {
		slog.Error("get project tags failed", "error", err, "user_id", userID, "project_id", projectID)
		return nil, err
	}

	items := dto.DomainTagsToAPI(tags)

	resp := generatedapi.ListProjectTagsResponse(items)

	return &resp, nil
}
