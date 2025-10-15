package event_bus

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/infra/mq"
)

const (
	TopicSDKErrorReports   = "sdk__error_reports"
	TopicSDKFeedbackEvents = "sdk__feedback_events"
)

const (
	eventsBatchSize = 50
)

type Service struct {
	bus  mq.MQ
	wg   sync.WaitGroup
	stop chan struct{}

	errorReportUseCase contract.ErrorReportsUseCase
	feedbackEventsRepo contract.FeedbackEventsRepository
	algProcessor       contract.AlgorithmsProcessor
}

func New(
	bus mq.MQ,
	errorReportUseCase contract.ErrorReportsUseCase,
	feedbackEventsRepo contract.FeedbackEventsRepository,
	algProcessor contract.AlgorithmsProcessor,
) *Service {
	return &Service{
		bus:                bus,
		stop:               make(chan struct{}),
		errorReportUseCase: errorReportUseCase,
		feedbackEventsRepo: feedbackEventsRepo,
		algProcessor:       algProcessor,
	}
}

//nolint:contextcheck // false positive
func (s *Service) Start(context.Context) error {
	s.dispatchConsumer(TopicSDKErrorReports, s.processSDKErrorReportEvent)
	s.dispatchBatchConsumer(TopicSDKFeedbackEvents, s.processSDKFeedbackEvents)

	return nil
}

func (s *Service) Stop(context.Context) error {
	close(s.stop)
	s.wg.Wait()

	return nil
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

func (s *Service) dispatchBatchConsumer(topic string, processFn func(ctx context.Context, messages [][]byte) error) {
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

		err := s.bus.SubscribeBatch(ctx, topic, eventsBatchSize, processFn)
		if err != nil {
			message := fmt.Sprintf("batch subscribe to topic %q: %v", topic, err)
			slog.Error(message)
		}
	}()
}
