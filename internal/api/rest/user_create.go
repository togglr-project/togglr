package rest

import (
	"context"
	"errors"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggl/internal/context"
	"github.com/rom8726/etoggl/internal/domain"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) CreateUser(
	ctx context.Context,
	req *generatedapi.CreateUserRequest,
) (generatedapi.CreateUserRes, error) {
	userID := etogglcontext.UserID(ctx)
	currentUser, err := r.usersUseCase.GetByID(ctx, userID)
	if err != nil {
		slog.Error("get current user failed", "error", err)

		return nil, err
	}

	username := req.Username
	email := req.Email
	password := req.Password

	isSuperuser := false
	if req.IsSuperuser.Set {
		isSuperuser = req.IsSuperuser.Value
	}

	user, err := r.usersUseCase.Create(
		ctx,
		currentUser,
		username,
		email,
		password,
		isSuperuser,
	)
	if err != nil {
		slog.Error("create user failed", "error", err)

		switch {
		case errors.Is(err, domain.ErrPermissionDenied):
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("Only superusers can create new users"),
				},
			}, nil
		case errors.Is(err, domain.ErrUsernameAlreadyInUse):
			return &generatedapi.ErrorBadRequest{
				Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("username already in use"),
				},
			}, nil
		case errors.Is(err, domain.ErrEmailAlreadyInUse):
			return &generatedapi.ErrorBadRequest{
				Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("email already in use"),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.CreateUserResponse{
		User: generatedapi.User{
			ID:              uint(user.ID),
			Username:        user.Username,
			Email:           user.Email,
			IsSuperuser:     user.IsSuperuser,
			IsActive:        user.IsActive,
			IsTmpPassword:   user.IsTmpPassword,
			IsExternal:      user.IsExternal,
			LicenseAccepted: user.LicenseAccepted,
			CreatedAt:       user.CreatedAt,
		},
	}, nil
}
