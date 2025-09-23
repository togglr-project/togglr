package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListCategories(ctx context.Context) (generatedapi.ListCategoriesRes, error) {
	userID := appcontext.UserID(ctx)

	// Get all categories
	categories, err := r.categoriesUseCase.ListCategories(ctx)
	if err != nil {
		slog.Error("get categories failed", "error", err, "user_id", userID)
		return nil, err
	}

	items := dto.DomainCategoriesToAPI(categories)

	resp := generatedapi.ListCategoriesResponse(items)

	return &resp, nil
}
