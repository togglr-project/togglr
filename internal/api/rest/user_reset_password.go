package rest

import (
	"context"
	"errors"

	"github.com/rom8726/etoggl/internal/domain"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) ResetPassword(
	ctx context.Context,
	req *generatedapi.ResetPasswordRequest,
) (generatedapi.ResetPasswordRes, error) {
	if err := r.usersUseCase.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			return &generatedapi.ErrorUnauthorized{
				Error: generatedapi.ErrorUnauthorizedError{
					Message: generatedapi.NewOptString(domain.ErrInvalidToken.Error()),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.ResetPasswordNoContent{}, nil
}
