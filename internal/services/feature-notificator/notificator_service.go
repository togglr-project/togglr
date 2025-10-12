package notificator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/rom8726/di"
	"github.com/rom8726/resilience"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

const (
	defaultBatchSize   = 100
	defaultInterval    = time.Second * 10
	defaultWorkerCount = 4
)

var _ di.Servicer = (*Service)(nil)

type notificationResult struct {
	notificationID domain.FeatureNotificationID
	skipped        bool
	skipReason     string
	err            error
}

type Service struct {
	txManager db.TxManager

	channelsMap          map[domain.NotificationType]contract.NotificationChannel
	notificationsUseCase contract.FeatureNotificationsUseCase
	projectsRepo         contract.ProjectsRepository
	envsRepo             contract.EnvironmentsRepository
	featuresUseCase      contract.FeaturesUseCase

	stop chan struct{}

	batchSize   uint
	interval    time.Duration
	workerCount int

	circuitBreaker resilience.CircuitBreaker
}

func New(
	channels []contract.NotificationChannel,
	txManager db.TxManager,
	notificationsUseCase contract.FeatureNotificationsUseCase,
	projectsRepo contract.ProjectsRepository,
	envsRepo contract.EnvironmentsRepository,
	featuresUseCase contract.FeaturesUseCase,
	workerCount int,
) *Service {
	if workerCount <= 0 {
		workerCount = defaultWorkerCount
	}

	channelsMap := make(map[domain.NotificationType]contract.NotificationChannel, len(channels))
	for i := range channels {
		channel := channels[i]
		channelsMap[channel.Type()] = channel
	}

	return &Service{
		channelsMap:          channelsMap,
		txManager:            txManager,
		notificationsUseCase: notificationsUseCase,
		projectsRepo:         projectsRepo,
		envsRepo:             envsRepo,
		featuresUseCase:      featuresUseCase,
		stop:                 make(chan struct{}),
		batchSize:            defaultBatchSize,
		interval:             defaultInterval,
		workerCount:          max2Ints(workerCount, 1),
		circuitBreaker:       resilience.NewDefaultCircuitBreaker("feature-notifications"),
	}
}

// Start starts the worker.
func (s *Service) Start(context.Context) error {
	go s.run() //nolint:contextcheck // it's ok to ignore context check here

	slog.Info("Feature notificator started")

	return nil
}

// Stop stops the worker.
func (s *Service) Stop(context.Context) error {
	close(s.stop)

	return nil
}

// run is the main loop of the worker.
func (s *Service) run() {
	for {
		select {
		case <-s.stop:
			return
		case <-time.After(s.interval):
			s.ProcessOutbox()
		}
	}
}

// ProcessOutbox processes pending notifications in the outbox.
func (s *Service) ProcessOutbox() {
	ctx, cancel := context.WithTimeout(context.Background(), s.interval)
	defer cancel()

	for {
		if ctx.Err() != nil {
			slog.Error("context error", "error", ctx.Err())

			break
		}

		if processed := s.processBatch(ctx); processed == 0 {
			break
		}
	}
}

func (s *Service) processBatch(ctx context.Context) (processed uint) {
	err := func() error {
		sent := 0

		notifications, err := s.notificationsUseCase.TakePendingNotificationsWithSettings(ctx, s.batchSize)
		if err != nil {
			err = fmt.Errorf("take pending notifications: %w", err)

			return err
		}

		if len(notifications) == 0 {
			slog.Debug("no pending feature notifications")

			return nil
		}

		slog.Debug("got pending feature notifications", "count", len(notifications))

		// Create channels for parallel processing
		notificationChan := make(chan *domain.FeatureNotificationWithSettings, len(notifications))
		resultChan := make(chan notificationResult, len(notifications))

		// Start worker goroutines
		var wg sync.WaitGroup
		for range s.workerCount {
			wg.Add(1)
			go s.worker(ctx, &wg, notificationChan, resultChan)
		}

		go func() {
			defer close(resultChan)

			// Send notifications to workers
			for i := range notifications {
				notification := notifications[i]
				notificationChan <- &notification
			}
			close(notificationChan)

			// Wait for all workers to complete
			wg.Wait()
		}()

		// Process results
		for result := range resultChan {
			if result.err != nil {
				slog.Error("check and notify failed",
					"error", result.err, "notification_id", result.notificationID)

				continue
			}

			if result.skipped {
				slog.Debug("notification skipped",
					"notification_id", result.notificationID, "reason", result.skipReason)
				err = s.notificationsUseCase.MarkNotificationAsSkipped(ctx, result.notificationID, result.skipReason)
				if err != nil {
					slog.Error("mark notification as skipped failed",
						"error", err, "notification_id", result.notificationID)
				}
			} else {
				sent++
			}
		}

		if sent > 0 {
			slog.Info("sent notifications", "sent", sent)

			processed = uint(sent)
		}

		return nil
	}()
	if err != nil {
		slog.Error("process feature notifications batch failed", "error", err)
	}

	return processed
}

