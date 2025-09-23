package apibackend

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) CreateCategory(
	ctx context.Context,
	req *generatedapi.CreateCategoryRequest,
) (generatedapi.CreateCategoryRes, error) {
	userID := etogglcontext.UserID(ctx)

	// Check if user is superuser
	if !etogglcontext.IsSuper(ctx) {
		slog.Error("permission denied", "user_id", userID)

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

	// Create category
	category, err := r.categoriesUseCase.CreateCategory(
		ctx,
		req.Name,
		req.Slug,
		description,
		color,
		domain.CategoryType(req.CategoryType),
	)
	if err != nil {
		slog.Error("create category failed", "error", err, "user_id", userID)

		return nil, err
	}

	// Convert to response
	item := generatedapi.Category{
		ID:           uuid.MustParse(category.ID.String()),
		Name:         category.Name,
		Slug:         category.Slug,
		Kind:         generatedapi.CategoryKind(category.Kind),
		CategoryType: generatedapi.CategoryCategoryType(req.CategoryType),
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
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
