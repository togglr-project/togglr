package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) Setup2FA(ctx context.Context) (generatedapi.Setup2FARes, error) {
	userID := appcontext.UserID(ctx)
	secret, qrURL, qrImage, err := r.usersUseCase.Setup2FA(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		default:
			slog.Error("failed to setup 2FA", "error", err)

			return nil, r.NewError(ctx, err)
		}
	}
	resp := &generatedapi.TwoFASetupResponse{
		Secret:  secret,
		QrURL:   qrURL,
		QrImage: qrImage, // base64 PNG
	}

	return resp, nil
}
