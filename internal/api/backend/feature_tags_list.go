package apibackend

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) ListFeatureTags(
	ctx context.Context,
	params generatedapi.ListFeatureTagsParams,
) (generatedapi.ListFeatureTagsRes, error) {
	userID := etogglcontext.UserID(ctx)
	featureID := domain.FeatureID(params.FeatureID.String())

	// Get feature tags
	tags, err := r.featureTagsUseCase.ListFeatureTags(ctx, featureID)
	if err != nil {
		slog.Error("list feature tags failed", "error", err, "user_id", userID, "feature_id", featureID)

		if err == domain.ErrEntityNotFound {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	// Convert to response
	items := make([]generatedapi.ProjectTag, len(tags))
	for i, tag := range tags {
		item := generatedapi.ProjectTag{
			ID:        uuid.MustParse(tag.ID.String()),
			ProjectID: uuid.MustParse(tag.ProjectID.String()),
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}

		if tag.CategoryID != nil {
			item.CategoryID = generatedapi.NewOptNilUUID(uuid.MustParse(tag.CategoryID.String()))
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

		items[i] = item
	}

	resp := generatedapi.ListProjectTagsResponse(items)

	return &resp, nil
}
