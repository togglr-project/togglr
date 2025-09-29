package features

import (
	"context"
	"fmt"
	"reflect"
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
	featureTagsRep        contract.FeatureTagsRepository
	tagsRep               contract.TagsRepository
	environmentsRep       contract.EnvironmentsRepository
	guardService          contract.GuardService
	guardEngine           contract.GuardEngine
	pendingChangesUseCase contract.PendingChangesUseCase
}

func New(
	txManager db.TxManager,
	repo contract.FeaturesRepository,
	flagVariantsRep contract.FlagVariantsRepository,
	rulesRep contract.RulesRepository,
	schedulesRep contract.FeatureSchedulesRepository,
	featureParamsRep contract.FeatureParamsRepository,
	featureTagsRep contract.FeatureTagsRepository,
	tagsRep contract.TagsRepository,
	environmentsRep contract.EnvironmentsRepository,
	guardService contract.GuardService,
	guardEngine contract.GuardEngine,
	pendingChangesUseCase contract.PendingChangesUseCase,
) *Service {
	return &Service{
		txManager:             txManager,
		repo:                  repo,
		flagVariantsRep:       flagVariantsRep,
		rulesRep:              rulesRep,
		schedulesRep:          schedulesRep,
		featureParamsRep:      featureParamsRep,
		featureTagsRep:        featureTagsRep,
		tagsRep:               tagsRep,
		environmentsRep:       environmentsRep,
		guardService:          guardService,
		guardEngine:           guardEngine,
		pendingChangesUseCase: pendingChangesUseCase,
	}
}

func (s *Service) CreateWithChildren(
	ctx context.Context,
	feature domain.Feature,
	variants []domain.FlagVariant,
	rules []domain.Rule,
	tagsIDs []domain.TagID,
) (domain.FeatureExtended, error) {
	var result domain.FeatureExtended

	envs, err := s.environmentsRep.ListByProjectID(ctx, feature.ProjectID)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("get env: %w", err)
	}

	if len(envs) == 0 {
		return domain.FeatureExtended{}, fmt.Errorf("no environments found: %w", domain.ErrEntityNotFound)
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
		result.Enabled = feature.Enabled
		result.DefaultValue = feature.DefaultValue

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

			// Create tags
			for _, tagID := range tagsIDs {
				err = s.featureTagsRep.AddFeatureTag(ctx, createdFeature.ID, tagID)
				if err != nil {
					return fmt.Errorf("add feature tag: %w", err)
				}
			}
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

	// Use new guard engine
	pendingChange, conflict, proceed, err := s.guardEngine.CheckGuardedOperation(
		ctx,
		contract.GuardRequest{
			ProjectID:     existing.ProjectID,
			EnvironmentID: env.ID,
			FeatureID:     id,
			Reason:        "Delete feature via API",
			Origin:        "feature-delete",
			Action:        domain.EntityActionDelete,
			OldEntity:     &existing,
			NewEntity:     nil, // No new entity for delete
		},
	)
	if err != nil {
		return domain.GuardedResult{}, fmt.Errorf("guard check failed: %w", err)
	}

	if conflict {
		return domain.GuardedResult{ChangeConflict: true}, nil
	}

	if pendingChange != nil {
		return domain.GuardedResult{
			Pending:       true,
			PendingChange: pendingChange,
		}, nil
	}

	if !proceed {
		return domain.GuardedResult{}, fmt.Errorf("unexpected guard result: no pending change but proceed=false")
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

	// Use new guard engine
	pendingChange, conflict, proceed, err := s.guardEngine.CheckGuardedOperation(
		ctx,
		contract.GuardRequest{
			ProjectID:     existing.ProjectID,
			EnvironmentID: env.ID,
			FeatureID:     id,
			Reason:        "Toggle feature via API",
			Origin:        "feature-toggle",
			Action:        domain.EntityActionUpdate,
			OldEntity:     &existing,
			NewEntity:     &newFeature,
		},
	)
	if err != nil {
		return domain.Feature{}, domain.GuardedResult{}, fmt.Errorf("guard check failed: %w", err)
	}

	if conflict {
		return domain.Feature{}, domain.GuardedResult{ChangeConflict: true}, nil
	}

	if pendingChange != nil {
		return domain.Feature{}, domain.GuardedResult{
			Pending:       true,
			PendingChange: pendingChange,
		}, nil
	}

	if !proceed {
		return domain.Feature{}, domain.GuardedResult{}, fmt.Errorf("unexpected guard result: no pending change but proceed=false")
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

// computeVariantChanges computes changes for flag variants by comparing existing and new variants.
func (s *Service) computeVariantChanges(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	existingVariants []domain.FlagVariant,
	newVariants []domain.FlagVariant,
) ([]domain.EntityChange, error) {
	// Build maps for comparison
	existingByID := make(map[domain.FlagVariantID]domain.FlagVariant)
	existingByName := make(map[string]domain.FlagVariant)
	for _, v := range existingVariants {
		existingByID[v.ID] = v
		existingByName[v.Name] = v
	}

	newByID := make(map[domain.FlagVariantID]domain.FlagVariant)
	for _, v := range newVariants {
		newByID[v.ID] = v
	}

	var changes []domain.EntityChange

	// Process variants for changes
	for _, newVariant := range newVariants {
		// Set environment ID for new variants
		newVariant.EnvironmentID = envID

		if newVariant.ID != "" {
			// Check if this is an update to existing variant
			if existing, exists := existingByID[newVariant.ID]; exists {
				// Compare and create change if different
				if s.variantsAreDifferent(existing, newVariant) {
					changes = append(changes, domain.EntityChange{
						Entity:   string(domain.EntityFlagVariant),
						EntityID: string(newVariant.ID),
						Action:   domain.EntityActionUpdate,
						Changes:  s.buildVariantChanges(existing, newVariant),
					})
				}
			} else {
				// This is a new variant
				changes = append(changes, domain.EntityChange{
					Entity:   string(domain.EntityFlagVariant),
					EntityID: string(newVariant.ID),
					Action:   domain.EntityActionInsert,
					Changes:  s.buildVariantChanges(domain.FlagVariant{}, newVariant),
				})
			}
		} else {
			// New variant without ID - will be created
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityFlagVariant),
				EntityID: "", // Will be generated
				Action:   domain.EntityActionInsert,
				Changes:  s.buildVariantChanges(domain.FlagVariant{}, newVariant),
			})
		}
	}

	// Check for deleted variants
	for _, existing := range existingVariants {
		if _, exists := newByID[existing.ID]; !exists {
			// This variant was deleted
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityFlagVariant),
				EntityID: string(existing.ID),
				Action:   domain.EntityActionDelete,
				Changes:  nil,
			})
		}
	}

	return changes, nil
}

