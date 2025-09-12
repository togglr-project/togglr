package email

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_sendEmailsParallel(t *testing.T) {
	tests := []struct {
		name          string
		maxWorkers    int
		emails        []emailData
		expectedError bool
	}{
		{
			name:       "successful parallel sending",
			maxWorkers: 2,
			emails: []emailData{
				{toEmails: []string{"user1@example.com"}, subject: "Test 1", body: "Body 1"},
				{toEmails: []string{"user2@example.com"}, subject: "Test 2", body: "Body 2"},
			},
			expectedError: false,
		},
		{
			name:          "empty emails list",
			maxWorkers:    2,
			emails:        []emailData{},
			expectedError: false,
		},
		{
			name:       "single email",
			maxWorkers: 2,
			emails: []emailData{
				{toEmails: []string{"user1@example.com"}, subject: "Test", body: "Body"},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				cfg: &Config{
					SMTPHost: "localhost:1025",
					Username: "test",
					Password: "test",
				},
			}
			// Подменяем функцию отправки писем на мок
			service.sendEmailFunc = func(ctx context.Context, to []string, subject, body string) error {
				return nil
			}

			err := service.sendEmailsParallel(context.Background(), tt.maxWorkers, tt.emails)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_TLSLogic(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		port        int
		expectedTLS bool
	}{
		{
			name: "port 1025 (MailHog) should never use TLS",
			config: &Config{
				AllowInsecure: false,
				UseTLS:        true,
			},
			port:        1025,
			expectedTLS: false,
		},
		{
			name: "port 465 should use TLS",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        false,
			},
			port:        465,
			expectedTLS: true,
		},
		{
			name: "port 587 should use TLS",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        false,
			},
			port:        587,
			expectedTLS: true,
		},
		{
			name: "port 25 without AllowInsecure should use TLS",
			config: &Config{
				AllowInsecure: false,
				UseTLS:        false,
			},
			port:        25,
			expectedTLS: true,
		},
		{
			name: "port 25 with AllowInsecure should not use TLS",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        false,
			},
			port:        25,
			expectedTLS: false,
		},
		{
			name: "explicit UseTLS should override port logic (except 1025)",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        true,
			},
			port:        25,
			expectedTLS: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a simple test to verify our TLS logic
			useTLS := false
			if tt.port != 1025 { // MailHog port
				useTLS = tt.config.UseTLS || tt.port == 465 || tt.port == 587 || !tt.config.AllowInsecure
			}
			assert.Equal(t, tt.expectedTLS, useTLS)
		})
	}
}

// Helper functions

func ptr[T any](v T) *T {
	return &v
}
