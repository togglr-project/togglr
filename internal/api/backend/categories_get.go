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

func (r *RestAPI) GetCategory(
	ctx context.Context,
	params generatedapi.GetCategoryParams,
) (generatedapi.GetCategoryRes, error) {
	userID := etogglcontext.UserID(ctx)
	categoryID := domain.CategoryID(params.CategoryID.String())

	// Get category
	category, err := r.categoriesUseCase.GetCategory(ctx, categoryID)
	if err != nil {
		slog.Error("get category failed", "error", err, "user_id", userID, "category_id", categoryID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("category not found"),
			}}, nil
		}

		return nil, err
	}

	// Convert to response
	item := generatedapi.Category{
		ID:        uuid.MustParse(category.ID.String()),
		Name:      category.Name,
		Slug:      category.Slug,
		Kind:      generatedapi.CategoryKind(category.Kind),
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if category.Description != nil {
		item.Description = generatedapi.NewOptNilString(*category.Description)
	}
	if category.Color != nil {
		item.Color = generatedapi.NewOptNilString(*category.Color)
	}

	resp := generatedapi.CategoryResponse{
		Category: item,
	}

	return &resp, nil
}