// computeRuleChanges computes changes for rules by comparing existing and new rules.
func (s *Service) computeRuleChanges(
	ctx context.Context,
	featureID domain.FeatureID,
	envID domain.EnvironmentID,
	existingRules []domain.Rule,
	newRules []domain.Rule,
) ([]domain.EntityChange, error) {
	// Build maps for comparison
	existingByID := make(map[domain.RuleID]domain.Rule)
	for _, r := range existingRules {
		existingByID[r.ID] = r
	}

	newByID := make(map[domain.RuleID]domain.Rule)
	for _, r := range newRules {
		newByID[r.ID] = r
	}

	var changes []domain.EntityChange

	// Process rules for changes
	for _, newRule := range newRules {
		// Set environment ID for new rules
		newRule.EnvironmentID = envID

		if newRule.ID != "" {
			// Check if this is an update to existing rule
			if existing, exists := existingByID[newRule.ID]; exists {
				// Compare and create change if different
				if s.rulesAreDifferent(existing, newRule) {
					changes = append(changes, domain.EntityChange{
						Entity:   string(domain.EntityRule),
						EntityID: string(newRule.ID),
						Action:   domain.EntityActionUpdate,
						Changes:  s.buildRuleChanges(existing, newRule),
					})
				}
			} else {
				// This is a new rule
				changes = append(changes, domain.EntityChange{
					Entity:   string(domain.EntityRule),
					EntityID: string(newRule.ID),
					Action:   domain.EntityActionInsert,
					Changes:  s.buildRuleChanges(domain.Rule{}, newRule),
				})
			}
		} else {
			// New rule without ID - will be created
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityRule),
				EntityID: "", // Will be generated
				Action:   domain.EntityActionInsert,
				Changes:  s.buildRuleChanges(domain.Rule{}, newRule),
			})
		}
	}

	// Check for deleted rules
	for _, existing := range existingRules {
		if _, exists := newByID[existing.ID]; !exists {
			// This rule was deleted
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityRule),
				EntityID: string(existing.ID),
				Action:   domain.EntityActionDelete,
				Changes:  nil,
			})
		}
	}

	return changes, nil
}

