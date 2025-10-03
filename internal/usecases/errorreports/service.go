package errorreports

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

var _ contract.ErrorReportsUseCase = (*Service)(nil)

// Service implements ErrorReportsUseCase.
// It orchestrates saving error reports and performing auto-disable logic.
//
// Dependency directions follow Clean Architecture: UseCase depends on interfaces in contract package.
//
//nolint:structcheck // fields are used across methods
type Service struct {
	txManager           db.TxManager
	repo                contract.ErrorReportRepository
	featuresRepo        contract.FeaturesRepository
	featureParams       contract.FeatureParamsRepository
	featureTags         contract.FeatureTagsRepository
	tagsRepo            contract.TagsRepository
	tagsUC              contract.TagsUseCase
	projectSettingsRepo contract.ProjectSettingsRepository
	projectSettingsUC   contract.ProjectSettingsUseCase
	pendingUC           contract.PendingChangesUseCase
	envsRepo            contract.EnvironmentsRepository
}

func New(
	txManager db.TxManager,
	repo contract.ErrorReportRepository,
	featuresRepo contract.FeaturesRepository,
	featureParams contract.FeatureParamsRepository,
	featureTags contract.FeatureTagsRepository,
	tagsRepo contract.TagsRepository,
	tagsUC contract.TagsUseCase,
	projectSettingsRepo contract.ProjectSettingsRepository,
	projectSettingsUC contract.ProjectSettingsUseCase,
	pendingUC contract.PendingChangesUseCase,
	envsRepo contract.EnvironmentsRepository,
) *Service {
	return &Service{
		txManager:           txManager,
		repo:                repo,
		featuresRepo:        featuresRepo,
		featureParams:       featureParams,
		featureTags:         featureTags,
		tagsRepo:            tagsRepo,
		tagsUC:              tagsUC,
		projectSettingsRepo: projectSettingsRepo,
		projectSettingsUC:   projectSettingsUC,
		pendingUC:           pendingUC,
		envsRepo:            envsRepo,
	}
}

