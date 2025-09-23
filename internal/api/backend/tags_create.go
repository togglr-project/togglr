package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateProjectTag(
	ctx context.Context,
	req *generatedapi.CreateProjectTagRequest,
	params generatedapi.CreateProjectTagParams,
) (generatedapi.CreateProjectTagRes, error) {
	userID := appcontext.UserID(ctx)
	projectID := domain.ProjectID(params.ProjectID.String())

	// Check if user can manage the project
	if err := r.permissionsService.CanManageProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", projectID)
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	// Convert request to domain
	var description *string
	if req.Description.Set {
		description = &req.Description.Value
	}

	var color *string
	if req.Color.Set {
		color = &req.Color.Value
	}

	var categoryID *domain.CategoryID
	if req.CategoryID.Set {
		str := req.CategoryID.Value.String()
		categoryID = (*domain.CategoryID)(&str)
	}

	// Create tag
	tag, err := r.tagsUseCase.CreateTag(
		ctx,
		projectID,
		categoryID,
		req.Name,
		req.Slug,
		description,
		color,
	)
	if err != nil {
		slog.Error("create tag failed", "error", err, "user_id", userID, "project_id", projectID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	// Convert to response
	item := dto.DomainTagToAPI(tag)

	resp := generatedapi.ProjectTagResponse{
		Tag: item,
	}

	return &resp, nil
}
