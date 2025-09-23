package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetProductInfo(ctx context.Context) (generatedapi.GetProductInfoRes, error) {
	// Check if the user is a superuser
	if !appcontext.IsSuper(ctx) {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("unauthorized"),
		}}, nil
	}

	productInfo, err := r.productInfoUseCase.GetProductInfo(ctx)
	if err != nil {
		slog.Error("get product info failed", "error", err)

		return nil, err
	}

	return &generatedapi.ProductInfoResponse{
		ClientID:  productInfo.ClientID,
		CreatedAt: productInfo.CreatedAt,
	}, nil
}
