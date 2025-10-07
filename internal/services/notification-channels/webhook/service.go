package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
	httpClient *http.Client
	baseURL    string
}

func New(baseURL string) *Service {
	return &Service{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (s *Service) Type() domain.NotificationType {
	return domain.NotificationTypeWebhook
}

func (s *Service) Send(
	ctx context.Context,
	project *domain.Project,
	feature *domain.Feature,
	configData json.RawMessage,
) error {
	var cfg WebhookConfig
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	payload := map[string]any{
		"project_id": project.ID.String(),
		"feature_id": feature.ID.String(),
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.WebhookURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
