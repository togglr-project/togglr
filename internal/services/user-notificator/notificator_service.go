package notificator

import (
	"context"
	"encoding/json"
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
	defaultInterval    = time.Minute
	defaultWorkerCount = 4
)

var _ di.Servicer = (*Service)(nil)

type notificationResult struct {
	notificationID domain.UserNotificationID
	skipped        bool
	skipReason     string
	err            error
}

type Service struct {
	userNotificationsRepo contract.UserNotificationsRepository
	emailService          contract.Emailer
	userRepo              contract.UsersRepository
	txManager             db.TxManager

	stopCh chan struct{}

	batchSize   uint
	interval    time.Duration
	workerCount int

	circuitBreaker resilience.CircuitBreaker
}

func New(
	txManager db.TxManager,
	userNotificationsRepo contract.UserNotificationsRepository,
	emailService contract.Emailer,
	userRepo contract.UsersRepository,
	workerCount int,
) *Service {
	if workerCount == 0 {
		workerCount = defaultWorkerCount
	}

	return &Service{
		userNotificationsRepo: userNotificationsRepo,
		emailService:          emailService,
		userRepo:              userRepo,
		txManager:             txManager,
		stopCh:                make(chan struct{}),
		batchSize:             defaultBatchSize,
		interval:              defaultInterval,
		workerCount:           max2Ints(workerCount, 1),
		circuitBreaker:        resilience.NewDefaultCircuitBreaker("user-notifications"),
	}
}

// Start starts the worker.
func (s *Service) Start(context.Context) error {
	go s.run() //nolint:contextcheck // it's ok to ignore context check here

	slog.Info("User notificator started")

	return nil
}

// Stop stops the worker.
func (s *Service) Stop(context.Context) error {
	close(s.stopCh)

	return nil
}

// run is the main loop of the worker.
func (s *Service) run() {
	for {
		select {
		case <-s.stopCh:
			return
		case <-time.After(s.interval):
			s.ProcessOutbox()
		}
	}
}

// ProcessOutbox processes pending email notifications in the outbox.
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
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		sent := 0
		notifications, err := s.userNotificationsRepo.GetPendingEmailNotificationsForUpdate(ctx, s.batchSize)
		if err != nil {
			return fmt.Errorf("get pending email notifications: %w", err)
		}

		if len(notifications) == 0 {
			slog.Debug("no pending email notifications")

			return nil
		}

		slog.Debug("got pending email notifications", "count", len(notifications))

		// Create channels for parallel processing
		notificationChan := make(chan *domain.UserNotification, len(notifications))
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
				err = s.userNotificationsRepo.MarkEmailAsFailed(ctx, result.notificationID, result.skipReason)
				if err != nil {
					slog.Error("mark notification as failed",
						"error", err, "notification_id", result.notificationID)
				}
			} else {
				sent++
			}
		}

		if sent > 0 {
			slog.Info("sent email notifications", "sent", sent)
			processed = uint(sent)
		}

		return nil
	})
	if err != nil {
		slog.Error("process notifications batch failed", "error", err)
	}

	return processed
}

func (s *Service) checkAndNotify(
	ctx context.Context,
	notification *domain.UserNotification,
) (skipped bool, skipReason string) {
	// Get user to get email
	user, err := s.userRepo.GetByID(ctx, notification.UserID)
	if err != nil {
		slog.Error("user not found for notification", "user_id", notification.UserID, "error", err)

		return true, "user not found"
	}
	if user.Email == "" {
		slog.Warn("user has no email, skip notification", "user_id", notification.UserID)

		return true, "user has no email"
	}

	var content domain.UserNotificationContent
	if err := json.Unmarshal(notification.Content, &content); err != nil {
		slog.Error("failed to unmarshal notification content",
			"error", err, "notification_id", notification.ID)

		return true, "invalid content format"
	}

	err = resilience.WithCircuitBreakerAndRetry(
		ctx,
		s.circuitBreaker,
		func(ctx context.Context) error {
			return s.emailService.SendUserNotificationEmail(ctx, user.Email, notification.Type, content)
		},
		resilience.DefaultRetryOptions()...,
	)

	if err != nil {
		slog.Error("send email notification failed",
			"error", err, "notification_id", notification.ID)

		err = s.userNotificationsRepo.MarkEmailAsFailed(ctx, notification.ID, err.Error())
		if err != nil {
			slog.Error("mark notification as failed",
				"error", err, "notification_id", notification.ID)
		}
	} else {
		slog.Debug("sent email notification",
			"notification_id", notification.ID)

		err = s.userNotificationsRepo.MarkEmailAsSent(ctx, notification.ID)
		if err != nil {
			slog.Error("mark notification as sent failed",
				"error", err, "notification_id", notification.ID)
		}
	}

	return false, ""
}

func (s *Service) worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	notificationChan <-chan *domain.UserNotification,
	resultChan chan<- notificationResult,
) {
	defer wg.Done()

	for notification := range notificationChan {
		skipped, skipReason := s.checkAndNotify(ctx, notification)
		resultChan <- notificationResult{
			notificationID: notification.ID,
			skipped:        skipped,
			skipReason:     skipReason,
		}
	}
}

func max2Ints(a, b int) int {
	if a > b {
		return a
	}

	return b
}
