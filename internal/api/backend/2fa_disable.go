package apibackend

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) Disable2FA(
	ctx context.Context,
	req *generatedapi.TwoFADisableRequest,
) (generatedapi.Disable2FARes, error) {
	userID := appcontext.UserID(ctx)

	slog.Warn("User wants to disable 2FA", "user_id", userID)

	err := r.usersUseCase.Disable2FA(ctx, userID, req.EmailCode)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmailCode):
			return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString("invalid code"),
			}}, nil
		case errors.Is(err, domain.ErrUserNotFound):
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		default:
			slog.Error("failed to disable 2FA", "error", err)

			return nil, r.NewError(ctx, err)
		}
	}

	return &generatedapi.Disable2FANoContent{}, nil
}
