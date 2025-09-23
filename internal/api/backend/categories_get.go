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

func (r *RestAPI) GetCategory(
	ctx context.Context,
	params generatedapi.GetCategoryParams,
) (generatedapi.GetCategoryRes, error) {
	userID := appcontext.UserID(ctx)
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
	item := dto.DomainCategoryToAPI(category)

	resp := generatedapi.CategoryResponse{
		Category: item,
	}

	return &resp, nil
}
