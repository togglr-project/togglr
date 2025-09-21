package license

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	serverURL string
	client    *http.Client
}

type LicenseResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Data    LicenseResponseData `json:"data,omitempty"`
}

type LicenseResponseData struct {
	// License ID.
	ID string `json:"id"`
	// Base64-encoded license string.
	LicenseString string `json:"license_string"`
	// When the license was issued.
	IssuedAt time.Time `json:"issued_at"`
	// When the license expires.
	ExpiresAt time.Time `json:"expires_at"`
	// Type of license.
	LicenseType string `json:"license_type"`
}

func NewClient(serverURL string) *Client {
	if os.Getenv("ENVIRONMENT") == "test" {
		serverURL = "https://example.com"
	}

	return &Client{
		serverURL: serverURL,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) RequestTrialLicense(
	ctx context.Context,
	clientID, hostname, mac, ipAddr, fingerprint string,
) (string, error) {
	reqBody := map[string]interface{}{
		"client_id":   clientID,
		"hostname":    hostname,
		"mac":         mac,
		"ip":          ipAddr,
		"fingerprint": fingerprint,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost, c.serverURL+"/api/license/trial", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w. send request: %w", ErrNetworkError, err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		bodyData, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("read response body: %w", err)
		}

		// slog.Debug("Received trial license", "body", string(bodyData))

		var respBody LicenseResponse
		if err := json.Unmarshal(bodyData, &respBody); err != nil {
			return "", fmt.Errorf("unmarshal response body: %w", err)
		}

		if !respBody.Success {
			slog.Error("Failed to request trial license", "error_message", respBody.Message)

			if strings.Contains(respBody.Message, "already exists") {
				return "", ErrTrialAlreadyIssued
			}

			return "", fmt.Errorf("%w. request trial license: %s", ErrRequestTrial, respBody.Message)
		}

		// slog.Debug("Received trial license string", "license_string", respBody.Data.LicenseString)

		return respBody.Data.LicenseString, nil
	case http.StatusConflict:
		return "", ErrTrialAlreadyIssued
	default:
		return "", fmt.Errorf("%w: unexpected status code %d", ErrNetworkError, resp.StatusCode)
	}
}