func (s *Service) checkAndNotify(
	ctx context.Context,
	notification *domain.FeatureNotificationWithSettings,
	envKey string,
) (skipped bool, skipReason string, err error) {
	feature, err := s.featuresUseCase.GetByIDWithEnv(ctx, notification.FeatureID, envKey)
	if err != nil {
		slog.Error("get feature failed", "error", err, "feature_id", notification.FeatureID)
		if errors.Is(err, domain.ErrEntityNotFound) {
			return true, "feature not found", nil
		}

		return false, "", err
	}

	project, err := s.projectsRepo.GetByID(ctx, feature.ProjectID)
	if err != nil {
		slog.Error("get project failed", "error", err, "project_id", feature.ProjectID)
		if errors.Is(err, domain.ErrEntityNotFound) {
			return true, "project not found", nil
		}

		return false, "", err
	}

	settings := filterSettings(notification)
	if len(settings) == 0 {
		return true, "no settings", nil
	}

	for _, setting := range settings {
		channel := s.channelsMap[setting.Type]
		if channel == nil {
			continue
		}

		err := resilience.WithCircuitBreakerAndRetry(
			ctx,
			s.circuitBreaker,
			func(ctx context.Context) error {
				return channel.Send(ctx, &project, &feature, envKey, setting.Config, notification.Payload)
			},
			resilience.DefaultRetryOptions()...,
		)

		if err != nil {
			slog.Error("send notification failed",
				"error", err, "channel", channel.Type())

			err = s.notificationsUseCase.MarkNotificationAsFailed(ctx, notification.ID, err.Error())
			if err != nil {
				slog.Error("mark notification as failed",
					"error", err, "notification_id", notification.ID)
			}
		} else {
			slog.Debug("sent notification",
				"notification_id", notification.ID, "channel", channel.Type())

			err = s.notificationsUseCase.MarkNotificationAsSent(ctx, notification.ID)
			if err != nil {
				slog.Error("mark notification as sent failed",
					"error", err, "notification_id", notification.ID)
			}
		}
	}

	return false, "", nil
}

func (s *Service) worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	notificationChan <-chan *domain.FeatureNotificationWithSettings,
	resultChan chan<- notificationResult,
) {
	defer wg.Done()

	for notification := range notificationChan {
		env, err := s.envsRepo.GetByID(ctx, notification.EnvironmentID)
		if err != nil {
			slog.Error("notifications worker: get environment failed",
				"error", err, "environment_id", notification.EnvironmentID)

			continue
		}

		skipped, skipReason, err := s.checkAndNotify(ctx, notification, env.Key)
		resultChan <- notificationResult{
			notificationID: notification.ID,
			skipped:        skipped,
			skipReason:     skipReason,
			err:            err,
		}
	}
}

func filterSettings(notification *domain.FeatureNotificationWithSettings) []domain.NotificationSetting {
	var availableSettings []domain.NotificationSetting //nolint:prealloc // it's ok

	for _, setting := range notification.Settings {
		if !setting.Enabled {
			continue
		}

		availableSettings = append(availableSettings, setting)
	}

	return availableSettings
}

func max2Ints(a, b int) int {
	if a > b {
		return a
	}

	return b
}
