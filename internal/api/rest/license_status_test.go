package rest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rom8726/etoggl/internal/domain"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
	mockcontract "github.com/rom8726/etoggl/test_mocks/internal_/contract"
)

func TestGetLicenseStatus(t *testing.T) {
	tests := []struct {
		name           string
		licenseStatus  domain.LicenseStatus
		expectedResult *generatedapi.LicenseStatusResponse
		expectError    bool
	}{
		{
			name: "valid trial license",
			licenseStatus: domain.LicenseStatus{
				ID:              "test-id",
				Type:            domain.Trial,
				IssuedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpiresAt:       time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				IsValid:         true,
				IsExpired:       false,
				DaysUntilExpiry: 15,
				LicenseText:     "Trial license Text",
			},
			expectedResult: &generatedapi.LicenseStatusResponse{
				License: generatedapi.LicenseStatusResponseLicense{
					ID:              generatedapi.NewOptString("test-id"),
					Type:            generatedapi.NewOptLicenseType(generatedapi.LicenseTypeTrial),
					IssuedAt:        generatedapi.NewOptDateTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ExpiresAt:       generatedapi.NewOptDateTime(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
					IsValid:         generatedapi.NewOptBool(true),
					IsExpired:       generatedapi.NewOptBool(false),
					DaysUntilExpiry: generatedapi.NewOptInt(15),
					LicenseText:     generatedapi.NewOptString("Trial license Text"),
					Features: []generatedapi.LicenseFeature{
						generatedapi.LicenseFeatureSSO,
						generatedapi.LicenseFeatureLdap,
						generatedapi.LicenseFeatureCorpNotifChannels,
					},
				},
			},
			expectError: false,
		},
		{
			name: "expired commercial license",
			licenseStatus: domain.LicenseStatus{
				ID:              "expired-id",
				Type:            domain.Commercial,
				IssuedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpiresAt:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				IsValid:         false,
				IsExpired:       true,
				DaysUntilExpiry: -5,
				LicenseText:     "Commercial license Text",
			},
			expectedResult: &generatedapi.LicenseStatusResponse{
				License: generatedapi.LicenseStatusResponseLicense{
					ID:              generatedapi.NewOptString("expired-id"),
					Type:            generatedapi.NewOptLicenseType(generatedapi.LicenseTypeCommercial),
					IssuedAt:        generatedapi.NewOptDateTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ExpiresAt:       generatedapi.NewOptDateTime(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
					IsValid:         generatedapi.NewOptBool(false),
					IsExpired:       generatedapi.NewOptBool(true),
					DaysUntilExpiry: generatedapi.NewOptInt(-5),
					LicenseText:     generatedapi.NewOptString("Commercial license Text"),
					Features: []generatedapi.LicenseFeature{
						generatedapi.LicenseFeatureSSO,
						generatedapi.LicenseFeatureLdap,
						generatedapi.LicenseFeatureCorpNotifChannels,
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLicenseUseCase := mockcontract.NewMockLicenseUseCase(t)
			mockLicenseUseCase.EXPECT().GetLicenseStatus(mock.Anything).Return(tt.licenseStatus, nil)

			api := &RestAPI{
				licenseUseCase: mockLicenseUseCase,
			}

			result, err := api.GetLicenseStatus(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
			mockLicenseUseCase.AssertExpectations(t)
		})
	}
}
