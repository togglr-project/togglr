// Package tokenizer creates and verifies JWT tokens
package tokenizer

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
)

var (
	errUnexpectedMethod = errors.New("unexpected token signing method")
	errInvalidType      = errors.New("invalid token type")
	errInvalidToken     = errors.New("invalid token")
)

type Service struct {
	secretKey        []byte
	accessTTL        time.Duration
	refreshTTL       time.Duration
	resetPasswordTTL time.Duration
}

type ServiceParams struct {
	SecretKey                               []byte
	AccessTTL, RefreshTTL, ResetPasswordTTL time.Duration
}

func New(
	params *ServiceParams,
) *Service {
	return &Service{
		secretKey:        params.SecretKey,
		accessTTL:        params.AccessTTL,
		refreshTTL:       params.RefreshTTL,
		resetPasswordTTL: params.ResetPasswordTTL,
	}
}

func (s *Service) SecretKey() string {
	return string(s.secretKey)
}

func (s *Service) AccessToken(user *domain.User) (string, error) {
	return s.generateToken(user, domain.TokenTypeAccess, s.accessTTL)
}

func (s *Service) RefreshToken(user *domain.User) (string, error) {
	return s.generateToken(user, domain.TokenTypeRefresh, s.refreshTTL)
}

func (s *Service) ResetPasswordToken(user *domain.User) (string, time.Duration, error) {
	token, err := s.generateToken(user, domain.TokenTypeResetPassword, s.resetPasswordTTL)
	if err != nil {
		return "", 0, err
	}

	return token, s.resetPasswordTTL, nil
}

func (s *Service) AccessTokenTTL() time.Duration {
	return s.accessTTL
}

func (s *Service) VerifyToken(token string, tokenType domain.TokenType) (*domain.TokenClaims, error) {
	claims, err := s.verifyToken(token, tokenType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err) //nolint:errorlint // ok
	}

	return claims, nil
}

func (s *Service) verifyToken(token string, tokenType domain.TokenType) (*domain.TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(
		token,
		&domain.TokenClaims{},
		func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errUnexpectedMethod
			}

			return s.secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*domain.TokenClaims)
	if !ok {
		return nil, errInvalidToken
	}

	if claims.TokenType != tokenType {
		return nil, errInvalidType
	}

	return claims, nil
}

func (s *Service) generateToken(user *domain.User, tokenType domain.TokenType, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
		TokenType:   tokenType,
		UserID:      uint(user.ID),
		Username:    user.Username,
		IsSuperuser: user.IsSuperuser,
	})

	return token.SignedString(s.secretKey)
}
