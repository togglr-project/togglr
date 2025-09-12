//go:build skip

package users

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/etoggl/internal/domain"
	mockcontract "github.com/rom8726/etoggl/test_mocks/internal_/contract"
	mockusers "github.com/rom8726/etoggl/test_mocks/internal_/usecases/users"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	ssoManager := mockcontract.NewMockSSOProviderManager(t)
	mockLicensesUseCase := mockcontract.NewMockLicenseUseCase(t)
	// Create service
	service := New(
		mockUsersRepo,
		mockTokenizer,
		mockEmailer,
		mockRateLimiter,
		ssoManager,
		mockLicensesUseCase,
		[]AuthProvider{mockAuthProvider},
	)

	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockUsersRepo, service.usersRepo)
	require.Equal(t, mockTokenizer, service.tokenizer)
	require.Equal(t, mockEmailer, service.emailer)
	require.Equal(t, ssoManager, service.ssoManager)

	// Verify that the authProvider is set
	require.NotNil(t, service.authProvider)
}

func TestLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockAuthProvider *mockusers.MockAuthProvider,
			mockTokenizer *mockcontract.MockTokenizer,
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockLicensesUseCase *mockcontract.MockLicenseUseCase,
		)
		username             string
		password             string
		expectedAccessToken  string
		expectedRefreshToken string
		expectedTmpPasswd    bool
		expectedError        bool
		errorContains        string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockLicensesUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicensesUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				user := &domain.User{
					ID:            1,
					Username:      "user1",
					PasswordHash:  "hash",
					IsActive:      true,
					IsTmpPassword: false,
				}

				mockAuthProvider.EXPECT().CanHandle("user1").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user1",
					"password1",
				).Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(user).Return("access_token", nil)
				mockTokenizer.EXPECT().RefreshToken(user).Return("refresh_token", nil)
				mockUsersRepo.EXPECT().UpdateLastLogin(mock.Anything, domain.UserID(1)).Return(nil)
			},
			username:             "user1",
			password:             "password1",
			expectedAccessToken:  "access_token",
			expectedRefreshToken: "refresh_token",
			expectedTmpPasswd:    false,
			expectedError:        false,
		},
		{
			name: "Success with temporary password",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockLicensesUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicensesUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				user := &domain.User{
					ID:            2,
					Username:      "user2",
					PasswordHash:  "hash",
					IsActive:      true,
					IsTmpPassword: true,
				}

				mockAuthProvider.EXPECT().CanHandle("user2").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user2",
					"password2",
				).Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(user).Return("access_token", nil)
				mockTokenizer.EXPECT().RefreshToken(user).Return("refresh_token", nil)
				mockUsersRepo.EXPECT().UpdateLastLogin(mock.Anything, domain.UserID(2)).Return(nil)
			},
			username:             "user2",
			password:             "password2",
			expectedAccessToken:  "access_token",
			expectedRefreshToken: "refresh_token",
			expectedTmpPasswd:    true,
			expectedError:        false,
		},
		{
			name: "Authentication error of external auth provider, local provider succeeded",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockLicensesUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicensesUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				mockAuthProvider.EXPECT().CanHandle("user3").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user3",
					"password3",
				).Return(nil, errors.New("authentication error"))

				user := domain.User{
					ID:            3,
					Username:      "user3",
					PasswordHash:  "$2a$10$v.9vN/U/WBX3oyVIF06e5OEQQUfyFFbm/kWUOYRyzAeYMBuksrrPi",
					IsActive:      true,
					IsTmpPassword: true,
				}
				mockUsersRepo.EXPECT().GetByUsername(mock.Anything, "user3").Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(&user).Return("access_token", nil)
				mockTokenizer.EXPECT().RefreshToken(&user).Return("refresh_token", nil)
				mockUsersRepo.EXPECT().UpdateLastLogin(mock.Anything, domain.UserID(3)).Return(nil)
			},
			username:             "user3",
			password:             "password3",
			expectedAccessToken:  "access_token",
			expectedRefreshToken: "refresh_token",
			expectedTmpPasswd:    true,
			expectedError:        false,
		},
		{
			name: "Access token creation error",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockLicensesUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicensesUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				user := &domain.User{
					ID:           4,
					Username:     "user4",
					PasswordHash: "hash",
					IsActive:     true,
				}

				mockAuthProvider.EXPECT().CanHandle("user4").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user4",
					"password4",
				).Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(user).Return("", errors.New("token error"))
			},
			username:             "user4",
			password:             "password4",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedTmpPasswd:    false,
			expectedError:        true,
			errorContains:        "token error",
		},
		{
			name: "Refresh token creation error",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockLicensesUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicensesUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				user := &domain.User{
					ID:           5,
					Username:     "user5",
					PasswordHash: "hash",
					IsActive:     true,
				}

				mockAuthProvider.EXPECT().CanHandle("user5").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user5",
					"password5",
				).Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(user).Return("access_token", nil)
				mockTokenizer.EXPECT().RefreshToken(user).Return("", errors.New("token error"))
			},
			username:             "user5",
			password:             "password5",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedTmpPasswd:    false,
			expectedError:        true,
			errorContains:        "token error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
			ssoManager := mockcontract.NewMockSSOProviderManager(t)
			mockLicensesUseCase := mockcontract.NewMockLicenseUseCase(t)
			// Setup mocks
			tt.setupMocks(mockAuthProvider, mockTokenizer, mockUsersRepo, mockLicensesUseCase)

			// Create service
			service := New(
				mockUsersRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				ssoManager,
				mockLicensesUseCase,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			accessToken, refreshToken, _, isTmpPasswd, err := service.Login(
				context.Background(),
				tt.username,
				tt.password,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAccessToken, accessToken)
				require.Equal(t, tt.expectedRefreshToken, refreshToken)
				require.Equal(t, tt.expectedTmpPasswd, isTmpPasswd)
			}
		})
	}
}

