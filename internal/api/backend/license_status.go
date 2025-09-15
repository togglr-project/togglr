package apibackend

import (
	"context"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) GetLicenseStatus(ctx context.Context) (generatedapi.GetLicenseStatusRes, error) {
	licenseStatus, err := r.licenseUseCase.GetLicenseStatus(ctx)
	if err != nil {
		slog.Error("get license status failed", "error", err)

		return nil, err
	}

	// Get available features for this license type
	availableFeatures := domain.GetAvailableFeatures(licenseStatus.Type)

	// Convert domain features to API features
	apiFeatures := make([]generatedapi.LicenseFeature, 0, len(availableFeatures))
	for _, feature := range availableFeatures {
		apiFeatures = append(apiFeatures, generatedapi.LicenseFeature(feature))
	}

	return &generatedapi.LicenseStatusResponse{
		License: generatedapi.LicenseStatusResponseLicense{
			ID:              generatedapi.NewOptString(licenseStatus.ID),
			Type:            generatedapi.NewOptLicenseType(generatedapi.LicenseType(licenseStatus.Type)),
			IssuedAt:        generatedapi.NewOptDateTime(licenseStatus.IssuedAt),
			ExpiresAt:       generatedapi.NewOptDateTime(licenseStatus.ExpiresAt),
			IsValid:         generatedapi.NewOptBool(licenseStatus.IsValid),
			IsExpired:       generatedapi.NewOptBool(licenseStatus.IsExpired),
			DaysUntilExpiry: generatedapi.NewOptInt(licenseStatus.DaysUntilExpiry),
			LicenseText:     generatedapi.NewOptString(licenseStatus.LicenseText),
			Features:        apiFeatures,
		},
	}, nil
}
