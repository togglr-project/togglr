package rest

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
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

	// Build project permissions for projects where the user has membership
	projectPermissions := generatedapi.UserProjectPermissions{}
	permsByProject, err := r.permissionsService.GetMyProjectPermissions(ctx)
	if err != nil {
		slog.Error("get my project permissions failed", "error", err)
		return nil, err
	}
	for projectID, keys := range permsByProject {
		arr := make([]string, 0, len(keys))
		for _, permKey := range keys {
			arr = append(arr, string(permKey))
		}
		projectPermissions[string(projectID)] = arr
	}

	return &generatedapi.User{
		ID:                 uint(userInfo.ID),
		Username:           userInfo.Username,
		Email:              userInfo.Email,
		IsSuperuser:        userInfo.IsSuperuser,
		IsActive:           userInfo.IsActive,
		IsTmpPassword:      userInfo.IsTmpPassword,
		IsExternal:         userInfo.IsExternal,
		LicenseAccepted:    userInfo.LicenseAccepted,
		TwoFaEnabled:       userInfo.TwoFAEnabled,
		CreatedAt:          userInfo.CreatedAt,
		LastLogin:          lastLogin,
		ProjectPermissions: generatedapi.NewOptUserProjectPermissions(projectPermissions),
	}, nil
}
