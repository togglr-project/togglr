package rest

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	etogglcontext "github.com/rom8726/etoggl/internal/context"
	"github.com/rom8726/etoggl/internal/domain"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) Reset2FA(
	ctx context.Context,
	req *generatedapi.TwoFAResetRequest,
) (generatedapi.Reset2FARes, error) {
	userID := etogglcontext.UserID(ctx)
	secret, qrURL, qrImage, err := r.usersUseCase.Reset2FA(ctx, userID, req.EmailCode)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalid2FACode):
			return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString("invalid code"),
			}}, nil
		case errors.Is(err, domain.ErrUserNotFound):
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		default:
			slog.Error("failed to reset 2FA", "error", err)

			return nil, r.NewError(ctx, err)
		}
	}
	resp := &generatedapi.TwoFASetupResponse{
		Secret:  secret,
		QrURL:   qrURL,
		QrImage: qrImage,
	}

	return resp, nil
}
