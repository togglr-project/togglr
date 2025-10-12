package feature_notifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

var _ contract.FeatureNotificationsUseCase = (*Service)(nil)

type Service struct {
	txManager                db.TxManager
	notificationSettingsRepo contract.NotificationSettingsRepository
	featureNotificationsRepo contract.FeatureNotificationRepository
	projectsRepo             contract.ProjectsRepository
	envsRepo                 contract.EnvironmentsRepository
	featuresRepo             contract.FeaturesRepository

	notificationChannels []contract.NotificationChannel
}

func New(
	txManager db.TxManager,
	notificationSettingsRepo contract.NotificationSettingsRepository,
	featureNotificationsRepo contract.FeatureNotificationRepository,
	projectsRepo contract.ProjectsRepository,
	envsRepo contract.EnvironmentsRepository,
	featuresRepo contract.FeaturesRepository,
	notificationChannels []contract.NotificationChannel,
) *Service {
	return &Service{
		txManager:                txManager,
		notificationSettingsRepo: notificationSettingsRepo,
		featureNotificationsRepo: featureNotificationsRepo,
		projectsRepo:             projectsRepo,
		envsRepo:                 envsRepo,
		featuresRepo:             featuresRepo,
		notificationChannels:     notificationChannels,
	}
}

// CreateNotificationSetting creates a new notification setting.
func (s *Service) CreateNotificationSetting(
	ctx context.Context,
	settingDTO domain.NotificationSettingDTO,
) (domain.NotificationSetting, error) {
	if _, err := s.projectsRepo.GetByID(ctx, settingDTO.ProjectID); err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("get project by ID: %w", err)
	}

	if settingDTO.Type == domain.NotificationTypeEmail {
		var list, err = s.notificationSettingsRepo.ListSettings(ctx, settingDTO.ProjectID, settingDTO.EnvironmentID)
		if err != nil {
			return domain.NotificationSetting{}, fmt.Errorf("list notification settings: %w", err)
		}

		for _, setting := range list {
			if setting.Type == domain.NotificationTypeEmail {
				return domain.NotificationSetting{}, errors.New("email notification already exists")
			}
		}
	}

	result, err := s.notificationSettingsRepo.CreateSetting(ctx, settingDTO)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("create notification setting: %w", err)
	}

	return result, nil
}

// GetNotificationSetting gets a notification setting by ID.
func (s *Service) GetNotificationSetting(
	ctx context.Context,
	id domain.NotificationSettingID,
) (domain.NotificationSetting, error) {
	setting, err := s.notificationSettingsRepo.GetSettingByID(ctx, id)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("get notification setting: %w", err)
	}

	return setting, nil
}

// UpdateNotificationSetting updates a notification setting.
func (s *Service) UpdateNotificationSetting(
	ctx context.Context,
	setting domain.NotificationSetting,
) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if _, err := s.notificationSettingsRepo.GetSettingByID(ctx, setting.ID); err != nil {
			return fmt.Errorf("get notification setting: %w", err)
		}

		err := s.notificationSettingsRepo.UpdateSetting(ctx, setting)
		if err != nil {
			return fmt.Errorf("update notification setting: %w", err)
		}

		return nil
	})

	return err
}

// DeleteNotificationSetting deletes a notification setting.
func (s *Service) DeleteNotificationSetting(
	ctx context.Context,
	id domain.NotificationSettingID,
) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if _, err := s.notificationSettingsRepo.GetSettingByID(ctx, id); err != nil {
			return fmt.Errorf("get notification setting: %w", err)
		}

		err := s.notificationSettingsRepo.DeleteSetting(ctx, id)
		if err != nil {
			return fmt.Errorf("delete notification setting: %w", err)
		}

		return nil
	})

	return err
}

