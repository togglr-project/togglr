package license

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/togglr-project/togglr/internal/domain"
	mockcontract "github.com/togglr-project/togglr/test_mocks/internal_/contract"
)

func TestGetLicenseStatus(t *testing.T) {
	var tests []struct {
		name           string
		license        domain.License
		expectedStatus domain.LicenseStatus
		expectError    bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mockcontract.NewMockLicensesRepository(t)
			mockRepo.EXPECT().GetLastByExpiresAt(mock.Anything).Return(tt.license, nil)

			service := &Service{
				licenseRepo: mockRepo,
			}

			result, err := service.GetLicenseStatus(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus.ID, result.ID)
			assert.Equal(t, tt.expectedStatus.Type, result.Type)
			assert.Equal(t, tt.expectedStatus.IssuedAt, result.IssuedAt)
			assert.Equal(t, tt.expectedStatus.ExpiresAt, result.ExpiresAt)
			assert.Equal(t, tt.expectedStatus.IsValid, result.IsValid)
			assert.Equal(t, tt.expectedStatus.IsExpired, result.IsExpired)
			assert.NotZero(t, result.DaysUntilExpiry)

			mockRepo.AssertExpectations(t)
		})
	}
}
