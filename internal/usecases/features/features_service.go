package features

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Service struct {
	txManager             db.TxManager
	repo                  contract.FeaturesRepository
	flagVariantsRep       contract.FlagVariantsRepository
	rulesRep              contract.RulesRepository
	schedulesRep          contract.FeatureSchedulesRepository
	featureParamsRep      contract.FeatureParamsRepository
	environmentsRep       contract.EnvironmentsRepository
	guardService          contract.GuardService
	pendingChangesUseCase contract.PendingChangesUseCase
}

func New(
	txManager db.TxManager,
	repo contract.FeaturesRepository,
	flagVariantsRep contract.FlagVariantsRepository,
	rulesRep contract.RulesRepository,
	schedulesRep contract.FeatureSchedulesRepository,
	featureParamsRep contract.FeatureParamsRepository,
	environmentsRep contract.EnvironmentsRepository,
	guardService contract.GuardService,
	pendingChangesUseCase contract.PendingChangesUseCase,
) *Service {
	return &Service{
		txManager:             txManager,
		repo:                  repo,
		flagVariantsRep:       flagVariantsRep,
		rulesRep:              rulesRep,
		schedulesRep:          schedulesRep,
		featureParamsRep:      featureParamsRep,
		environmentsRep:       environmentsRep,
		guardService:          guardService,
		pendingChangesUseCase: pendingChangesUseCase,
	}
}

func (s *Service) CreateWithChildren(
	ctx context.Context,
	feature domain.Feature,
	variants []domain.FlagVariant,
	rules []domain.Rule,
) (domain.FeatureExtended, error) {
	var result domain.FeatureExtended

	envs, err := s.environmentsRep.ListByProjectID(ctx, feature.ProjectID)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("get env: %w", err)
	}

	envProd := envs[0]
	for _, env := range envs {
		if env.Key == "prod" {
			envProd = env

			break
		}
	}

	ruleVariantsMap := make(map[domain.FlagVariantID][]domain.RuleID)
	rulesEnvMap := make(map[string][]*domain.Rule)
	for _, env := range envs {
		rulesEnv := make([]*domain.Rule, 0, len(rules))
		for _, rule := range rules {
			ruleNew := rule
			ruleNew.ID = domain.RuleID(uuid.NewString())
			rulesEnv = append(rulesEnv, &ruleNew)
			if rule.FlagVariantID != nil {
				ruleIDs := ruleVariantsMap[*rule.FlagVariantID]
				ruleIDs = append(ruleIDs, ruleNew.ID)
				ruleVariantsMap[*rule.FlagVariantID] = ruleIDs
			}
		}

		rulesEnvMap[env.Key] = rulesEnv
	}

	variantsEnvMap := make(map[string][]domain.FlagVariant)
	for _, env := range envs {
		variantEnv := make([]domain.FlagVariant, 0, len(variants))
		for _, variant := range variants {
			variantNew := variant
			variantNew.ID = domain.FlagVariantID(uuid.New().String())
			variantEnv = append(variantEnv, variantNew)

			if ruleIDs, ok := ruleVariantsMap[variant.ID]; ok {
				rulesEnv := rulesEnvMap[env.Key]
				for _, ruleEnv := range rulesEnv {
					for _, ruleID := range ruleIDs {
						if ruleEnv.ID == ruleID {
							ruleEnv.FlagVariantID = &variantNew.ID
						}
					}
				}
			}
		}

		variantsEnvMap[env.Key] = variantEnv
	}

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Create feature first
		createdFeature, err := s.repo.Create(ctx, envProd.ID, feature.BasicFeature)
		if err != nil {
			return fmt.Errorf("create feature: %w", err)
		}
		result.BasicFeature = createdFeature

		for _, env := range envs {
			// Create feature params for the environment
			featureParams := domain.FeatureParams{
				FeatureID:     createdFeature.ID,
				EnvironmentID: env.ID,
				Enabled:       feature.Enabled,
				DefaultValue:  feature.DefaultValue,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			_, err = s.featureParamsRep.Create(ctx, feature.ProjectID, featureParams)
			if err != nil {
				return fmt.Errorf("create feature params: %w", err)
			}

			// Create variants
			variantsEnv := variantsEnvMap[env.Key]
			createdVariants := make([]domain.FlagVariant, 0, len(variantsEnv))
			for _, variant := range variantsEnv {
				variant.EnvironmentID = env.ID
				variant.ProjectID = createdFeature.ProjectID
				variant.FeatureID = createdFeature.ID
				cv, err := s.flagVariantsRep.Create(ctx, variant)
				if err != nil {
					return fmt.Errorf("create flag variant: %w", err)
				}
				createdVariants = append(createdVariants, cv)
			}
			result.FlagVariants = createdVariants

			// Create rules
			rulesEnv := rulesEnvMap[env.Key]
			createdRules := make([]domain.Rule, 0, len(rulesEnv))
			for _, rule := range rulesEnv {
				rule.EnvironmentID = env.ID
				rule.FeatureID = createdFeature.ID
				rule.ProjectID = createdFeature.ProjectID
				cr, err := s.rulesRep.Create(ctx, *rule)
				if err != nil {
					return fmt.Errorf("create rule: %w", err)
				}
				createdRules = append(createdRules, cr)
			}
			result.Rules = createdRules
		}

		return nil
	}); err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("tx create feature with children: %w", err)
	}

	return result, nil
}

