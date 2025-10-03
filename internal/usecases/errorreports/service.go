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
	simplecache "github.com/togglr-project/togglr/pkg/simple-cache"
)

const PendingChangeCacheTTL = time.Minute

var _ contract.ErrorReportsUseCase = (*Service)(nil)

type Service struct {
	txManager         db.TxManager
	repo              contract.ErrorReportRepository
	featuresRepo      contract.FeaturesRepository
	featuresUC        contract.FeaturesUseCase
	featureParamsRepo contract.FeatureParamsRepository
	featureTagsRepo   contract.FeatureTagsRepository
	tagsUC            contract.TagsUseCase
	projectSettingsUC contract.ProjectSettingsUseCase
	pendingUC         contract.PendingChangesUseCase
	envsUC            contract.EnvironmentsUseCase

	pendingCache *simplecache.Cache[string, bool]
}

func New(
	txManager db.TxManager,
	repo contract.ErrorReportRepository,
	featuresRepo contract.FeaturesRepository,
	featuresUC contract.FeaturesUseCase,
	featureParams contract.FeatureParamsRepository,
	featureTags contract.FeatureTagsRepository,
	tagsUC contract.TagsUseCase,
	projectSettingsUC contract.ProjectSettingsUseCase,
	pendingUC contract.PendingChangesUseCase,
	envsUC contract.EnvironmentsUseCase,
) *Service {
	pendingCache := simplecache.New[string, bool]()
	pendingCache.StartCleanup(30 * time.Second)

	return &Service{
		txManager:         txManager,
		repo:              repo,
		featuresRepo:      featuresRepo,
		featuresUC:        featuresUC,
		featureParamsRepo: featureParams,
		featureTagsRepo:   featureTags,
		tagsUC:            tagsUC,
		projectSettingsUC: projectSettingsUC,
		pendingUC:         pendingUC,
		envsUC:            envsUC,
		pendingCache:      pendingCache,
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
	feature, err := s.featuresUC.GetByKeyWithEnvCached(ctx, featureKey, envKey)
	if err != nil {
		return false, err
	}

	env, err := s.envsUC.GetByProjectIDAndKeyCached(ctx, projectID, envKey)
	if err != nil {
		return false, err
	}

	var accepted bool

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
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

		tagAutoDisable, err := s.tagsUC.GetAutoDisableTagCached(txCtx, projectID)
		if err != nil {
			slog.Warn("auto-disable tag not found, skipping auto-disable", "error", err)
		} else {
			hasAutoDisableTag, err := s.featureTagsRepo.HasFeatureTag(txCtx, feature.ID, tagAutoDisable.ID)
			if err != nil {
				return err
			}
			if hasAutoDisableTag {
				enabled, err := s.projectSettingsUC.GetAutoDisableEnabledCached(txCtx, projectID)
				if err != nil {
					return err
				}
				if enabled {
					threshold, err := s.projectSettingsUC.GetAutoDisableErrorThresholdCached(txCtx, projectID)
					if err != nil {
						return err
					}
					windowSec, err := s.projectSettingsUC.GetAutoDisableTimeWindowSecCached(txCtx, projectID)
					if err != nil {
						return err
					}

					cnt, err := s.repo.CountRecent(txCtx, feature.ID, env.ID, time.Duration(windowSec)*time.Second)
					if err != nil {
						return err
					}
					if cnt >= threshold {
						currentParams, err := s.featureParamsRepo.GetForUpdate(txCtx, feature.ID, env.ID)
						if err != nil {
							return fmt.Errorf("get feature params for update: %w", err)
						}

						if !currentParams.Enabled {
							slog.Debug("feature already disabled, skipping auto-disable",
								"feature_id", feature.ID, "env_id", env.ID)
						} else {
							requiresApproval, err := s.projectSettingsUC.
								GetAutoDisableRequiresApprovalCached(txCtx, projectID)
							if err != nil {
								return err
							}
							if !requiresApproval {
								tagGuarded, err := s.tagsUC.GetGuardedTagCached(txCtx, projectID)
								if err != nil {
									slog.Error("get guarded tag failed", "error", err)
								} else {
									requiresApproval, err = s.featureTagsRepo.HasFeatureTag(
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
								pendingCacheKey := s.makePendingCacheKey(feature.ID, envKey)
								if _, found := s.pendingCache.Get(pendingCacheKey); found {
									return nil
								}

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
										s.pendingCache.Set(pendingCacheKey, true, PendingChangeCacheTTL)

										return nil
									}

									return fmt.Errorf("create pending change: %w", err)
								}

								s.pendingCache.Set(pendingCacheKey, true, PendingChangeCacheTTL)
								accepted = true
							} else {
								params := domain.FeatureParams{
									FeatureID:     feature.ID,
									EnvironmentID: env.ID,
									Enabled:       false,
									DefaultValue:  currentParams.DefaultValue,
									UpdatedAt:     time.Now(),
								}
								if _, err := s.featureParamsRepo.Update(txCtx, projectID, params); err != nil {
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
	feature, err := s.featuresUC.GetByKeyWithEnvCached(ctx, featureKey, envKey)
	if err != nil {
		return domain.FeatureHealth{}, err
	}
	env, err := s.envsUC.GetByProjectIDAndKeyCached(ctx, projectID, envKey)
	if err != nil {
		return domain.FeatureHealth{}, err
	}
	timeWindow, err := s.projectSettingsUC.GetAutoDisableTimeWindowCached(ctx, projectID)
	if err != nil {
		timeWindow = 60 * time.Second
	}
	agg, err := s.repo.GetHealth(ctx, feature.ID, env.ID, timeWindow)
	if err != nil {
		return domain.FeatureHealth{}, err
	}
	params, err := s.featureParamsRepo.GetByFeatureWithEnv(ctx, feature.ID, env.ID)
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

func generateEventID() string {
	return uuid.New().String()
}

func (s *Service) makePendingCacheKey(featureID domain.FeatureID, envKey string) string {
	return featureID.String() + ":" + envKey
}