// ListNotificationSettings lists all notification settings for a project.
func (s *Service) ListNotificationSettings(
	ctx context.Context,
	projectID domain.ProjectID,
	envID domain.EnvironmentID,
) ([]domain.NotificationSetting, error) {
	if _, err := s.projectsRepo.GetByID(ctx, projectID); err != nil {
		return nil, fmt.Errorf("get project by ID: %w", err)
	}

	settings, err := s.notificationSettingsRepo.ListSettings(ctx, projectID, envID)
	if err != nil {
		return nil, fmt.Errorf("list notification settings: %w", err)
	}

	return settings, nil
}

func (s *Service) TakePendingNotificationsWithSettings(
	ctx context.Context,
	limit uint,
) ([]domain.FeatureNotificationWithSettings, error) {
	notifications, err := s.featureNotificationsRepo.TakePendingForUpdate(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("take pending notifications: %w", err)
	}

	projectsMap := make(map[domain.ProjectID]struct{})
	for _, notification := range notifications {
		projectsMap[notification.ProjectID] = struct{}{}
	}

	projectEnvsMap := make(map[projectIDWithEnvID][]domain.NotificationSetting)

	for projectID := range projectsMap {
		settings, err := s.notificationSettingsRepo.ListSettingsAll(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("list notification settings: %w", err)
		}

		for _, setting := range settings {
			projectEnvsMap[projectIDWithEnvID{
				ProjectID: setting.ProjectID,
				EnvID:     setting.EnvironmentID,
			}] = append(projectEnvsMap[projectIDWithEnvID{
				ProjectID: setting.ProjectID,
				EnvID:     setting.EnvironmentID,
			}], setting)
		}
	}

	result := make([]domain.FeatureNotificationWithSettings, 0, len(notifications))
	for i := range notifications {
		notification := notifications[i]
		settings := projectEnvsMap[projectIDWithEnvID{
			ProjectID: notification.ProjectID,
			EnvID:     notification.EnvironmentID,
		}]
		result = append(result, domain.FeatureNotificationWithSettings{
			FeatureNotification: notification,
			Settings:            settings,
		})
	}

	return result, nil
}

func (s *Service) MarkNotificationAsSent(ctx context.Context, id domain.FeatureNotificationID) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.featureNotificationsRepo.MarkAsSent(ctx, id)
	})
}

func (s *Service) MarkNotificationAsFailed(ctx context.Context, id domain.FeatureNotificationID, reason string) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.featureNotificationsRepo.MarkAsFailed(ctx, id, reason)
	})
}

func (s *Service) MarkNotificationAsSkipped(ctx context.Context, id domain.FeatureNotificationID, reason string) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.featureNotificationsRepo.MarkAsSkipped(ctx, id, reason)
	})
}

func (s *Service) SendTestNotification(
	ctx context.Context,
	projectID domain.ProjectID,
	envID domain.EnvironmentID,
	notificationSettingID domain.NotificationSettingID,
) error {
	project, err := s.projectsRepo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project by ID: %w", err)
	}

	settings, err := s.notificationSettingsRepo.GetSettingByID(ctx, notificationSettingID)
	if err != nil {
		return fmt.Errorf("get notification setting: %w", err)
	}

	env, err := s.envsRepo.GetByID(ctx, envID)
	if err != nil {
		return fmt.Errorf("get environment by ID: %w", err)
	}

	for _, channel := range s.notificationChannels {
		if channel.Type() == settings.Type {
			feature := domain.Feature{
				BasicFeature: domain.BasicFeature{
					ID:          domain.FeatureID(uuid.NewString()),
					ProjectID:   projectID,
					Key:         "test_feature",
					Kind:        domain.FeatureKindSimple,
					Name:        "Test Feature",
					Description: "Test Feature",
				},
				EnvironmentID: envID,
				Enabled:       true,
				DefaultValue:  "on",
			}

			payload := domain.FeatureNotificationPayload{
				State: &domain.FeatureNotificationStatePayload{
					Enabled:   true,
					ChangedBy: "test",
				},
			}

			return channel.Send(ctx, &project, &feature, env.Key, settings.Config, payload)
		}
	}

	return errors.New("channel not found")
}
