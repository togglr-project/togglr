package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateCategory(
	ctx context.Context,
	req *generatedapi.CreateCategoryRequest,
) (generatedapi.CreateCategoryRes, error) {
	userID := appcontext.UserID(ctx)

	// Check if user is superuser
	if !appcontext.IsSuper(ctx) {
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
	item := dto.DomainCategoryToAPI(category)

	resp := generatedapi.CategoryResponse{
		Category: item,
	}

	return &resp, nil
}