func (s *Service) GetByIDWithEnv(ctx context.Context, id domain.FeatureID, envKey string) (domain.Feature, error) {
	f, err := s.repo.GetByIDWithEnv(ctx, id, envKey)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("get feature by id with environment: %w", err)
	}

	return f, nil
}

func (s *Service) GetExtendedByID(
	ctx context.Context,
	id domain.FeatureID,
	envKey string,
) (domain.FeatureExtended, error) {
	feature, err := s.repo.GetByIDWithEnv(ctx, id, envKey)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("get feature by id with environment: %w", err)
	}

	env, err := s.environmentsRep.GetByProjectIDAndKey(ctx, feature.ProjectID, envKey)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("get environment: %w", err)
	}

	variants, err := s.flagVariantsRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("list flag variants: %w", err)
	}

	rules, err := s.rulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("list rules: %w", err)
	}

	schedules, err := s.schedulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("list schedules: %w", err)
	}

	return domain.FeatureExtended{
		Feature:      feature,
		FlagVariants: variants,
		Rules:        rules,
		Schedules:    schedules,
	}, nil
}

func (s *Service) GetByKeyWithEnv(ctx context.Context, key, envKey string) (domain.Feature, error) {
	f, err := s.repo.GetByKeyWithEnv(ctx, key, envKey)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("get feature by key with environment: %w", err)
	}

	return f, nil
}

func (s *Service) List(ctx context.Context, envKey string) ([]domain.Feature, error) {
	items, err := s.repo.List(ctx, envKey)
	if err != nil {
		return nil, fmt.Errorf("list features: %w", err)
	}

	return items, nil
}

func (s *Service) ListByProjectID(
	ctx context.Context,
	projectID domain.ProjectID,
	envKey string,
) ([]domain.Feature, error) {
	items, err := s.repo.ListByProjectID(ctx, projectID, envKey)
	if err != nil {
		return nil, fmt.Errorf("list features by projectID: %w", err)
	}

	return items, nil
}

func (s *Service) ListByProjectIDFiltered(
	ctx context.Context,
	projectID domain.ProjectID,
	envKey string,
	filter contract.FeaturesListFilter,
) ([]domain.Feature, int, error) {
	items, total, err := s.repo.ListByProjectIDFiltered(ctx, projectID, envKey, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list features by projectID filtered: %w", err)
	}

	return items, total, nil
}

