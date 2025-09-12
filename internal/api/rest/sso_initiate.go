package rest

import (
	"context"
	"log/slog"

	"github.com/rom8726/etoggl/internal/domain"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) SSOInitiate(
	ctx context.Context,
	params generatedapi.SSOInitiateParams,
) (generatedapi.SSOInitiateRes, error) {
	// Check if SSO feature is available in the current license
	isAvailable, err := r.licenseUseCase.IsFeatureAvailable(ctx, domain.FeatureSSO)
	if err != nil {
		slog.Error("Failed to check license for SSO feature", "error", err)

		return nil, err
	}

	if !isAvailable {
		slog.Error("SSO feature not available in current license")

		return &generatedapi.ErrorBadRequest{
			Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString("SSO feature not available in current license"),
			},
		}, nil
	}

	redirectURL, err := r.usersUseCase.SSOInitiate(ctx, params.Provider)
	if err != nil {
		slog.Error("SSO initiate failed", "error", err)

		return nil, err
	}

	return &generatedapi.SSOInitiateResponse{
		RedirectURL: redirectURL,
	}, nil
}
