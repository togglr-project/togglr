package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type LicenseUseCase interface {
	GetLicenseStatus(ctx context.Context) (domain.LicenseStatus, error)
	UpdateLicense(ctx context.Context, licenseText string) (domain.LicenseStatus, error)
	IsFeatureAvailable(ctx context.Context, feature domain.LicenseFeature) (bool, error)
}

type LicensesRepository interface {
	GetLastByExpiresAt(ctx context.Context) (domain.License, error)
	GetByID(ctx context.Context, id string) (domain.License, error)
	Create(ctx context.Context, license domain.License) (domain.License, error)
	UpdateLicense(ctx context.Context, license domain.License) (domain.License, error)
}

type ProductInfoRepository interface {
	GetClientID(ctx context.Context) (string, error)
}

type ProductInfoUseCase interface {
	GetProductInfo(ctx context.Context) (domain.ProductInfo, error)
}
