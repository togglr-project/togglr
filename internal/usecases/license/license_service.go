package license

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/internal/license"
	"github.com/rom8726/etoggle/pkg/db"
)

type Service struct {
	txManager       db.TxManager
	licenseRepo     contract.LicensesRepository
	productInfoRepo contract.ProductInfoRepository
	usersRepo       contract.UsersRepository
	settingsSrv     contract.SettingsUseCase
	validator       *license.Validator
}

func New(
	txManager db.TxManager,
	licenseRepo contract.LicensesRepository,
	productInfoRepo contract.ProductInfoRepository,
	usersRepo contract.UsersRepository,
	settingsSrv contract.SettingsUseCase,
) *Service {
	validator, err := license.NewValidator()
	if err != nil {
		validator = nil
	}

	return &Service{
		txManager:       txManager,
		licenseRepo:     licenseRepo,
		productInfoRepo: productInfoRepo,
		usersRepo:       usersRepo,
		settingsSrv:     settingsSrv,
		validator:       validator,
	}
}

func (s *Service) GetLicenseStatus(ctx context.Context) (domain.LicenseStatus, error) {
	lic, err := s.licenseRepo.GetLastByExpiresAt(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.LicenseStatus{}, domain.ErrEntityNotFound
		}

		return domain.LicenseStatus{}, fmt.Errorf("get license: %w", err)
	}

	now := time.Now()
	isExpired := lic.ExpiresAt.Before(now)
	daysUntilExpiry := int(lic.ExpiresAt.Sub(now).Hours() / 24)

	return domain.LicenseStatus{
		ID:              lic.ID,
		Type:            lic.Type,
		IssuedAt:        lic.IssuedAt,
		ExpiresAt:       lic.ExpiresAt,
		IsValid:         !isExpired,
		IsExpired:       isExpired,
		DaysUntilExpiry: daysUntilExpiry,
		LicenseText:     lic.LicenseText,
	}, nil
}

func (s *Service) UpdateLicense(ctx context.Context, licenseText string) (domain.LicenseStatus, error) {
	if s.validator == nil {
		return domain.LicenseStatus{}, errors.New("license validator not available")
	}

	clientID, err := s.productInfoRepo.GetClientID(ctx)
	if err != nil {
		return domain.LicenseStatus{}, fmt.Errorf("get clientID: %w", err)
	}

	validatedLicense, err := s.validator.ValidateLicense(licenseText, clientID)
	if err != nil {
		return domain.LicenseStatus{}, fmt.Errorf("validate license: %w", err)
	}

	lic := domain.License{
		ID:          validatedLicense.ID,
		ClientID:    validatedLicense.ClientID,
		Type:        validatedLicense.Type,
		IssuedAt:    validatedLicense.IssuedAt,
		ExpiresAt:   validatedLicense.ExpiresAt,
		LicenseText: licenseText,
		CreatedAt:   time.Now(),
	}

	var updatedLicense domain.License
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		updatedLicense, err = s.licenseRepo.UpdateLicense(ctx, lic)
		if err != nil {
			return fmt.Errorf("update license: %w", err)
		}

		return nil
	})
	if err != nil {
		return domain.LicenseStatus{}, fmt.Errorf("update license tx: %w", err)
	}

	now := time.Now()
	isExpired := updatedLicense.ExpiresAt.Before(now)
	daysUntilExpiry := int(updatedLicense.ExpiresAt.Sub(now).Hours() / 24)

	return domain.LicenseStatus{
		ID:              updatedLicense.ID,
		Type:            updatedLicense.Type,
		IssuedAt:        updatedLicense.IssuedAt,
		ExpiresAt:       updatedLicense.ExpiresAt,
		IsValid:         !isExpired,
		IsExpired:       isExpired,
		DaysUntilExpiry: daysUntilExpiry,
		LicenseText:     updatedLicense.LicenseText,
	}, nil
}

// IsFeatureAvailable checks if a specific feature is available based on the current license.
func (s *Service) IsFeatureAvailable(ctx context.Context, feature domain.LicenseFeature) (bool, error) {
	licStatus, err := s.GetLicenseStatus(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("get license status: %w", err)
	}

	if licStatus.IsExpired {
		return false, nil
	}

	return domain.IsFeatureAvailable(licStatus.Type, feature), nil
}
