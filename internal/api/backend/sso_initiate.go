package apibackend

import (
	"context"
	"log/slog"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) SSOInitiate(
	ctx context.Context,
	params generatedapi.SSOInitiateParams,
) (generatedapi.SSOInitiateRes, error) {
	redirectURL, err := r.usersUseCase.SSOInitiate(ctx, params.Provider)
	if err != nil {
		slog.Error("SSO initiate failed", "error", err)

		return nil, err
	}

	return &generatedapi.SSOInitiateResponse{
		RedirectURL: redirectURL,
	}, nil
}
