package errorreports

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

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
) (domain.FeatureHealth, bool, int, error) {
	// find feature by key and env
	feature, err := s.featuresRepo.GetByKeyWithEnv(ctx, featureKey, envKey)
	if err != nil {
		return domain.FeatureHealth{}, false, 0, err
	}

	// ensure environment exists to get its ID
	env, err := s.envsRepo.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		return domain.FeatureHealth{}, false, 0, err
	}

	// insert error + possibly auto-disable in a single transaction
	var (
		accepted  bool
		health    domain.FeatureHealth
		threshold int
	)

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
				enabled := getBoolSetting(
					txCtx,
					s.projectSettings,
					projectID,
					psAutoDisableEnabledKey,
					defaultAutoDisableEnabled,
				)
				if enabled {
					threshold = getIntSetting(
						txCtx,
						s.projectSettings,
						projectID,
						psAutoDisableErrorThresholdKey,
						defaultErrorThreshold,
					)
					windowSec := getIntSetting(
						txCtx,
						s.projectSettings,
						projectID,
						psAutoDisableTimeWindowSecKey,
						int(defaultTimeWindow/time.Second),
					)

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
							slog.Info("feature already disabled, skipping auto-disable",
								"feature_id", feature.ID, "env_id", env.ID)
						} else {
							requiresApproval := getBoolSetting(
								txCtx,
								s.projectSettings,
								projectID,
								psAutoDisableRequiresApprovalKey,
								defaultRequiresApproval,
							)
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

								slog.Info("feature auto-disabled",
									"feature_id", feature.ID, "error_count", cnt, "threshold", threshold)
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
		return domain.FeatureHealth{}, false, 0, err
	}

	return health, accepted, threshold, nil
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
func getBoolSetting(
	ctx context.Context,
	repo contract.ProjectSettingsRepository,
	projectID domain.ProjectID,
	name string,
	def bool,
) bool {
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

func getIntSetting(
	ctx context.Context,
	repo contract.ProjectSettingsRepository,
	projectID domain.ProjectID,
	name string,
	def int,
) int {
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

// writeAutoDisableAuditLog writes an audit log entry for auto-disable action
// func (s *Service) writeAutoDisableAuditLog(
//	ctx context.Context,
//	projectID domain.ProjectID,
//	featureID domain.FeatureID,
//	envID domain.EnvironmentID,
//	oldParams domain.FeatureParams,
//	newParams domain.FeatureParams,
//	errorCount int,
//	threshold int,
//) error {
//	// Create audit log entry with auto-disable context
//	auditData := map[string]any{
//		"auto_disable": true,
//		"error_count":  errorCount,
//		"threshold":    threshold,
//		"reason":       "auto-disable threshold exceeded",
//		"origin":       "sdk",
//	}
//
//	// Use the audit log writer directly
//	executor := s.getExecutor(ctx)
//	if executor != nil {
//		if err := s.writeAuditLog(ctx, executor, projectID, featureID, envID, oldParams, newParams, auditData); err != nil {
//			return fmt.Errorf("write auto-disable audit log: %w", err)
//		}
//	} else {
//		slog.Warn("no transaction context available, skipping audit log", "feature_id", featureID)
//		return nil // Don't fail the operation
//	}
//
//	return nil
//}

// writeAuditLog is a helper to write audit log entries
// func (s *Service) writeAuditLog(
//	ctx context.Context,
//	executor db.Tx,
//	projectID domain.ProjectID,
//	featureID domain.FeatureID,
//	envID domain.EnvironmentID,
//	oldVal any,
//	newVal any,
//	meta map[string]any,
//) error {
//	// Import auditlog package at the top of the file
//	// For now, we'll create a simple audit log entry
//	const query = `
//		INSERT INTO audit_log (project_id, feature_id, entity_id, entity, actor,
//		                       username, action, old_value, new_value, request_id, environment_id)
//		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
//	`
//
//	// Marshal old and new values
//	var oldJSON, newJSON []byte
//	var err error
//
//	if oldVal != nil {
//		oldJSON, err = json.Marshal(oldVal)
//		if err != nil {
//			return fmt.Errorf("marshal old value: %w", err)
//		}
//	}
//
//	if newVal != nil {
//		newJSON, err = json.Marshal(newVal)
//		if err != nil {
//			return fmt.Errorf("marshal new value: %w", err)
//		}
//	}
//
//	// Add metadata to new value
//	if meta != nil && newJSON != nil {
//		var newValMap map[string]any
//		if err := json.Unmarshal(newJSON, &newValMap); err == nil {
//			for k, v := range meta {
//				newValMap[k] = v
//			}
//			newJSON, err = json.Marshal(newValMap)
//			if err != nil {
//				return fmt.Errorf("marshal new value with metadata: %w", err)
//			}
//		}
//	}
//
//	_, err = executor.Exec(
//		ctx,
//		query,
//		projectID,
//		featureID,
//		featureID.String(), // entity_id
//		"feature_params",   // entity
//		"system",           // actor
//		"",                 // username
//		"auto_disable",     // action
//		oldJSON,
//		newJSON,
//		"", // request_id
//		int64(envID),
//	)
//
//	return err
//}

// getExecutor is a helper to get executor from context
// func (s *Service) getExecutor(ctx context.Context) db.Tx {
//	// Get executor from transaction context if available
//	if tx := db.TxFromContext(ctx); tx != nil {
//		return tx
//	}
//	// If no transaction context available, return nil
//	// This should not happen in normal flow as we're always within a transaction
//	return nil
//}

// generateEventID generates a new UUID for error report event.
func generateEventID() string {
	return uuid.New().String()
}
