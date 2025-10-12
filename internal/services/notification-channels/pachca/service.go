package pachca

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/togglr-project/togglr/internal/domain"
)

//go:embed message.tmpl
var messageTmpl string

type ServiceParams struct {
	BaseURL string
}

type Service struct {
	httpClient *http.Client
	cfg        *ServiceParams
}

func New(cfg *ServiceParams) *Service {
	return &Service{
		httpClient: &http.Client{},
		cfg:        cfg,
	}
}

func (s *Service) Type() domain.NotificationType {
	return domain.NotificationTypePachca
}

func (s *Service) Send(
	ctx context.Context,
	project *domain.Project,
	feature *domain.Feature,
	envKey string,
	configData json.RawMessage,
	payload domain.FeatureNotificationPayload,
) error {
	var cfg PachcaConfig
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	message, err := renderMessage(feature, project, envKey, payload, s.cfg.BaseURL)
	if err != nil {
		return fmt.Errorf("render message: %w", err)
	}

	reqBody, err := json.Marshal(map[string]any{
		"message": message,
	})
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

func renderMessage(
	feature *domain.Feature,
	project *domain.Project,
	envKey string,
	payload domain.FeatureNotificationPayload,
	_ string,
) (string, error) {
	tmpl, err := template.New("pachca").Parse(messageTmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var templateData map[string]any

	switch {
	case payload.State != nil:
		status := "disabled"
		if payload.State.Enabled {
			status = "enabled"
		}
		templateData = map[string]any{
			"ProjectName": project.Name,
			"FeatureName": feature.Name,
			"EnvKey":      envKey,
			"State":       true,
			"Status":      status,
			"ChangedBy":   payload.State.ChangedBy,
		}

	case payload.AutoDisabled != nil:
		templateData = map[string]any{
			"ProjectName":  project.Name,
			"FeatureName":  feature.Name,
			"EnvKey":       envKey,
			"AutoDisabled": true,
			"DisabledAt":   payload.AutoDisabled.DisabledAt.Format("2006-01-02 15:04:05"),
		}

	case payload.ChangeRequest != nil:
		templateData = map[string]any{
			"ProjectName":   project.Name,
			"FeatureName":   feature.Name,
			"EnvKey":        envKey,
			"ChangeRequest": true,
			"RequestedBy":   payload.ChangeRequest.RequestedBy,
		}

	default:
		templateData = map[string]any{
			"ProjectName": project.Name,
			"FeatureName": feature.Name,
			"EnvKey":      envKey,
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}
