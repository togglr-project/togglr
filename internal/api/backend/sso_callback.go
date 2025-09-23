package apibackend

import (
	"context"
	"errors"
	"log/slog"
	"time"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) SSOCallback(
	ctx context.Context,
	request *generatedapi.SSOCallbackRequest,
) (generatedapi.SSOCallbackRes, error) {
	accessToken, refreshToken, _, err := r.usersUseCase.SSOCallback(
		ctx, request.Provider, appcontext.RawRequest(ctx), request.Response, request.State)
	if err != nil {
		slog.Error("SSO callback failed", "error", err)

		if errors.Is(err, domain.ErrInactiveUser) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("user is inactive"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.LoginResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		ExpiresIn:     int(time.Now().Add(r.tokenizer.AccessTokenTTL()).Unix()),
		IsTmpPassword: false,
	}, nil
}
