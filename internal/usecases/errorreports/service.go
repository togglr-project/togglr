package errorreports

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

const (
	// project settings keys.
	psAutoDisableEnabledKey          = "auto_disable_enabled"
	psAutoDisableRequiresApprovalKey = "auto_disable_requires_approval"
	psAutoDisableErrorThresholdKey   = "auto_disable_error_threshold"
	psAutoDisableTimeWindowSecKey    = "auto_disable_time_window_sec"

	// defaults.
	defaultAutoDisableEnabled = true
	defaultRequiresApproval   = false
	defaultErrorThreshold     = 20
	defaultTimeWindow         = 60 * time.Second

	// tag slug which enables auto-disable for a feature.
	autoDisableTagSlug = "auto-disable"
)

var _ contract.ErrorReportsUseCase = (*Service)(nil)

// Service implements ErrorReportsUseCase.
// It orchestrates saving error reports and performing auto-disable logic.
//
// Dependency directions follow Clean Architecture: UseCase depends on interfaces in contract package.
//
//nolint:structcheck // fields are used across methods
type Service struct {
	txManager       db.TxManager
	repo            contract.ErrorReportRepository
	featuresRepo    contract.FeaturesRepository
	featureParams   contract.FeatureParamsRepository
	featureTags     contract.FeatureTagsRepository
	tagsRepo        contract.TagsRepository
	projectSettings contract.ProjectSettingsRepository
	pendingUC       contract.PendingChangesUseCase
	envsRepo        contract.EnvironmentsRepository
}

func New(
	txManager db.TxManager,
	repo contract.ErrorReportRepository,
	featuresRepo contract.FeaturesRepository,
	featureParams contract.FeatureParamsRepository,
	featureTags contract.FeatureTagsRepository,
	tagsRepo contract.TagsRepository,
	projectSettings contract.ProjectSettingsRepository,
	pendingUC contract.PendingChangesUseCase,
	envsRepo contract.EnvironmentsRepository,
) *Service {
	return &Service{
		txManager:       txManager,
		repo:            repo,
		featuresRepo:    featuresRepo,
		featureParams:   featureParams,
		featureTags:     featureTags,
		tagsRepo:        tagsRepo,
		projectSettings: projectSettings,
		pendingUC:       pendingUC,
		envsRepo:        envsRepo,
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
) (domain.FeatureHealth, bool, error) {
	// find feature by key and env
	feature, err := s.featuresRepo.GetByKeyWithEnv(ctx, featureKey, envKey)
	if err != nil {
		return domain.FeatureHealth{}, false, err
	}

	// ensure environment exists to get its ID
	env, err := s.envsRepo.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		return domain.FeatureHealth{}, false, err
	}

	// insert error + possibly auto-disable in a single transaction
	var (
		accepted bool
		health   domain.FeatureHealth
	)

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		// 1) insert error report
		report := domain.ErrorReport{
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
		tag, err := s.tagsRepo.GetByProjectAndSlug(txCtx, projectID, autoDisableTagSlug)
		if err != nil {
			// if tag not found, just return health without auto-disable
			slog.Warn("auto-disable tag not found, skipping auto-disable", "error", err)
		} else {
			hasTag, err := s.featureTags.HasFeatureTag(txCtx, feature.ID, tag.ID)
			if err != nil {
				return err
			}
			if hasTag {
				// 3) read project settings
				enabled := getBoolSetting(txCtx, s.projectSettings, projectID, psAutoDisableEnabledKey, defaultAutoDisableEnabled)
				if enabled {
					threshold := getIntSetting(txCtx, s.projectSettings, projectID, psAutoDisableErrorThresholdKey, defaultErrorThreshold)
					windowSec := getIntSetting(txCtx, s.projectSettings, projectID, psAutoDisableTimeWindowSecKey, int(defaultTimeWindow/time.Second))

					cnt, err := s.repo.CountRecent(txCtx, feature.ID, env.ID, time.Duration(windowSec)*time.Second)
					if err != nil {
						return err
					}
					if cnt >= threshold {
						requiresApproval := getBoolSetting(txCtx, s.projectSettings, projectID, psAutoDisableRequiresApprovalKey, defaultRequiresApproval)
						if requiresApproval {
							// create pending change to disable feature
							payload := domain.PendingChangePayload{
								Entities: []domain.EntityChange{
									{
										Entity:   "feature_params",
										EntityID: feature.ID.String(),
										Action:   domain.EntityActionUpdate,
										Changes: map[string]domain.ChangeValue{
											"enabled": {Old: feature.Enabled, New: false},
										},
									},
								},
								Meta: domain.PendingChangeMeta{
									Reason: "auto-disable threshold exceeded",
									Origin: "auto-disable",
									Client: "sdk",
								},
							}
							_, err := s.pendingUC.Create(txCtx, projectID, env.ID, "system", nil, payload)
							if err != nil {
								return fmt.Errorf("create pending change: %w", err)
							}
							accepted = true
						} else {
							// perform immediate disable by updating feature params
							params := domain.FeatureParams{
								FeatureID:     feature.ID,
								EnvironmentID: env.ID,
								Enabled:       false,
								DefaultValue:  feature.DefaultValue,
							}
							if _, err := s.featureParams.Update(txCtx, projectID, params); err != nil {
								return fmt.Errorf("disable feature: %w", err)
							}
						}
					}
				}
			}
		}

		// 4) build health snapshot (using repo aggregates and current enabled state)
		agg, err := s.repo.GetHealth(txCtx, feature.ID, env.ID, defaultTimeWindow)
		if err != nil {
			return err
		}

		// load latest params for enabled state
		params, err := s.featureParams.GetByFeatureWithEnv(txCtx, feature.ID, env.ID)
		if err != nil {
			return err
		}
		health = domain.FeatureHealth{
			FeatureID:     feature.ID,
			EnvironmentID: env.ID,
			Enabled:       params.Enabled,
			Status:        deriveStatus(params.Enabled, agg.LastErrorAt),
			ErrorRate:     agg.ErrorRate,
			LastErrorAt:   agg.LastErrorAt,
		}

		return nil
	})
	if err != nil {
		return domain.FeatureHealth{}, false, err
	}

	return health, accepted, nil
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
	agg, err := s.repo.GetHealth(ctx, feature.ID, env.ID, defaultTimeWindow)
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
		Status:        deriveStatus(params.Enabled, agg.LastErrorAt),
		ErrorRate:     agg.ErrorRate,
		LastErrorAt:   agg.LastErrorAt,
	}, nil
}

func deriveStatus(enabled bool, lastErr time.Time) string {
	if !enabled {
		return "disabled"
	}
	if lastErr.IsZero() {
		return "healthy"
	}
	// simplification: if there were recent errors, mark as degraded
	if time.Since(lastErr) < defaultTimeWindow {
		return "degraded"
	}

	return "healthy"
}

// helpers for settings.
func getBoolSetting(ctx context.Context, repo contract.ProjectSettingsRepository, projectID domain.ProjectID, name string, def bool) bool {
	st, err := repo.GetByName(ctx, projectID, name)
	if err != nil || st == nil {
		return def
	}
	v, ok := st.Value.(bool)
	if !ok {
		return def
	}

	return v
}

func getIntSetting(ctx context.Context, repo contract.ProjectSettingsRepository, projectID domain.ProjectID, name string, def int) int {
	st, err := repo.GetByName(ctx, projectID, name)
	if err != nil || st == nil {
		return def
	}
	// JSON numbers come as float64 in interface{} commonly
	switch val := st.Value.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return def
	}
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