func TestLoginReissue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTokenizer *mockcontract.MockTokenizer,
			mockUsersRepo *mockcontract.MockUsersRepository,
		)
		refreshToken         string
		expectedAccessToken  string
		expectedRefreshToken string
		expectedError        bool
		errorContains        string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				claims := &domain.TokenClaims{
					UserID: 1,
				}
				mockTokenizer.EXPECT().VerifyToken(
					"valid_refresh_token",
					domain.TokenTypeRefresh,
				).Return(claims, nil)

				user := domain.User{
					ID:       1,
					Username: "user1",
					IsActive: true,
				}
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(&user).Return("new_access_token", nil)
				mockTokenizer.EXPECT().RefreshToken(&user).Return("new_refresh_token", nil)
				mockUsersRepo.EXPECT().UpdateLastLogin(mock.Anything, domain.UserID(1)).Return(nil)
			},
			refreshToken:         "valid_refresh_token",
			expectedAccessToken:  "new_access_token",
			expectedRefreshToken: "new_refresh_token",
			expectedError:        false,
		},
		{
			name: "Invalid refresh token",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				mockTokenizer.EXPECT().VerifyToken(
					"invalid_refresh_token",
					domain.TokenTypeRefresh,
				).Return(nil, errors.New("invalid token"))
			},
			refreshToken:         "invalid_refresh_token",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedError:        true,
			errorContains:        "invalid token",
		},
		{
			name: "User not found",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				claims := &domain.TokenClaims{
					UserID: 999,
				}
				mockTokenizer.EXPECT().VerifyToken(
					"valid_refresh_token",
					domain.TokenTypeRefresh,
				).Return(claims, nil)

				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(999),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			refreshToken:         "valid_refresh_token",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedError:        true,
			errorContains:        "get user by uuid",
		},
		{
			name: "Access token creation error",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				claims := &domain.TokenClaims{
					UserID: 1,
				}
				mockTokenizer.EXPECT().VerifyToken(
					"valid_refresh_token",
					domain.TokenTypeRefresh,
				).Return(claims, nil)

				user := domain.User{
					ID:       1,
					Username: "user1",
					IsActive: true,
				}
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(&user).Return("", errors.New("token error"))
			},
			refreshToken:         "valid_refresh_token",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedError:        true,
			errorContains:        "token error",
		},
		{
			name: "Refresh token creation error",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				claims := &domain.TokenClaims{
					UserID: 1,
				}
				mockTokenizer.EXPECT().VerifyToken(
					"valid_refresh_token",
					domain.TokenTypeRefresh,
				).Return(claims, nil)

				user := domain.User{
					ID:       1,
					Username: "user1",
					IsActive: true,
				}
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(user, nil)

				mockTokenizer.EXPECT().AccessToken(&user).Return("new_access_token", nil)
				mockTokenizer.EXPECT().RefreshToken(&user).Return("", errors.New("token error"))
			},
			refreshToken:         "valid_refresh_token",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedError:        true,
			errorContains:        "token error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
			ssoManager := mockcontract.NewMockSSOProviderManager(t)
			mockLicensesUseCase := mockcontract.NewMockLicenseUseCase(t)
			// Setup mocks
			tt.setupMocks(mockTokenizer, mockUsersRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				ssoManager,
				mockLicensesUseCase,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			accessToken, refreshToken, err := service.LoginReissue(
				context.Background(),
				tt.refreshToken,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAccessToken, accessToken)
				require.Equal(t, tt.expectedRefreshToken, refreshToken)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mockUsersRepo *mockcontract.MockUsersRepository)
		userID        domain.UserID
		expectedUser  domain.User
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:       1,
					Username: "user1",
					IsActive: true,
				}, nil)
			},
			userID: domain.UserID(1),
			expectedUser: domain.User{
				ID:       1,
				Username: "user1",
				IsActive: true,
			},
			expectedError: false,
		},
		{
			name: "User not found",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(999),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			userID:        domain.UserID(999),
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
			ssoManager := mockcontract.NewMockSSOProviderManager(t)
			mockLicensesUseCase := mockcontract.NewMockLicenseUseCase(t)
			// Setup mocks
			tt.setupMocks(mockUsersRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				ssoManager,
				mockLicensesUseCase,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			user, err := service.GetByID(context.Background(), tt.userID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockEmailer *mockcontract.MockEmailer,
			mockLicenseUseCase *mockcontract.MockLicenseUseCase,
		)
		currentUser   domain.User
		username      string
		email         string
		password      string
		isSuperuser   bool
		expectedUser  domain.User
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockEmailer *mockcontract.MockEmailer,
				mockLicenseUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicenseUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				// Check if username exists
				mockUsersRepo.EXPECT().GetByUsername(
					mock.Anything,
					"newuser",
				).Return(domain.User{}, domain.ErrEntityNotFound)

				// Check if email exists
				mockUsersRepo.EXPECT().GetByEmail(
					mock.Anything,
					"newuser@example.com",
				).Return(domain.User{}, domain.ErrEntityNotFound)

				// Create user
				mockUsersRepo.EXPECT().Create(
					mock.Anything,
					mock.AnythingOfType("domain.UserDTO"),
				).Return(domain.User{
					ID:          1,
					Username:    "newuser",
					Email:       "newuser@example.com",
					IsActive:    true,
					IsSuperuser: true,
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "admin",
				IsSuperuser: true,
			},
			username:    "newuser",
			email:       "newuser@example.com",
			password:    "password",
			isSuperuser: true,
			expectedUser: domain.User{
				ID:          1,
				Username:    "newuser",
				Email:       "newuser@example.com",
				IsActive:    true,
				IsSuperuser: true,
				CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: false,
		},
		{
			name: "Current user not superuser",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockEmailer *mockcontract.MockEmailer,
				mockLicenseUseCase *mockcontract.MockLicenseUseCase,
			) {
				// No other mocks needed
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "regular",
				IsSuperuser: false,
			},
			username:      "newuser",
			email:         "newuser@example.com",
			password:      "password",
			isSuperuser:   true,
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "forbidden",
		},
		{
			name: "Username already exists",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockEmailer *mockcontract.MockEmailer,
				mockLicenseUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicenseUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				// Check if username exists
				mockUsersRepo.EXPECT().GetByUsername(
					mock.Anything,
					"existinguser",
				).Return(domain.User{
					ID:       1,
					Username: "existinguser",
				}, nil)
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "admin",
				IsSuperuser: true,
			},
			username:      "existinguser",
			email:         "newuser@example.com",
			password:      "password",
			isSuperuser:   true,
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "username already in use",
		},
		{
			name: "Email already exists",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockEmailer *mockcontract.MockEmailer,
				mockLicenseUseCase *mockcontract.MockLicenseUseCase,
			) {
				mockLicenseUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(domain.LicenseStatus{Type: domain.Commercial}, nil)
				// Check if username exists
				mockUsersRepo.EXPECT().GetByUsername(
					mock.Anything,
					"newuser",
				).Return(domain.User{}, domain.ErrEntityNotFound)

				// Check if email exists
				mockUsersRepo.EXPECT().GetByEmail(
					mock.Anything,
					"existing@example.com",
				).Return(domain.User{
					ID:    2,
					Email: "existing@example.com",
				}, nil)
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "admin",
				IsSuperuser: true,
			},
			username:      "newuser",
			email:         "existing@example.com",
			password:      "password",
			isSuperuser:   true,
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "email already in use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
			ssoManager := mockcontract.NewMockSSOProviderManager(t)
			mockLicenseUseCase := mockcontract.NewMockLicenseUseCase(t)
			// Setup mocks
			tt.setupMocks(mockUsersRepo, mockEmailer, mockLicenseUseCase)

			// Create service
			service := New(
				mockUsersRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				ssoManager,
				mockLicenseUseCase,
				[]AuthProvider{mockAuthProvider},
			)

			// Create context
			ctx := context.Background()

			// Call method
			user, err := service.Create(
				ctx,
				tt.currentUser,
				tt.username,
				tt.email,
				tt.password,
				tt.isSuperuser,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUser, user)
			}
		})
	}
}