// variantsAreDifferent checks if two variants are different.
func (s *Service) variantsAreDifferent(existing, new domain.FlagVariant) bool {
	return existing.Name != new.Name ||
		existing.RolloutPercent != new.RolloutPercent
}

// rulesAreDifferent checks if two rules are different.
func (s *Service) rulesAreDifferent(existing, new domain.Rule) bool {
	return existing.IsCustomized != new.IsCustomized ||
		existing.Action != new.Action ||
		existing.Priority != new.Priority ||
		!s.flagVariantIDsEqual(existing.FlagVariantID, new.FlagVariantID) ||
		!s.segmentIDsEqual(existing.SegmentID, new.SegmentID) ||
		!s.conditionsEqual(existing.Conditions, new.Conditions)
}

// conditionsEqual compares two conditions expressions.
func (s *Service) conditionsEqual(existing, new any) bool {
	// Use reflect.DeepEqual for reliable deep comparison
	return reflect.DeepEqual(existing, new)
}

// computeFeatureChanges computes changes for BasicFeature fields
func (s *Service) computeFeatureChanges(oldFeature, newFeature *domain.Feature) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	// Compare BasicFeature fields
	if oldFeature.Name != newFeature.Name {
		changes["name"] = domain.ChangeValue{
			Old: oldFeature.Name,
			New: newFeature.Name,
		}
	}

	if oldFeature.Description != newFeature.Description {
		changes["description"] = domain.ChangeValue{
			Old: oldFeature.Description,
			New: newFeature.Description,
		}
	}

	if oldFeature.RolloutKey != newFeature.RolloutKey {
		changes["rollout_key"] = domain.ChangeValue{
			Old: oldFeature.RolloutKey,
			New: newFeature.RolloutKey,
		}
	}

	return changes
}

// computeFeatureParamsChanges computes changes for FeatureParams fields
func (s *Service) computeFeatureParamsChanges(oldFeature, newFeature *domain.Feature) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	// Compare FeatureParams fields
	if oldFeature.Enabled != newFeature.Enabled {
		changes["enabled"] = domain.ChangeValue{
			Old: oldFeature.Enabled,
			New: newFeature.Enabled,
		}
	}

	if oldFeature.DefaultValue != newFeature.DefaultValue {
		changes["default_value"] = domain.ChangeValue{
			Old: oldFeature.DefaultValue,
			New: newFeature.DefaultValue,
		}
	}

	return changes
}

// flagVariantIDsEqual compares two flag variant ID pointers.
func (s *Service) flagVariantIDsEqual(a, b *domain.FlagVariantID) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// segmentIDsEqual compares two segment ID pointers.
func (s *Service) segmentIDsEqual(a, b *domain.SegmentID) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// buildVariantChanges builds changes map for a variant.
func (s *Service) buildVariantChanges(existing, new domain.FlagVariant) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if existing.Name != new.Name {
		changes["name"] = domain.ChangeValue{
			Old: existing.Name,
			New: new.Name,
		}
	}

	if existing.RolloutPercent != new.RolloutPercent {
		changes["rollout_percent"] = domain.ChangeValue{
			Old: existing.RolloutPercent,
			New: new.RolloutPercent,
		}
	}

	// For new variants, add required fields
	if existing.ID == "" {
		changes["project_id"] = domain.ChangeValue{New: new.ProjectID}
		changes["feature_id"] = domain.ChangeValue{New: new.FeatureID}
		changes["environment_id"] = domain.ChangeValue{New: new.EnvironmentID}
	}

	return changes
}

