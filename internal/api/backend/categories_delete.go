package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteCategory(
	ctx context.Context,
	params generatedapi.DeleteCategoryParams,
) (generatedapi.DeleteCategoryRes, error) {
	userID := appcontext.UserID(ctx)

	// Check if user can manage categories
	if err := r.permissionsService.CanManageCategories(ctx); err != nil {
		slog.Error("permission denied", "error", err, "user_id", userID)

		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	categoryID := domain.CategoryID(params.CategoryID.String())

	// Delete category
	err := r.categoriesUseCase.DeleteCategory(ctx, categoryID)
	if err != nil {
		slog.Error("delete category failed", "error", err, "user_id", userID, "category_id", categoryID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("category not found"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteCategoryNoContent{}, nil
}