func (s *Service) ListExtendedByProjectID(
	ctx context.Context,
	projectID domain.ProjectID,
	envKey string,
) ([]domain.FeatureExtended, error) {
	features, err := s.repo.ListByProjectID(ctx, projectID, envKey)
	if err != nil {
		return nil, fmt.Errorf("list features by projectID: %w", err)
	}

	// Resolve environment to ensure child entities are scoped correctly
	env, err := s.environmentsRep.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		return nil, fmt.Errorf("get environment: %w", err)
	}

	result := make([]domain.FeatureExtended, 0, len(features))

	for _, feature := range features {
		variants, err := s.flagVariantsRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return nil, fmt.Errorf("list flag variants: %w", err)
		}

		rules, err := s.rulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return nil, fmt.Errorf("list rules: %w", err)
		}

		schedules, err := s.schedulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return nil, fmt.Errorf("list schedules: %w", err)
		}

		result = append(result, domain.FeatureExtended{
			Feature:      feature,
			FlagVariants: variants,
			Rules:        rules,
			Schedules:    schedules,
		})
	}

	return result, nil
}

func (s *Service) ListExtendedByProjectIDFiltered(
	ctx context.Context,
	projectID domain.ProjectID,
	envKey string,
	filter contract.FeaturesListFilter,
) ([]domain.FeatureExtended, int, error) {
	features, total, err := s.repo.ListByProjectIDFiltered(ctx, projectID, envKey, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list features by projectID: %w", err)
	}

	// Resolve environment to ensure child entities are scoped correctly
	env, err := s.environmentsRep.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		return nil, 0, fmt.Errorf("get environment: %w", err)
	}

	result := make([]domain.FeatureExtended, 0, len(features))

	for _, feature := range features {
		variants, err := s.flagVariantsRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("list flag variants: %w", err)
		}

		rules, err := s.rulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("list rules: %w", err)
		}

		schedules, err := s.schedulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("list schedules: %w", err)
		}

		result = append(result, domain.FeatureExtended{
			Feature:      feature,
			FlagVariants: variants,
			Rules:        rules,
			Schedules:    schedules,
		})
	}

	return result, total, nil
}

func (s *Service) Delete(ctx context.Context, id domain.FeatureID, envKey string) (domain.GuardedResult, error) {
	// Load existing feature to check guard status
	existing, err := s.repo.GetByIDWithEnv(ctx, id, envKey)
	if err != nil {
		return domain.GuardedResult{}, fmt.Errorf("get feature by id: %w", err)
	}

	env, err := s.environmentsRep.GetByProjectIDAndKey(ctx, existing.ProjectID, envKey)
	if err != nil {
		return domain.GuardedResult{}, fmt.Errorf("get env: %w", err)
	}

	// Check if a feature is guarded and create pending change if needed
	guardResult := s.checkFeatureGuardedAndCreatePendingChange(
		ctx,
		id,
		env.ID,
		domain.EntityActionDelete,
		&existing,
		nil, // No new feature for delete
	)

	// If there's a conflict or error, return early
	if guardResult.ChangeConflict || guardResult.Error != nil {
		return guardResult, nil
	}

	// If pending change was created, return it
	if guardResult.Pending {
		return guardResult, nil
	}

	// Feature is not guarded, proceed with normal delete
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.Delete(ctx, env.ID, id); err != nil {
			return fmt.Errorf("delete feature: %w", err)
		}

		return nil
	}); err != nil {
		return domain.GuardedResult{}, fmt.Errorf("tx delete feature: %w", err)
	}

	return domain.GuardedResult{Pending: false}, nil
}

