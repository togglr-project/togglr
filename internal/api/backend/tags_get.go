package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) GetProjectTag(
	ctx context.Context,
	params generatedapi.GetProjectTagParams,
) (generatedapi.GetProjectTagRes, error) {
	userID := etogglcontext.UserID(ctx)
	projectID := domain.ProjectID(params.ProjectID.String())
	tagID := domain.TagID(params.TagID.String())

	// Check if user can manage the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID, "project_id", projectID)
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	// Get tag
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

	// Convert to response
	item := generatedapi.ProjectTag{
		ID:        uuid.MustParse(tag.ID.String()),
		ProjectID: uuid.MustParse(tag.ProjectID.String()),
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}

	if tag.Description != nil {
		item.Description = generatedapi.NewOptNilString(*tag.Description)
	}
	if tag.Color != nil {
		item.Color = generatedapi.NewOptNilString(*tag.Color)
	}

	// Convert category
	if tag.Category != nil {
		catItem := generatedapi.Category{
			ID:        uuid.MustParse(tag.Category.ID.String()),
			Name:      tag.Category.Name,
			Slug:      tag.Category.Slug,
			Kind:      generatedapi.CategoryKind(tag.Category.Kind),
			CreatedAt: tag.Category.CreatedAt,
			UpdatedAt: tag.Category.UpdatedAt,
		}

		if tag.Category.Description != nil {
			catItem.Description = generatedapi.NewOptNilString(*tag.Category.Description)
		}
		if tag.Category.Color != nil {
			catItem.Color = generatedapi.NewOptNilString(*tag.Category.Color)
		}

		item.Category = generatedapi.NewOptCategory(catItem)
	}

	resp := generatedapi.ProjectTagResponse{
		Tag: item,
	}

	return &resp, nil
}
