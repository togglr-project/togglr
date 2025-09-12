package license

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/rom8726/etoggl/internal/domain"
	"github.com/rom8726/etoggl/pkg/crypt"
)

const selfSignedTrialTTL = time.Hour * 24 * 60

type Middleware struct {
	validator       *Validator
	client          *Client
	repo            LicensesRepository
	productInfoRepo ProductInfoRepository
	projectsRepo    ProjectsRepository

	clientID string
}

func NewMiddleware(
	repo LicensesRepository,
	productInfoRepo ProductInfoRepository,
	projectsRepo ProjectsRepository,
	licenseServerURL string,
) (*Middleware, error) {
	validator, err := NewValidator()
	if err != nil {
		return nil, fmt.Errorf("create validator: %w", err)
	}

	clientID, err := productInfoRepo.GetClientID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client ID: %w", err)
	}

	client := NewClient(licenseServerURL)

	return &Middleware{
		repo:            repo,
		productInfoRepo: productInfoRepo,
		projectsRepo:    projectsRepo,
		validator:       validator,
		client:          client,
		clientID:        clientID,
	}, nil
}

func (m *Middleware) ValidateLicense(ctx context.Context) error {
	license, err := m.repo.GetLastByExpiresAt(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			slog.Info("License not found in repository, requesting trial license")

			return m.requestTrialLicense(ctx)
		}

		return fmt.Errorf("get last license: %w", err)
	}

	if license.ExpiresAt.Before(time.Now()) {
		slog.Error("License is expired",
			"license_id", license.ID,
			"expires_at", license.ExpiresAt)

		return ErrLicenseExpired
	}

	_, err = m.validator.ValidateLicense(license.LicenseText, m.clientID)
	if err != nil {
		slog.Error("Failed to validate license", "error", err)

		return ErrLicenseInvalid
	}

	slog.Info("License validated successfully",
		"license_id", license.ID,
		"type", license.Type,
		"expires_at", license.ExpiresAt,
	)

	return nil
}

//nolint:nestif // need refactoring
func (m *Middleware) requestTrialLicense(ctx context.Context) error {
	projectsCnt, err := m.projectsRepo.Count(ctx)
	if err == nil {
		if projectsCnt > 0 {
			return ErrNotEmptyDB
		}
	} else {
		slog.Error("Failed to get projects count", "error", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		slog.Warn("Failed to get hostname", "error", err)
		hostname = ""
	}

	mac := ""                                // GetMACAddress()
	ipAddr := ""                             // GetIPAddress()
	fingerprint := GetMultiHostFingerprint() // GetSystemFingerprint()

	licenseString, err := m.client.RequestTrialLicense(
		ctx, m.clientID, hostname, mac, ipAddr, fingerprint)
	if err != nil {
		slog.Error("Failed to request trial license from license server", "error", err)

		if errors.Is(err, ErrNetworkError) || errors.Is(err, ErrRequestTrial) {
			licenseString, err = m.issueSelfSignedTrialLicense()
			if err != nil {
				slog.Error("Failed to issue self-signed trial license", "error", err)

				return err
			}
		} else {
			return fmt.Errorf("request trial license: %w", err)
		}
	}

	license, err := m.validator.ValidateLicense(licenseString, m.clientID)
	if err != nil {
		slog.Error("Failed to validate received trial license", "error", err)

		return fmt.Errorf("validate received license: %w", err)
	}

	domainLicense := domain.License{
		ID:          license.ID,
		ClientID:    license.ClientID,
		Type:        license.Type,
		IssuedAt:    license.IssuedAt,
		ExpiresAt:   license.ExpiresAt,
		LicenseText: licenseString,
		CreatedAt:   time.Now(),
	}

	_, err = m.repo.Create(ctx, domainLicense)
	if err != nil {
		slog.Error("Failed to save trial license to repository", "error", err)

		return fmt.Errorf("save trial license: %w", err)
	}

	slog.Info("Trial license requested and saved successfully",
		"license_id", license.ID,
		"client_id", license.ClientID)

	return nil
}

func (m *Middleware) issueSelfSignedTrialLicense() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Warn("Failed to get hostname", "error", err)
	}

	fingerprint := GetMultiHostFingerprint()

	license := License{
		ID:          uuid.NewString(),
		ClientID:    m.clientID,
		Type:        domain.TrialSelfSigned,
		Status:      "",
		IssuedAt:    time.Now().Truncate(time.Second).UTC(),
		ExpiresAt:   time.Now().Add(selfSignedTrialTTL).Truncate(time.Second).UTC(),
		Hostname:    &hostname,
		Fingerprint: &fingerprint,
	}

	licenseBytes, err := json.Marshal(license)
	if err != nil {
		return "", fmt.Errorf("marshal license: %w", err)
	}

	privateKey := m.validator.TrialPrivateKey()

	cipherText, err := crypt.Sign(licenseBytes, privateKey)
	if err != nil {
		return "", fmt.Errorf("encrypt license: %w", err)
	}

	licenseData := LicenseData{
		License:   license,
		Signature: cipherText,
	}

	licenseDataBytes, err := json.Marshal(licenseData)
	if err != nil {
		return "", fmt.Errorf("marshal license data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(licenseDataBytes), nil
}