func (s *Service) Toggle(
	ctx context.Context,
	id domain.FeatureID,
	enabled bool,
	envKey string,
) (domain.Feature, domain.GuardedResult, error) {
	// Load existing feature to check guard status
	existing, err := s.repo.GetByIDWithEnv(ctx, id, envKey)
	if err != nil {
		return domain.Feature{}, domain.GuardedResult{}, fmt.Errorf("get feature by id: %w", err)
	}

	env, err := s.environmentsRep.GetByProjectIDAndKey(ctx, existing.ProjectID, envKey)
	if err != nil {
		return domain.Feature{}, domain.GuardedResult{}, fmt.Errorf("get env: %w", err)
	}

	// Create new feature with updated enabled status
	newFeature := existing
	newFeature.Enabled = enabled

	// Check if feature is guarded and create pending change if needed
	guardResult := s.checkFeatureGuardedAndCreatePendingChange(
		ctx,
		id,
		env.ID,
		domain.EntityActionUpdate,
		&existing,
		&newFeature,
	)

	// If there's a conflict or error, return early
	if guardResult.Error != nil {
		return domain.Feature{}, domain.GuardedResult{}, guardResult.Error
	}

	if guardResult.ChangeConflict {
		return domain.Feature{}, guardResult, nil
	}

	// If pending change was created, return it
	if guardResult.Pending {
		return domain.Feature{}, guardResult, nil
	}

	// Feature is not guarded, proceed with normal update of environment-specific params
	var updated domain.Feature

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Update feature_params for this environment (enabled flag)
		params := domain.FeatureParams{
			FeatureID:     existing.ID,
			EnvironmentID: env.ID,
			Enabled:       enabled,
			DefaultValue:  existing.DefaultValue,
			UpdatedAt:     time.Now(),
		}

		if _, err := s.featureParamsRep.Update(ctx, existing.ProjectID, params); err != nil {
			return fmt.Errorf("update feature params: %w", err)
		}

		// Reload feature with environment-specific fields
		reloaded, err := s.repo.GetByIDWithEnv(ctx, id, envKey)
		if err != nil {
			return fmt.Errorf("reload feature after toggle: %w", err)
		}
		updated = reloaded

		return nil
	}); err != nil {
		return domain.Feature{}, domain.GuardedResult{}, fmt.Errorf("tx toggle feature: %w", err)
	}

	return updated, domain.GuardedResult{Pending: false}, nil
}

