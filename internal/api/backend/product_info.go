package apibackend

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) GetProductInfo(ctx context.Context) (generatedapi.GetProductInfoRes, error) {
	// Check if the user is a superuser
	if !etogglcontext.IsSuper(ctx) {
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
