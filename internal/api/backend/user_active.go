package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) SetUserActiveStatus(
	ctx context.Context,
	req *generatedapi.SetUserActiveStatusRequest,
	params generatedapi.SetUserActiveStatusParams,
) (generatedapi.SetUserActiveStatusRes, error) {
	userID := domain.UserID(params.UserID)
	user, err := r.usersUseCase.SetActiveStatus(ctx, userID, req.IsActive)
	if err != nil {
		slog.Error("set user active status failed", "error", err, "user_id", userID, "is_active", req.IsActive)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		case errors.Is(err, domain.ErrPermissionDenied):
			return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	var lastLogin generatedapi.OptDateTime
	if user.LastLogin != nil {
		lastLogin.Value = *user.LastLogin
		lastLogin.Set = true
	}

	return &generatedapi.User{
		ID:              uint(user.ID),
		Username:        user.Username,
		Email:           user.Email,
		IsSuperuser:     user.IsSuperuser,
		IsActive:        user.IsActive,
		IsTmpPassword:   user.IsTmpPassword,
		IsExternal:      user.IsExternal,
		LicenseAccepted: user.LicenseAccepted,
		CreatedAt:       user.CreatedAt,
		LastLogin:       lastLogin,
	}, nil
}
