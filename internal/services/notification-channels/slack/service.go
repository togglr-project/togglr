package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type ServiceParams struct {
	BaseURL string
}

type Service struct {
	httpClient *http.Client
	cfg        *ServiceParams
}

type slackMessage struct {
	Channel string       `json:"channel"`
	Blocks  []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string `json:"type"`
	Text struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"text"`
}

func New(cfg *ServiceParams) *Service {
	return &Service{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cfg: cfg,
	}
}

func (s *Service) Type() domain.NotificationType {
	return domain.NotificationTypeSlack
}

func (s *Service) Send(
	ctx context.Context,
	project *domain.Project,
	feature *domain.Feature,
	envKey string,
	configData json.RawMessage,
	payload domain.FeatureNotificationPayload,
) error {
	var cfg SlackConfig
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.WebhookURL == "" {
		return errors.New("webhook URL is required")
	}

	if cfg.ChannelName == "" {
		return errors.New("channel name is required")
	}

	text, err := renderMessage(feature, project, envKey, payload, s.cfg.BaseURL)
	if err != nil {
		return fmt.Errorf("render message: %w", err)
	}

	msg := slackMessage{
		Channel: cfg.ChannelName,
		Blocks: []slackBlock{
			{
				Type: "section",
				Text: struct {
					Type string `json:"type"`
					Text string `json:"text"`
				}{
					Type: "mrkdwn",
					Text: text,
				},
			},
		},
	}

	reqBody, err := json.Marshal(msg)
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
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
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
	const maxMessageLength = 4000

	var msgTemplate string
	var templateData map[string]any

	switch {
	case payload.State != nil:
		status := "disabled"
		if payload.State.Enabled {
			status = "enabled"
		}
		msgTemplate = `*[{{.ProjectName}}] {{.FeatureName}}* ({{.EnvKey}})
Feature {{.Status}} by *{{.ChangedBy}}*`
		templateData = map[string]any{
			"ProjectName": project.Name,
			"FeatureName": feature.Name,
			"EnvKey":      envKey,
			"Status":      status,
			"ChangedBy":   payload.State.ChangedBy,
		}

	case payload.AutoDisabled != nil:
		msgTemplate = `*[{{.ProjectName}}] {{.FeatureName}}* ({{.EnvKey}})
Feature automatically disabled due to error threshold exceeded
Disabled at: {{.DisabledAt}}`
		templateData = map[string]any{
			"ProjectName": project.Name,
			"FeatureName": feature.Name,
			"EnvKey":      envKey,
			"DisabledAt":  payload.AutoDisabled.DisabledAt.Format("2006-01-02 15:04:05"),
		}

	case payload.ChangeRequest != nil:
		msgTemplate = `*[{{.ProjectName}}] {{.FeatureName}}* ({{.EnvKey}})
Change request created
Requested by: *{{.RequestedBy}}*
Approval required to apply changes`
		templateData = map[string]any{
			"ProjectName": project.Name,
			"FeatureName": feature.Name,
			"EnvKey":      envKey,
			"RequestedBy": payload.ChangeRequest.RequestedBy,
		}

	default:
		msgTemplate = `*[{{.ProjectName}}] {{.FeatureName}}* ({{.EnvKey}})`
		templateData = map[string]any{
			"ProjectName": project.Name,
			"FeatureName": feature.Name,
			"EnvKey":      envKey,
		}
	}

	tmpl, err := template.New("slack").Parse(msgTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	message := buf.String()
	if len(message) > maxMessageLength {
		message = message[:maxMessageLength-3] + "..."
	}

	return message, nil
}
