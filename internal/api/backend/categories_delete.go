package apibackend

import (
	"context"
	"errors"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) DeleteCategory(
	ctx context.Context,
	params generatedapi.DeleteCategoryParams,
) (generatedapi.DeleteCategoryRes, error) {
	userID := etogglcontext.UserID(ctx)

	// Check if user is superuser
	if !etogglcontext.IsSuper(ctx) {
		slog.Error("permission denied", "user_id", userID)
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