func (s *Service) ReportError(
	ctx context.Context,
	projectID domain.ProjectID,
	featureKey string,
	envKey string,
	reqCtx map[domain.RuleAttribute]any,
	reportType string,
	reportMsg string,
) (bool, error) {
	// find feature by key and env
	feature, err := s.featuresRepo.GetByKeyWithEnv(ctx, featureKey, envKey)
	if err != nil {
		return false, err
	}

	// ensure environment exists to get its ID
	env, err := s.envsRepo.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		return false, err
	}

	// insert error + possibly auto-disable in a single transaction
	var accepted bool

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		// 1) insert error report
		report := domain.ErrorReport{
			EventID:       generateEventID(),
			ProjectID:     projectID,
			FeatureID:     feature.ID,
			EnvironmentID: env.ID,
			ErrorType:     reportType,
			ErrorMessage:  reportMsg,
			Context:       anyMap(reqCtx),
			CreatedAt:     time.Now(),
		}
		if err := s.repo.Insert(txCtx, report); err != nil {
			return err
		}

		// 2) check tag auto-disable
		tagAutoDisable, err := s.tagsUC.GetAutoDisableTag(txCtx, projectID)
		if err != nil {
			// if tag not found, just return health without auto-disable
			slog.Warn("auto-disable tag not found, skipping auto-disable", "error", err)
		} else {
			hasAutoDisableTag, err := s.featureTags.HasFeatureTag(txCtx, feature.ID, tagAutoDisable.ID)
			if err != nil {
				return err
			}
			if hasAutoDisableTag {
				// 3) read project settings
				enabled, err := s.projectSettingsUC.GetAutoDisableEnabled(txCtx, projectID)
				if err != nil {
					return err
				}
				if enabled {
					threshold, err := s.projectSettingsUC.GetAutoDisableErrorThreshold(txCtx, projectID)
					if err != nil {
						return err
					}
					windowSec, err := s.projectSettingsUC.GetAutoDisableTimeWindowSec(txCtx, projectID)
					if err != nil {
						return err
					}

					cnt, err := s.repo.CountRecent(txCtx, feature.ID, env.ID, time.Duration(windowSec)*time.Second)
					if err != nil {
						return err
					}
					if cnt >= threshold {
						// Lock feature_params for update to prevent race conditions
						currentParams, err := s.featureParams.GetForUpdate(txCtx, feature.ID, env.ID)
						if err != nil {
							return fmt.Errorf("get feature params for update: %w", err)
						}

						// Check if already disabled to avoid duplicate operations
						if !currentParams.Enabled {
							// Feature already disabled, skip auto-disable
							slog.Debug("feature already disabled, skipping auto-disable",
								"feature_id", feature.ID, "env_id", env.ID)
						} else {
							requiresApproval, err := s.projectSettingsUC.GetAutoDisableRequiresApproval(txCtx, projectID)
							if err != nil {
								return err
							}
							if !requiresApproval {
								tagGuarded, err := s.tagsUC.GetGuardedTag(txCtx, projectID)
								if err != nil {
									slog.Error("get guarded tag failed", "error", err)
								} else {
									requiresApproval, err = s.featureTags.HasFeatureTag(
										txCtx,
										feature.ID,
										tagGuarded.ID,
									)
									if err != nil {
										return err
									}
								}
							}

							if requiresApproval {
								// create pending change to disable feature
								payload := domain.PendingChangePayload{
									Entities: []domain.EntityChange{
										{
											Entity:   "feature_params",
											EntityID: feature.ID.String(),
											Action:   domain.EntityActionUpdate,
											Changes: map[string]domain.ChangeValue{
												"enabled": {Old: currentParams.Enabled, New: false},
											},
										},
									},
									Meta: domain.PendingChangeMeta{
										Reason: "auto-disable threshold exceeded",
										Origin: "auto-disable",
										Client: "sdk",
									},
								}
								_, err := s.pendingUC.Create(
									txCtx,
									projectID,
									env.ID,
									"system",
									nil,
									payload,
								)
								if err != nil {
									if errors.Is(err, domain.ErrEntityAlreadyExists) {
										return nil
									}

									return fmt.Errorf("create pending change: %w", err)
								}
								accepted = true
							} else {
								// perform immediate disable by updating feature params
								params := domain.FeatureParams{
									FeatureID:     feature.ID,
									EnvironmentID: env.ID,
									Enabled:       false,
									DefaultValue:  currentParams.DefaultValue,
									UpdatedAt:     time.Now(),
								}
								if _, err := s.featureParams.Update(txCtx, projectID, params); err != nil {
									return fmt.Errorf("disable feature: %w", err)
								}

								slog.Warn("feature auto-disabled",
									"feature_id", feature.ID, "env", envKey,
									"error_count", cnt, "threshold", threshold)
							}
						}
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return accepted, nil
}

func (s *Service) GetFeatureHealth(
	ctx context.Context,
	projectID domain.ProjectID,
	featureKey string,
	envKey string,
) (domain.FeatureHealth, error) {
	feature, err := s.featuresRepo.GetByKeyWithEnv(ctx, featureKey, envKey)
	if err != nil {
		return domain.FeatureHealth{}, err
	}
	env, err := s.envsRepo.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		return domain.FeatureHealth{}, err
	}
	// Get time window from project settings, fallback to default
	timeWindow, err := s.projectSettingsUC.GetAutoDisableTimeWindow(ctx, projectID)
	if err != nil {
		// Fallback to default if settings are not available
		timeWindow = 60 * time.Second
	}
	agg, err := s.repo.GetHealth(ctx, feature.ID, env.ID, timeWindow)
	if err != nil {
		return domain.FeatureHealth{}, err
	}
	params, err := s.featureParams.GetByFeatureWithEnv(ctx, feature.ID, env.ID)
	if err != nil {
		return domain.FeatureHealth{}, err
	}

	return domain.FeatureHealth{
		FeatureID:     feature.ID,
		EnvironmentID: env.ID,
		Enabled:       params.Enabled,
		Status:        deriveStatus(params.Enabled, agg.LastErrorAt, timeWindow),
		ErrorRate:     agg.ErrorRate,
		LastErrorAt:   agg.LastErrorAt,
	}, nil
}

func deriveStatus(enabled bool, lastErr time.Time, timeWindow time.Duration) string {
	if !enabled {
		return "disabled"
	}
	if lastErr.IsZero() {
		return "healthy"
	}
	// simplification: if there were recent errors, mark as degraded
	if time.Since(lastErr) < timeWindow {
		return "degraded"
	}

	return "healthy"
}

func anyMap(m map[domain.RuleAttribute]any) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any, len(m))
	for k, v := range m {
		res[string(k)] = v
	}

	return res
}

// generateEventID generates a new UUID for the error report event.
func generateEventID() string {
	return uuid.New().String()
}
