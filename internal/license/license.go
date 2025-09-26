package license

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/crypt"
)

var (
	ErrLicenseExpired     = errors.New("license expired")
	ErrLicenseInvalid     = errors.New("license invalid")
	ErrClientIDMismatch   = errors.New("client_id mismatch")
	ErrTrialAlreadyIssued = errors.New("trial already issued")
	ErrNetworkError       = errors.New("network error")
	ErrNotEmptyDB         = errors.New("trial license can be issued only on empty database")
	ErrRequestTrial       = errors.New("request trial license failed")
)

// License represents a license.
type License struct {
	ID          string                   `json:"id"`
	ClientID    string                   `json:"client_id"`
	Type        domain.LicenseType       `json:"type"`
	Status      domain.LicenseStatusType `json:"status,omitempty"`
	IssuedAt    time.Time                `json:"issued_at"`
	ExpiresAt   time.Time                `json:"expires_at"`
	Hostname    *string                  `json:"hostname,omitempty"`
	MAC         *string                  `json:"mac,omitempty"`
	IP          *string                  `json:"ip,omitempty"`
	Fingerprint *string                  `json:"fingerprint,omitempty"`
	LicenseText string                   `json:"license_text,omitempty"`
}

// LicenseData represents the data that will be signed.
type LicenseData struct {
	License   License `json:"license"`
	Signature []byte  `json:"signature,omitempty"`
}

// IsExpired checks if the license is expired.
func (l *License) IsExpired() bool {
	return time.Now().UTC().After(l.ExpiresAt)
}

// Validate validates the license.
func (l *License) Validate() error {
	if l.ClientID == "" {
		return errors.New("client_id is required")
	}

	if l.Type == "" {
		return errors.New("type is required")
	}

	if l.IssuedAt.IsZero() {
		return errors.New("issued_at is required")
	}

	if l.ExpiresAt.IsZero() {
		return errors.New("expires_at is required")
	}

	if l.IsExpired() {
		return ErrLicenseExpired
	}

	// fingerprint := GetMultiHostFingerprint()
	// if l.Fingerprint != nil && *l.Fingerprint != fingerprint {
	//	return errors.New("host fingerprint mismatch")
	//}

	return nil
}

// Validator is the license validator.
type Validator struct {
	publicKey *rsa.PublicKey

	trialPrivateKey any
	trialPublicKey  *rsa.PublicKey
}

// NewValidator creates a new license validator.
func NewValidator() (*Validator, error) {
	publicKey, err := decodePublicKey(PublicKey)
	if err != nil {
		return nil, err
	}

	trialPublicKey, err := decodePublicKey(TrialPublicKey)
	if err != nil {
		return nil, err
	}

	trialPrivateKey, err := decodePrivateKey(TrialPrivateKey)
	if err != nil {
		return nil, err
	}

	return &Validator{
		publicKey:       publicKey,
		trialPrivateKey: trialPrivateKey,
		trialPublicKey:  trialPublicKey,
	}, nil
}

func (v *Validator) TrialPrivateKey() any {
	return v.trialPrivateKey
}

// ValidateLicense validates a license.
func (v *Validator) ValidateLicense(licenseString, clientID string) (*License, error) {
	// Decode the license
	license, err := v.decodeLicense(licenseString)
	if err != nil {
		return nil, fmt.Errorf("decode license: %w", err)
	}

	// Validate the license
	if err := license.Validate(); err != nil {
		return nil, fmt.Errorf("validate license: %w", err)
	}

	if license.ClientID != clientID {
		return nil, ErrClientIDMismatch
	}

	return license, nil
}

// decodeLicense decodes a license string.
func (v *Validator) decodeLicense(licenseString string) (*License, error) {
	// Base64 decode the license string
	decodedData, err := base64.StdEncoding.DecodeString(licenseString)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	// slog.Debug("License string decoded", "decoded_data", string(decodedData))

	// Decode the license data
	var licenseData LicenseData
	if err := json.Unmarshal(decodedData, &licenseData); err != nil {
		return nil, fmt.Errorf("unmarshal license data: %w", err)
	}

	// Encode the license to JSON for verification
	licenseBytes, err := json.Marshal(licenseData.License)
	if err != nil {
		return nil, fmt.Errorf("marshal license: %w", err)
	}

	// Verify the signature
	var publicKey *rsa.PublicKey
	if licenseData.License.Type == domain.TrialSelfSigned {
		publicKey = v.trialPublicKey
	} else {
		publicKey = v.publicKey
	}

	if err := crypt.Verify(licenseBytes, licenseData.Signature, publicKey); err != nil {
		return nil, fmt.Errorf("verify signature: %w", err)
	}

	return &licenseData.License, nil
}

func decodePublicKey(publicKeyBase64 string) (*rsa.PublicKey, error) {
	// Parse public key from PEM
	publicKey64Decoded, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode public key: %w", err)
	}

	block, _ := pem.Decode(publicKey64Decoded)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	// Check if a public key is of a supported type
	switch rsaPubKey := publicKey.(type) {
	case *rsa.PublicKey:
		return rsaPubKey, nil
	default:
		return nil, fmt.Errorf("unsupported public key type: %T", publicKey)
	}
}

func decodePrivateKey(privateKeyBase64 string) (any, error) {
	privateKey64Decoded, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode private key: %w", err)
	}

	block, _ := pem.Decode(privateKey64Decoded)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	return privateKey, nil
}
