package rest

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggl/internal/context"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) GetCurrentUser(ctx context.Context) (generatedapi.GetCurrentUserRes, error) {
	userInfo, err := r.usersUseCase.GetByID(ctx, etogglcontext.UserID(ctx))
	if err != nil {
		slog.Error("get current user info failed", "error", err)

		return nil, err
	}

	var lastLogin generatedapi.OptDateTime
	if userInfo.LastLogin != nil {
		lastLogin.Value = *userInfo.LastLogin
		lastLogin.Set = true
	}

	return &generatedapi.User{
		ID:              uint(userInfo.ID),
		Username:        userInfo.Username,
		Email:           userInfo.Email,
		IsSuperuser:     userInfo.IsSuperuser,
		IsActive:        userInfo.IsActive,
		IsTmpPassword:   userInfo.IsTmpPassword,
		IsExternal:      userInfo.IsExternal,
		LicenseAccepted: userInfo.LicenseAccepted,
		TwoFaEnabled:    userInfo.TwoFAEnabled,
		CreatedAt:       userInfo.CreatedAt,
		LastLogin:       lastLogin,
	}, nil
}