// buildRuleChanges builds changes map for a rule.
func (s *Service) buildRuleChanges(existing, new domain.Rule) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if existing.IsCustomized != new.IsCustomized {
		changes["is_customized"] = domain.ChangeValue{
			Old: existing.IsCustomized,
			New: new.IsCustomized,
		}
	}

	if existing.Action != new.Action {
		changes["action"] = domain.ChangeValue{
			Old: existing.Action,
			New: new.Action,
		}
	}

	if existing.Priority != new.Priority {
		changes["priority"] = domain.ChangeValue{
			Old: existing.Priority,
			New: new.Priority,
		}
	}

	if !s.flagVariantIDsEqual(existing.FlagVariantID, new.FlagVariantID) {
		changes["flag_variant_id"] = domain.ChangeValue{
			Old: existing.FlagVariantID,
			New: new.FlagVariantID,
		}
	}

	if !s.segmentIDsEqual(existing.SegmentID, new.SegmentID) {
		changes["segment_id"] = domain.ChangeValue{
			Old: existing.SegmentID,
			New: new.SegmentID,
		}
	}

	if !s.conditionsEqual(existing.Conditions, new.Conditions) {
		changes["condition"] = domain.ChangeValue{
			Old: existing.Conditions,
			New: new.Conditions,
		}
	}

	// For new rules, add required fields
	if existing.ID == "" {
		changes["project_id"] = domain.ChangeValue{New: new.ProjectID}
		changes["feature_id"] = domain.ChangeValue{New: new.FeatureID}
		changes["environment_id"] = domain.ChangeValue{New: new.EnvironmentID}
	}

	return changes
}

// UpdateWithChildren updates the feature and reconciles its child entities (variants and rules).
func (s *Service) UpdateWithChildren(
	ctx context.Context,
	envKey string,
	feature domain.Feature,
	variants []domain.FlagVariant,
	rules []domain.Rule,
	tags []domain.FeatureTags,
) (domain.FeatureExtended, domain.GuardedResult, error) {
	// Load existing feature to check guard status
	existing, err := s.repo.GetByIDWithEnv(ctx, feature.ID, envKey)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("get feature by id: %w", err)
	}

	env, err := s.environmentsRep.GetByProjectIDAndKey(ctx, existing.ProjectID, envKey)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("get env: %w", err)
	}

	// Check if a feature is guarded
	isGuarded, err := s.guardService.IsFeatureGuarded(ctx, feature.ID)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("check feature guarded: %w", err)
	}

	// If a feature is not guarded, proceed with a direct update
	if !isGuarded {
		// Proceed with normal update - use existing logic
		return s.updateFeatureWithChildrenDirect(ctx, feature, env, variants, rules, tags)
	}

	// Load existing variants and rules for comparison
	existingVariants, err := s.flagVariantsRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("list existing variants: %w", err)
	}

	existingRules, err := s.rulesRep.ListByFeatureIDWithEnvID(ctx, feature.ID, env.ID)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("list existing rules: %w", err)
	}

	// Compute changes for variants and rules
	variantChanges, err := s.computeVariantChanges(ctx, feature.ID, env.ID, existingVariants, variants)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("compute variant changes: %w", err)
	}

	ruleChanges, err := s.computeRuleChanges(ctx, feature.ID, env.ID, existingRules, rules)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("compute rule changes: %w", err)
	}

	// Load existing tags and compute tag changes
	existingTags, err := s.featureTagsRep.ListFeatureTags(ctx, feature.ID)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("list existing tags: %w", err)
	}

	// Convert []domain.Tag to []domain.FeatureTags
	existingFeatureTags := make([]domain.FeatureTags, len(existingTags))
	for i, tag := range existingTags {
		existingFeatureTags[i] = domain.FeatureTags{
			FeatureID: feature.ID,
			TagID:     tag.ID,
		}
	}

	tagChanges, err := s.computeTagChanges(ctx, feature.ID, existingFeatureTags, tags)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("compute tag changes: %w", err)
	}

	// Feature is guarded - create pending change
	// Compute feature changes using guard engine logic
	featureChanges := s.computeFeatureChanges(&existing, &feature)
	paramsChanges := s.computeFeatureParamsChanges(&existing, &feature)

	// Create entities list
	entities := make([]domain.EntityChange, 0, 4+len(variantChanges)+len(ruleChanges)+len(tagChanges))

	// Add feature changes if any
	if len(featureChanges) > 0 {
		entities = append(entities, domain.EntityChange{
			Entity:   string(domain.EntityFeature),
			EntityID: string(feature.ID),
			Action:   domain.EntityActionUpdate,
			Changes:  featureChanges,
		})
	}

	// Add feature_params changes if any
	if len(paramsChanges) > 0 {
		entities = append(entities, domain.EntityChange{
			Entity:   string(domain.EntityFeatureParams),
			EntityID: string(feature.ID),
			Action:   domain.EntityActionUpdate,
			Changes:  paramsChanges,
		})
	}

	// Add variant changes
	entities = append(entities, variantChanges...)

	// Add rule changes
	entities = append(entities, ruleChanges...)

	// Add tag changes
	entities = append(entities, tagChanges...)

	// If no changes at all, proceed normally
	if len(entities) == 0 {
		return domain.FeatureExtended{}, domain.GuardedResult{Pending: false}, nil
	}

	// Check if this is a single-user project for auto-approve
	activeUserCount, err := s.pendingChangesUseCase.GetProjectActiveUserCount(ctx, existing.ProjectID)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("count active users: %w", err)
	}

	// Create pending change payload
	payload := domain.PendingChangePayload{
		Entities: entities,
		Meta: domain.PendingChangeMeta{
			Reason:            "Update feature with children via API",
			Client:            "ui",
			Origin:            "feature-update-with-children",
			SingleUserProject: activeUserCount == 1,
		},
	}

	// Get user info from context
	requestedBy := appcontext.Username(ctx)
	requestUserID := appcontext.UserID(ctx)
	var requestUserIDPtr *int
	if requestUserID != 0 {
		userIDInt := int(requestUserID)
		requestUserIDPtr = &userIDInt
	}

	// Check for conflicts before creating pending change
	hasConflict, err := s.pendingChangesUseCase.CheckEntityConflict(ctx, payload.Entities)
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("check entity conflict: %w", err)
	}
	if hasConflict {
		return domain.FeatureExtended{}, domain.GuardedResult{ChangeConflict: true}, nil
	}

	// Create pending change
	var createdPendingChange domain.PendingChange
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		createdPendingChange, err = s.pendingChangesUseCase.Create(
			ctx,
			existing.ProjectID,
			env.ID,
			requestedBy,
			requestUserIDPtr,
			payload,
		)
		return err
	})
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("create pending change: %w", err)
	}

	return domain.FeatureExtended{}, domain.GuardedResult{
		Pending:       true,
		PendingChange: &createdPendingChange,
	}, nil
}