// UpdateWithChildren updates the feature and reconciles its child entities (variants and rules).
func (s *Service) UpdateWithChildren(
	ctx context.Context,
	envKey string,
	feature domain.Feature,
	variants []domain.FlagVariant,
	rules []domain.Rule,
) (domain.FeatureExtended, domain.GuardedResult, error) {
	var result domain.FeatureExtended

	// Load existing feature to check guard status
	existing, err := s.repo.GetByIDWithEnv(ctx, feature.ID, envKey)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("get feature by id: %w", err)
	}

	env, err := s.environmentsRep.GetByProjectIDAndKey(ctx, existing.ProjectID, envKey)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("get env: %w", err)
	}

	// Check if a feature is guarded and create pending change if needed
	guardResult := s.checkFeatureGuardedAndCreatePendingChange(
		ctx,
		feature.ID,
		env.ID,
		domain.EntityActionUpdate,
		&existing,
		&feature,
	)

	// If there's a conflict or error, return early
	if guardResult.ChangeConflict || guardResult.Error != nil {
		return domain.FeatureExtended{}, guardResult, nil
	}

	// If pending change was created, return it
	if guardResult.Pending {
		return domain.FeatureExtended{}, guardResult, nil
	}

	// Feature is not guarded, proceed with normal update
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		feature.ProjectID = existing.ProjectID

		updated, err := s.repo.Update(ctx, env.ID, feature.BasicFeature)
		if err != nil {
			return fmt.Errorf("update feature: %w", err)
		}

		// Update environment-specific parameters (enabled, default_value)
		params := domain.FeatureParams{
			FeatureID:     feature.ID,
			EnvironmentID: env.ID,
			Enabled:       feature.Enabled,
			DefaultValue:  feature.DefaultValue,
			UpdatedAt:     time.Now(),
		}
		if _, err := s.featureParamsRep.Update(ctx, feature.ProjectID, params); err != nil {
			// Fallback: if params row is missing (should not happen due to trigger), create it
			if errors.Is(err, domain.ErrEntityNotFound) {
				if _, cerr := s.featureParamsRep.Create(ctx, feature.ProjectID, domain.FeatureParams{
					FeatureID:     feature.ID,
					EnvironmentID: env.ID,
					Enabled:       feature.Enabled,
					DefaultValue:  feature.DefaultValue,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}); cerr != nil {
					return fmt.Errorf("create feature params: %w", cerr)
				}
			} else {
				return fmt.Errorf("update feature params: %w", err)
			}
		}

		// Temporarily set updated base fields; will reload with env-specific fields later
		result.Feature = domain.Feature{
			BasicFeature: updated,
			Enabled:      params.Enabled,
			DefaultValue: params.DefaultValue,
		}

		// Reconcile variants (environment-scoped)
		existingVariants, err := s.flagVariantsRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return fmt.Errorf("list flag variants: %w", err)
		}

		// Build lookup maps for existing variants by ID and by Name
		existingVByID := make(map[domain.FlagVariantID]domain.FlagVariant, len(existingVariants))
		existingVByName := make(map[string]domain.FlagVariant, len(existingVariants))
		for _, v := range existingVariants {
			existingVByID[v.ID] = v
			existingVByName[v.Name] = v
		}

		// Map incoming variant IDs to the final IDs that will be used in this environment
		variantIDMap := make(map[domain.FlagVariantID]domain.FlagVariantID)
		// Track which existing IDs should be kept (to avoid deleting updated ones)
		keepVariantIDs := make(map[domain.FlagVariantID]struct{})
		updatedVariants := make([]domain.FlagVariant, 0, len(variants))
		for _, incoming := range variants {
			origID := incoming.ID
			incoming.ProjectID = feature.ProjectID
			incoming.FeatureID = feature.ID
			incoming.EnvironmentID = env.ID

			if incoming.ID != "" {
				if _, ok := existingVByID[incoming.ID]; ok {
					// Update existing by ID in this env
					uv, uErr := s.flagVariantsRep.Update(ctx, incoming)
					if uErr != nil {
						return fmt.Errorf("update flag variant: %w", uErr)
					}
					updatedVariants = append(updatedVariants, uv)
					variantIDMap[origID] = uv.ID
					keepVariantIDs[uv.ID] = struct{}{}

					continue
				}
			}

			// Try match by name within env (IDs may belong to another env)
			if exist, ok := existingVByName[incoming.Name]; ok {
				// Update existing by name; preserve its ID
				savedID := exist.ID
				incoming.ID = exist.ID
				uv, uErr := s.flagVariantsRep.Update(ctx, incoming)
				if uErr != nil {
					return fmt.Errorf("update flag variant (by name): %w", uErr)
				}
				updatedVariants = append(updatedVariants, uv)
				if origID != "" {
					variantIDMap[origID] = savedID
				}
				keepVariantIDs[savedID] = struct{}{}

				continue
			}

			// Create new variant; clear ID to avoid PK conflicts with IDs from other envs
			incoming.ID = ""
			cv, cErr := s.flagVariantsRep.Create(ctx, incoming)
			if cErr != nil {
				return fmt.Errorf("create flag variant: %w", cErr)
			}
			updatedVariants = append(updatedVariants, cv)
			if origID != "" {
				variantIDMap[origID] = cv.ID
			}
			keepVariantIDs[cv.ID] = struct{}{}
		}

		// Delete variants not present in the request after reconciliation
		for id := range existingVByID {
			if _, ok := keepVariantIDs[id]; !ok {
				if dErr := s.flagVariantsRep.Delete(ctx, id); dErr != nil {
					return fmt.Errorf("delete flag variant: %w", dErr)
				}
			}
		}

		result.FlagVariants = updatedVariants

		// Reconcile rules (environment-scoped)
		existingRules, err := s.rulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
		if err != nil {
			return fmt.Errorf("list rules: %w", err)
		}

		existingRMap := make(map[domain.RuleID]domain.Rule, len(existingRules))
		for _, r := range existingRules {
			existingRMap[r.ID] = r
		}

		requestedRMap := make(map[domain.RuleID]domain.Rule, len(rules))
		updatedRules := make([]domain.Rule, 0, len(rules))
		for _, rule := range rules {
			// Remap FlagVariantID if it was changed during variant reconciliation
			if rule.FlagVariantID != nil {
				if newID, ok := variantIDMap[*rule.FlagVariantID]; ok {
					newIDCopy := newID
					rule.FlagVariantID = &newIDCopy
				}
			}
			rule.ProjectID = feature.ProjectID
			rule.FeatureID = feature.ID
			rule.EnvironmentID = env.ID
			if rule.ID != "" {
				requestedRMap[rule.ID] = rule
			}

			if rule.ID != "" {
				if _, ok := existingRMap[rule.ID]; ok {
					ur, uErr := s.rulesRep.Update(ctx, rule)
					if uErr != nil {
						return fmt.Errorf("update rule: %w", uErr)
					}
					updatedRules = append(updatedRules, ur)

					continue
				}
			}

			// Create new rule: clear ID to avoid PK conflicts if client passed an external ID
			rule.ID = ""
			cr, cErr := s.rulesRep.Create(ctx, rule)
			if cErr != nil {
				return fmt.Errorf("create rule: %w", cErr)
			}
			updatedRules = append(updatedRules, cr)
		}

		for id := range existingRMap {
			if _, ok := requestedRMap[id]; !ok {
				if dErr := s.rulesRep.Delete(ctx, id); dErr != nil {
					return fmt.Errorf("delete rule: %w", dErr)
				}
			}
		}

		result.Rules = updatedRules

		// Reload feature with environment-specific fields (enabled, default_value)
		reloaded, rErr := s.repo.GetByIDWithEnv(ctx, feature.ID, envKey)
		if rErr != nil {
			return fmt.Errorf("reload feature after update: %w", rErr)
		}
		result.Feature = reloaded

		return nil
	})
	if err != nil {
		err = fmt.Errorf("tx update feature with children: %w", err)

		return domain.FeatureExtended{}, domain.GuardedResult{}, err
	}

	return result, domain.GuardedResult{Pending: false}, nil
}

