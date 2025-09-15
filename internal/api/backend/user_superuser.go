package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) SetSuperuserStatus(
	ctx context.Context,
	req *generatedapi.SetSuperuserStatusRequest,
	params generatedapi.SetSuperuserStatusParams,
) (generatedapi.SetSuperuserStatusRes, error) {
	userID := domain.UserID(params.UserID)
	user, err := r.usersUseCase.SetSuperuserStatus(ctx, userID, req.IsSuperuser)
	if err != nil {
		slog.Error("set superuser status failed",
			"error", err, "user_id", userID, "is_superuser", req.IsSuperuser)

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

	// Convert domain.User to API response format
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
		IsExternal:      user.IsExternal,
		LicenseAccepted: user.LicenseAccepted,
		CreatedAt:       user.CreatedAt,
		LastLogin:       lastLogin,
	}, nil
}