func (s *Service) GetFeatureParams(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error) {
	return s.featureParamsRep.ListByFeatureID(ctx, featureID)
}

// updateFeatureWithChildrenDirect performs direct update without guard engine
func (s *Service) updateFeatureWithChildrenDirect(
	ctx context.Context,
	feature domain.Feature,
	env domain.Environment,
	variants []domain.FlagVariant,
	rules []domain.Rule,
	tags []domain.FeatureTags,
) (domain.FeatureExtended, domain.GuardedResult, error) {
	var result domain.FeatureExtended
	envKey := env.Key

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		feature.ProjectID = feature.ProjectID // Ensure project ID is set

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
			if err == domain.ErrEntityNotFound {
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

		// Reconcile tags
		if err := s.reconcileTags(ctx, feature.ProjectID, feature.ID, tags); err != nil {
			return fmt.Errorf("reconcile tags: %w", err)
		}

		// Reload feature with environment-specific fields (enabled, default_value)
		reloaded, rErr := s.repo.GetByIDWithEnv(ctx, feature.ID, envKey)
		if rErr != nil {
			return fmt.Errorf("reload feature after update: %w", rErr)
		}
		result.Feature = reloaded

		return nil
	})
	if err != nil {
		return domain.FeatureExtended{}, domain.GuardedResult{}, fmt.Errorf("tx update feature with children: %w", err)
	}

	return result, domain.GuardedResult{Pending: false}, nil
}