// checkFeatureGuardedAndCreatePendingChange checks if a feature is guarded and creates a pending change if needed.
func (s *Service) checkFeatureGuardedAndCreatePendingChange(
	ctx context.Context,
	featureID domain.FeatureID,
	environmentID domain.EnvironmentID,
	action domain.EntityAction,
	oldFeature *domain.Feature,
	newFeature *domain.Feature,
) domain.GuardedResult {
	// Extract user info from context
	requestedBy := appcontext.Username(ctx)
	requestUserID := appcontext.UserID(ctx)
	// Check if a feature is guarded
	isGuarded, err := s.guardService.IsFeatureGuarded(ctx, featureID)
	if err != nil {
		return domain.GuardedResult{
			Error: fmt.Errorf("check feature guarded: %w", err),
		}
	}

	if !isGuarded {
		return domain.GuardedResult{
			Pending: false,
		}
	}

	// Build changes diff for base feature fields
	featureChanges := make(map[string]domain.ChangeValue)
	// Build changes for environment-scoped feature_params
	paramsChanges := make(map[string]domain.ChangeValue)

	if oldFeature != nil && newFeature != nil {
		// Compare fields and build changes
		if oldFeature.Name != newFeature.Name {
			featureChanges["name"] = domain.ChangeValue{
				Old: oldFeature.Name,
				New: newFeature.Name,
			}
		}

		if oldFeature.Description != newFeature.Description {
			featureChanges["description"] = domain.ChangeValue{
				Old: oldFeature.Description,
				New: newFeature.Description,
			}
		}

		if oldFeature.RolloutKey != newFeature.RolloutKey {
			featureChanges["rollout_key"] = domain.ChangeValue{
				Old: oldFeature.RolloutKey.String(),
				New: newFeature.RolloutKey.String(),
			}
		}

		// Environment-scoped fields come from feature_params now
		if oldFeature.Enabled != newFeature.Enabled {
			paramsChanges["enabled"] = domain.ChangeValue{
				Old: oldFeature.Enabled,
				New: newFeature.Enabled,
			}
		}

		if oldFeature.DefaultValue != newFeature.DefaultValue {
			paramsChanges["default_value"] = domain.ChangeValue{
				Old: oldFeature.DefaultValue,
				New: newFeature.DefaultValue,
			}
		}
	}

	entities := make([]domain.EntityChange, 0, 2)

	// Add feature entity only if delete action or there are real changes
	if action == domain.EntityActionDelete || len(featureChanges) > 0 {
		entities = append(entities, domain.EntityChange{
			Entity:   string(domain.EntityFeature),
			EntityID: featureID.String(),
			Action:   action,
			Changes:  featureChanges,
		})
	}

	// Add feature_params entity for env-scoped fields when changed and action is update
	if action == domain.EntityActionUpdate && len(paramsChanges) > 0 {
		entities = append(entities, domain.EntityChange{
			Entity:   string(domain.EntityFeatureParams),
			EntityID: featureID.String(),
			Action:   domain.EntityActionUpdate,
			Changes:  paramsChanges,
		})
	}

	// Create a pending change payload
	payload := domain.PendingChangePayload{
		Entities: entities,
		Meta: domain.PendingChangeMeta{
			Reason: "Feature update via API",
			Client: "ui",
			Origin: "feature-update",
		},
	}

	// Get project ID from feature
	var projectID domain.ProjectID
	if oldFeature != nil {
		projectID = oldFeature.ProjectID
	} else if newFeature != nil {
		projectID = newFeature.ProjectID
	} else {
		// Fallback: get feature to get project ID
		feature, err := s.repo.GetByID(ctx, featureID)
		if err != nil {
			return domain.GuardedResult{
				Error: fmt.Errorf("get feature for project ID: %w", err),
			}
		}

		projectID = feature.ProjectID
	}

	// Convert UserID to *int
	var requestUserIDPtr *int

	if requestUserID != 0 {
		userIDInt := int(requestUserID)
		requestUserIDPtr = &userIDInt
	}

	// Check if this is a single-user project
	activeUserCount, err := s.pendingChangesUseCase.GetProjectActiveUserCount(ctx, projectID)
	if err != nil {
		return domain.GuardedResult{
			Error: fmt.Errorf("get project active user count: %w", err),
		}
	}

	// Check for conflicts before creating pending change
	hasConflict, err := s.pendingChangesUseCase.CheckEntityConflict(ctx, payload.Entities)
	if err != nil {
		return domain.GuardedResult{
			ChangeConflict: false,
			Error:          err,
		}
	}

	if hasConflict {
		return domain.GuardedResult{
			ChangeConflict: true,
		}
	}

	// For single-user projects, always create a pending change but mark it as requiring auto-approve
	// The frontend will handle showing the password/TOTP dialog
	var pendingChange domain.PendingChange
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		pendingChange, err = s.pendingChangesUseCase.Create(
			ctx,
			projectID,
			environmentID,
			requestedBy,
			requestUserIDPtr,
			payload,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.GuardedResult{
			Error: fmt.Errorf("create pending change: %w", err),
		}
	}

	// Add metadata about a single-user project for frontend
	// If exactly 1 active user, treat as a single-user project (enables auto-approve)
	if activeUserCount == 1 {
		// For single-user projects, the frontend should show an auto-approve dialog
		// We'll add this information to the pending change response
		pendingChange.Change.Meta.SingleUserProject = true
	}

	return domain.GuardedResult{
		Pending:       true,
		PendingChange: &pendingChange,
	}
}

func (s *Service) GetFeatureParams(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error) {
	return s.featureParamsRep.ListByFeatureID(ctx, featureID)
}
