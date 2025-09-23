package apibackend

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) ListProjectTags(
	ctx context.Context,
	params generatedapi.ListProjectTagsParams,
) (generatedapi.ListProjectTagsRes, error) {
	userID := etogglcontext.UserID(ctx)
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

	items := make([]generatedapi.ProjectTag, 0, len(tags))
	for i := range tags {
		tag := tags[i]
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

		items = append(items, item)
	}

	resp := generatedapi.ListProjectTagsResponse(items)

	return &resp, nil
}