// computeTagChanges computes changes for tags by comparing existing and new tags.
func (s *Service) computeTagChanges(
	ctx context.Context,
	featureID domain.FeatureID,
	existingTags []domain.FeatureTags,
	newTags []domain.FeatureTags,
) ([]domain.EntityChange, error) {
	// Build maps for comparison
	existingByID := make(map[domain.TagID]domain.FeatureTags)
	for _, t := range existingTags {
		existingByID[t.TagID] = t
	}

	newByID := make(map[domain.TagID]domain.FeatureTags)
	for _, t := range newTags {
		newByID[t.TagID] = t
	}

	var changes []domain.EntityChange

	// Process new tags for changes
	for _, newTag := range newTags {
		if newTag.TagID != "" {
			// Check if this is an update to existing tag
			if existing, exists := existingByID[newTag.TagID]; exists {
				// Compare and create change if different
				if s.tagsAreDifferent(existing, newTag) {
					changes = append(changes, domain.EntityChange{
						Entity:   string(domain.EntityFeatureTag),
						EntityID: string(newTag.TagID),
						Action:   domain.EntityActionUpdate,
						Changes:  s.buildTagChanges(existing, newTag, featureID),
					})
				}
			} else {
				// This is a new tag
				changes = append(changes, domain.EntityChange{
					Entity:   string(domain.EntityFeatureTag),
					EntityID: string(newTag.TagID),
					Action:   domain.EntityActionInsert,
					Changes:  s.buildTagChanges(domain.FeatureTags{}, newTag, featureID),
				})
			}
		} else {
			// New tag without ID - will be created
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityFeatureTag),
				EntityID: "", // Will be generated
				Action:   domain.EntityActionInsert,
				Changes:  s.buildTagChanges(domain.FeatureTags{}, newTag, featureID),
			})
		}
	}

	// Check for deleted tags
	for _, existing := range existingTags {
		if _, exists := newByID[existing.TagID]; !exists {
			// This tag was deleted - need to provide feature_id and tag_id for delete operation
			deleteChanges := map[string]domain.ChangeValue{
				"feature_id": {New: featureID},
				"tag_id":     {New: existing.TagID},
			}
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityFeatureTag),
				EntityID: string(existing.TagID),
				Action:   domain.EntityActionDelete,
				Changes:  deleteChanges,
			})
		}
	}

	return changes, nil
}

// tagsAreDifferent checks if two tags are different.
func (s *Service) tagsAreDifferent(existing, new domain.FeatureTags) bool {
	return existing.TagID != new.TagID
}

// buildTagChanges builds changes map for a tag.
func (s *Service) buildTagChanges(existing, new domain.FeatureTags, featureID domain.FeatureID) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	// Always add feature_id for both insert and update operations
	changes["feature_id"] = domain.ChangeValue{New: featureID}

	if existing.TagID != new.TagID {
		changes["tag_id"] = domain.ChangeValue{
			Old: existing.TagID,
			New: new.TagID,
		}
	}

	// For new tags, add tag_id
	if existing.TagID == "" {
		changes["tag_id"] = domain.ChangeValue{New: new.TagID}
	}

	return changes
}

// reconcileTags reconciles feature tags by comparing existing and new tags.
func (s *Service) reconcileTags(ctx context.Context, projectID domain.ProjectID, featureID domain.FeatureID, newTags []domain.FeatureTags) error {
	// Get existing tags
	existingTags, err := s.featureTagsRep.ListFeatureTags(ctx, featureID)
	if err != nil {
		return fmt.Errorf("list existing tags: %w", err)
	}

	// Build maps for comparison
	existingByID := make(map[domain.TagID]bool)
	for _, tag := range existingTags {
		existingByID[tag.ID] = true
	}

	newByID := make(map[domain.TagID]bool)
	for _, tag := range newTags {
		newByID[tag.TagID] = true
	}

	// Remove tags that are no longer present
	for _, existingTag := range existingTags {
		if !newByID[existingTag.ID] {
			if err := s.featureTagsRep.RemoveFeatureTag(ctx, featureID, existingTag.ID); err != nil {
				return fmt.Errorf("remove feature tag: %w", err)
			}
		}
	}

	// Add new tags
	for _, newTag := range newTags {
		if !existingByID[newTag.TagID] {
			if err := s.featureTagsRep.AddFeatureTag(ctx, featureID, newTag.TagID); err != nil {
				return fmt.Errorf("add feature tag: %w", err)
			}
		}
	}

	return nil
}
