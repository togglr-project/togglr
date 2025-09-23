package apibackend

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
	mockcontract "github.com/togglr-project/togglr/test_mocks/internal_/contract"
)

func TestSecurityHandler_HandleBearerAuth(t *testing.T) {
	t.Run("valid token for other operation", func(t *testing.T) {
		mockTokenizer := mockcontract.NewMockTokenizer(t)
		mockUsersService := mockcontract.NewMockUsersUseCase(t)

		handler := &SecurityHandler{
			tokenizer:    mockTokenizer,
			usersService: mockUsersService,
		}

		ctx := context.Background()
		tokenHolder := generatedapi.BearerAuth{
			Token: "valid_token",
		}

		expectedClaims := &domain.TokenClaims{
			UserID: 123,
		}
		expectedUser := domain.User{
			ID: 123,
		}

		mockTokenizer.EXPECT().
			VerifyToken("valid_token", domain.TokenTypeAccess).
			Return(expectedClaims, nil)

		mockUsersService.EXPECT().
			GetByID(mock.Anything, domain.UserID(123)).
			Return(expectedUser, nil)

		resultCtx, err := handler.HandleBearerAuth(ctx, generatedapi.LoginOperation, tokenHolder)

		require.NoError(t, err)
		require.NotNil(t, resultCtx)
		assert.Equal(t, domain.UserID(123), appcontext.UserID(resultCtx))
	})

	t.Run("invalid token", func(t *testing.T) {
		mockTokenizer := mockcontract.NewMockTokenizer(t)

		handler := &SecurityHandler{
			tokenizer: mockTokenizer,
		}

		ctx := context.Background()
		tokenHolder := generatedapi.BearerAuth{
			Token: "invalid_token",
		}

		expectedErr := errors.New("invalid token")
		mockTokenizer.EXPECT().
			VerifyToken("invalid_token", domain.TokenTypeAccess).
			Return(nil, expectedErr)

		resultCtx, err := handler.HandleBearerAuth(ctx, generatedapi.LoginOperation, tokenHolder)

		require.Error(t, err)
		assert.Nil(t, resultCtx)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockTokenizer := mockcontract.NewMockTokenizer(t)
		mockUsersService := mockcontract.NewMockUsersUseCase(t)

		handler := &SecurityHandler{
			tokenizer:    mockTokenizer,
			usersService: mockUsersService,
		}

		ctx := context.Background()
		tokenHolder := generatedapi.BearerAuth{
			Token: "valid_token",
		}

		expectedClaims := &domain.TokenClaims{
			UserID: 123,
		}
		expectedErr := errors.New("user not found")

		mockTokenizer.EXPECT().
			VerifyToken("valid_token", domain.TokenTypeAccess).
			Return(expectedClaims, nil)

		mockUsersService.EXPECT().
			GetByID(mock.Anything, domain.UserID(123)).
			Return(domain.User{}, expectedErr)

		resultCtx, err := handler.HandleBearerAuth(ctx, generatedapi.LoginOperation, tokenHolder)

		require.Error(t, err)
		assert.Nil(t, resultCtx)
		assert.Equal(t, expectedErr, err)
	})
}
