package rest

import (
	"context"
	"errors"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) UserChangeMyPassword(
	ctx context.Context,
	req *generatedapi.ChangeUserPasswordRequest,
) (generatedapi.UserChangeMyPasswordRes, error) {
	err := r.usersUseCase.UpdatePassword(ctx, etogglcontext.UserID(ctx), req.OldPassword, req.NewPassword)
	if err != nil {
		slog.Error("update password failed", "error", err)

		switch {
		case errors.Is(err, domain.ErrInvalidPassword):
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		case errors.Is(err, domain.ErrPermissionDenied):
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("External users can't change password"),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.UserChangeMyPasswordNoContent{}, nil
}
