package apibackend

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) ListCategories(ctx context.Context) (generatedapi.ListCategoriesRes, error) {
	userID := etogglcontext.UserID(ctx)

	// Get all categories
	categories, err := r.categoriesUseCase.ListCategories(ctx)
	if err != nil {
		slog.Error("get categories failed", "error", err, "user_id", userID)
		return nil, err
	}

	items := make([]generatedapi.Category, 0, len(categories))
	for i := range categories {
		category := categories[i]
		item := generatedapi.Category{
			ID:           uuid.MustParse(category.ID.String()),
			Name:         category.Name,
			Slug:         category.Slug,
			Kind:         generatedapi.CategoryKind(category.Kind),
			CategoryType: generatedapi.CategoryCategoryType(category.Type),
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
		}

		if category.Description != nil {
			item.Description = generatedapi.NewOptNilString(*category.Description)
		}
		if category.Color != nil {
			item.Color = generatedapi.NewOptNilString(*category.Color)
		}

		items = append(items, item)
	}

	resp := generatedapi.ListCategoriesResponse(items)

	return &resp, nil
}
