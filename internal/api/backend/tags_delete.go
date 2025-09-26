package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteProjectTag(
	ctx context.Context,
	params generatedapi.DeleteProjectTagParams,
) (generatedapi.DeleteProjectTagRes, error) {
	userID := appcontext.UserID(ctx)
	projectID := domain.ProjectID(params.ProjectID.String())
	tagID := domain.TagID(params.TagID.String())

	// Check if user can manage the project
	if err := r.permissionsService.CanManageProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", projectID)

		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	// Get tag to verify it belongs to the project
	tag, err := r.tagsUseCase.GetTag(ctx, tagID)
	if err != nil {
		slog.Error("get tag failed", "error", err, "user_id", userID, "tag_id", tagID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("tag not found"),
			}}, nil
		}

		return nil, err
	}

	// Verify tag belongs to the project
	if tag.ProjectID != projectID {
		return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
			Message: generatedapi.NewOptString("tag not found"),
		}}, nil
	}

	// Delete tag
	err = r.tagsUseCase.DeleteTag(ctx, tagID)
	if err != nil {
		slog.Error("delete tag failed", "error", err, "user_id", userID, "tag_id", tagID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("tag not found"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteProjectTagNoContent{}, nil
}
