package apibackend

import (
	"context"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

var _ generatedapi.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
	tokenizer    contract.Tokenizer
	usersService contract.UsersUseCase
}

func NewSecurityHandler(
	tokenizer contract.Tokenizer,
	usersService contract.UsersUseCase,
) *SecurityHandler {
	return &SecurityHandler{
		tokenizer:    tokenizer,
		usersService: usersService,
	}
}

func (r *SecurityHandler) HandleBearerAuth(
	ctx context.Context,
	_ generatedapi.OperationName,
	tokenHolder generatedapi.BearerAuth,
) (context.Context, error) {
	claims, err := r.tokenizer.VerifyToken(tokenHolder.Token, domain.TokenTypeAccess)
	if err != nil {
		return nil, err
	}

	user, err := r.usersService.GetByID(ctx, domain.UserID(claims.UserID))
	if err != nil {
		return nil, err
	}

	ctx = appcontext.WithUserID(ctx, user.ID)

	return ctx, nil
}
