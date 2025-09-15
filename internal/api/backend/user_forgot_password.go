package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) ForgotPassword(
	ctx context.Context,
	req *generatedapi.ForgotPasswordRequest,
) (generatedapi.ForgotPasswordRes, error) {
	if err := r.usersUseCase.ForgotPassword(ctx, req.Email); err != nil {
		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("External user can't change password"),
				},
			}, nil
		}

		slog.Error("forgot password failed", "error", err)

		return nil, err
	}

	return &generatedapi.ForgotPasswordNoContent{}, nil
}
