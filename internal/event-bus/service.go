package event_bus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/infra/mq"
)

const (
	topicSDKErrorReports = "sdk.error_reports"
)

type Service struct {
	bus  mq.MQ
	wg   sync.WaitGroup
	stop chan struct{}

	errorReportUseCase contract.ErrorReportsUseCase
}

func New(bus mq.MQ, errorReportUseCase contract.ErrorReportsUseCase) *Service {
	return &Service{
		bus:                bus,
		stop:               make(chan struct{}),
		errorReportUseCase: errorReportUseCase,
	}
}

func (s *Service) Start(context.Context) error {
	s.dispatchConsumer(topicSDKErrorReports, s.processSDKErrorReportEvent)

	return nil
}

func (s *Service) Stop(context.Context) error {
	close(s.stop)
	s.wg.Wait()

	return nil
}

func (s *Service) PublishErrorReport(ctx context.Context, event contract.ErrorReportEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return s.bus.Publish(ctx, topicSDKErrorReports, data)
}

func (s *Service) dispatchConsumer(topic string, processFn func(ctx context.Context, data []byte) error) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				slog.Error(fmt.Sprintf("panic: %v", r))
			}
		}()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			<-s.stop
			cancel()
		}()

		err := s.bus.Subscribe(ctx, topic, processFn)
		if err != nil {
			message := fmt.Sprintf("subscribe to topic %q: %v", topic, err)
			slog.Error(message)
		}
	}()
}

func (s *Service) processSDKErrorReportEvent(ctx context.Context, data []byte) error {
	var event contract.ErrorReportEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return domain.NewSkippableError(err)
	}

	_, accepted, _, err := s.errorReportUseCase.ReportError(
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
			"feature_key", event.FeatureKey, "project_id", event.ProjectID, "env_key", event.EnvKey)
	}

	slog.Debug("error report processed",
		"feature_key", event.FeatureKey, "project_id", event.ProjectID, "env_key", event.EnvKey)

	return nil
}
