package event_bus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

func (s *Service) PublishErrorReport(ctx context.Context, event contract.ErrorReportEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal report event: %w", err)
	}

	return s.bus.Publish(ctx, TopicSDKErrorReports, data)
}

func (s *Service) processSDKErrorReportEvent(ctx context.Context, data []byte) error {
	var event contract.ErrorReportEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return domain.NewSkippableError(err)
	}

	requestID := event.RequestID
	ctx = appcontext.WithUsername(ctx, "sdk")
	ctx = appcontext.WithRequestID(ctx, requestID)

	accepted, err := s.errorReportUseCase.ReportError(
		ctx,
		event.ProjectID,
		event.FeatureKey,
		event.EnvKey,
		event.Context,
		event.ErrorType,
		event.ErrorMessage,
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.NewSkippableError(err)
		}

		return err
	}

	if accepted {
		slog.Warn("error report accepted (pending change)",
			"feature_key", event.FeatureKey, "project_id", event.ProjectID,
			"env_key", event.EnvKey, "request_id", requestID)
	}

	slog.Debug("error report processed",
		"feature_key", event.FeatureKey, "project_id", event.ProjectID, "env_key", event.EnvKey)

	return nil
}
