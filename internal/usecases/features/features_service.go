//nolint:nestif,maintidx // fix it
package features

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
	simplecache "github.com/togglr-project/togglr/pkg/simple-cache"
)

const CacheTTL = 10 * time.Minute

type Service struct {
	txManager                db.TxManager
	repo                     contract.FeaturesRepository
	flagVariantsRep          contract.FlagVariantsRepository
	rulesRep                 contract.RulesRepository
	schedulesRep             contract.FeatureSchedulesRepository
	featureParamsRep         contract.FeatureParamsRepository
	featureTagsRep           contract.FeatureTagsRepository
	tagsRep                  contract.TagsRepository
	environmentsRep          contract.EnvironmentsRepository
	featureNotificationsRepo contract.FeatureNotificationRepository
	guardService             contract.GuardService
	guardEngine              contract.GuardEngine
	pendingChangesUseCase    contract.PendingChangesUseCase
	cache                    *simplecache.Cache[string, domain.Feature]
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
	featureNotificationsRepo contract.FeatureNotificationRepository,
	guardService contract.GuardService,
	guardEngine contract.GuardEngine,
	pendingChangesUseCase contract.PendingChangesUseCase,
) *Service {
	cache := simplecache.New[string, domain.Feature]()
	cache.StartCleanup(time.Minute)

	return &Service{
		txManager:                txManager,
		repo:                     repo,
		flagVariantsRep:          flagVariantsRep,
		rulesRep:                 rulesRep,
		schedulesRep:             schedulesRep,
		featureParamsRep:         featureParamsRep,
		featureTagsRep:           featureTagsRep,
		tagsRep:                  tagsRep,
		environmentsRep:          environmentsRep,
		featureNotificationsRepo: featureNotificationsRepo,
		guardService:             guardService,
		guardEngine:              guardEngine,
		pendingChangesUseCase:    pendingChangesUseCase,
		cache:                    cache,
	}
}

//nolint:gocognit // This is a complex function.
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

func (s *Service) GetByKeyWithEnvCached(ctx context.Context, key, envKey string) (domain.Feature, error) {
	cacheKey := makeFeatureCacheKey(key, envKey)

	if cached, found := s.cache.Get(cacheKey); found {
		return cached, nil
	}

	feature, err := s.repo.GetByKeyWithEnv(ctx, key, envKey)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("get feature by key with environment: %w", err)
	}

	s.cache.Set(cacheKey, feature, CacheTTL)

	return feature, nil
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

	// Check if feature is guarded
	isGuarded, err := s.guardService.IsFeatureGuarded(ctx, id)
	if err != nil {
		return domain.GuardedResult{}, fmt.Errorf("check feature guarded: %w", err)
	}

	// If feature is not guarded, proceed with direct delete
	if !isGuarded {
		if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			if err := s.repo.Delete(ctx, env.ID, id); err != nil {
				return fmt.Errorf("delete feature: %w", err)
			}

			return nil
		}); err != nil {
			return domain.GuardedResult{}, fmt.Errorf("tx delete feature: %w", err)
		}

		return domain.GuardedResult{}, nil
	}

	// Feature is guarded - create pending change
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
		return domain.GuardedResult{}, errors.New("unexpected guard result: no pending change but proceed=false")
	}

	// This should never be reached for guarded features
	return domain.GuardedResult{}, errors.New("unexpected: guarded feature but no pending change created")
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

	// Check if feature is guarded
	isGuarded, err := s.guardService.IsFeatureGuarded(ctx, id)
	if err != nil {
		return domain.Feature{}, domain.GuardedResult{}, fmt.Errorf("check feature guarded: %w", err)
	}

	// If a feature is not guarded, proceed with a direct update
	if !isGuarded {
		// Proceed with normal update of environment-specific params
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

			// Update the feature object to return
			updated = existing
			updated.Enabled = enabled

			errNotif := s.featureNotificationsRepo.AddNotification(
				ctx,
				existing.ProjectID,
				env.ID,
				existing.ID,
				makeStateNotificationPayload(ctx, enabled),
			)
			if errNotif != nil {
				slog.Error("failed to add notification", "error", errNotif)
			}

			return nil
		}); err != nil {
			return domain.Feature{}, domain.GuardedResult{}, fmt.Errorf("tx update feature: %w", err)
		}

		return updated, domain.GuardedResult{}, nil
	}

	// Feature is guarded - create pending change
	// Create new feature with updated enabled status
	newFeature := existing
	newFeature.Enabled = enabled

	// Use the new guard engine
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
		err = errors.New("unexpected guard result: no pending change but proceed=false")

		return domain.Feature{}, domain.GuardedResult{}, err
	}

	// This should never be reached for guarded features
	err = errors.New("unexpected: guarded feature but no pending change created")

	return domain.Feature{}, domain.GuardedResult{}, err
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
		return s.updateFeatureWithChildrenDirect(ctx, feature, existing, env, variants, rules, tags)
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
	variantChanges := s.computeVariantChanges(env.ID, existingVariants, variants)
	ruleChanges := s.computeRuleChanges(env.ID, existingRules, rules)

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

	tagChanges := s.computeTagChanges(feature.ID, existingFeatureTags, tags)

	// Feature is guarded - create pending change
	// Compute feature changes using guard engine logic
	featureChanges := s.computeFeatureChanges(&existing, &feature)
	paramsChanges := s.computeFeatureParamsChanges(&existing, &feature)

	// Create an entities list
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

// computeVariantChanges computes changes for flag variants by comparing existing and new variants.
func (s *Service) computeVariantChanges(
	envID domain.EnvironmentID,
	existingVariants []domain.FlagVariant,
	newVariants []domain.FlagVariant,
) []domain.EntityChange {
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
			// Check if this is an update to an existing variant
			if existing, exists := existingByID[newVariant.ID]; exists {
				// Use the guard engine to compute changes
				variantChanges := s.guardEngine.BuildChangeDiff(&existing, &newVariant)
				// Only create change if there are differences
				if len(variantChanges) > 0 {
					changes = append(changes, domain.EntityChange{
						Entity:   string(domain.EntityFlagVariant),
						EntityID: string(newVariant.ID),
						Action:   domain.EntityActionUpdate,
						Changes:  variantChanges,
					})
				}
			} else {
				// This is a new variant - use guard engine for insert changes
				variantChanges := s.guardEngine.BuildInsertChanges(&newVariant)
				changes = append(changes, domain.EntityChange{
					Entity:   string(domain.EntityFlagVariant),
					EntityID: string(newVariant.ID),
					Action:   domain.EntityActionInsert,
					Changes:  variantChanges,
				})
			}
		} else {
			// New variant without ID - will be created
			variantChanges := s.guardEngine.BuildInsertChanges(&newVariant)
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityFlagVariant),
				EntityID: "", // Will be generated
				Action:   domain.EntityActionInsert,
				Changes:  variantChanges,
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

	return changes
}

// computeRuleChanges computes changes for rules by comparing existing and new rules.
func (s *Service) computeRuleChanges(
	envID domain.EnvironmentID,
	existingRules []domain.Rule,
	newRules []domain.Rule,
) []domain.EntityChange {
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
				// Use guard engine to compute changes
				ruleChanges := s.guardEngine.BuildChangeDiff(&existing, &newRule)
				// Only create change if there are differences
				if len(ruleChanges) > 0 {
					changes = append(changes, domain.EntityChange{
						Entity:   string(domain.EntityRule),
						EntityID: string(newRule.ID),
						Action:   domain.EntityActionUpdate,
						Changes:  ruleChanges,
					})
				}
			} else {
				// This is a new rule - use guard engine for insert changes
				ruleChanges := s.guardEngine.BuildInsertChanges(&newRule)
				changes = append(changes, domain.EntityChange{
					Entity:   string(domain.EntityRule),
					EntityID: string(newRule.ID),
					Action:   domain.EntityActionInsert,
					Changes:  ruleChanges,
				})
			}
		} else {
			// New rule without ID - will be created
			ruleChanges := s.guardEngine.BuildInsertChanges(&newRule)
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityRule),
				EntityID: "", // Will be generated
				Action:   domain.EntityActionInsert,
				Changes:  ruleChanges,
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

	return changes
}

// computeFeatureChanges computes changes for BasicFeature fields using the guard engine.
func (s *Service) computeFeatureChanges(oldFeature, newFeature *domain.Feature) map[string]domain.ChangeValue {
	// Extract BasicFeature from both features
	oldBasic := oldFeature.BasicFeature
	newBasic := newFeature.BasicFeature

	// Use guard engine to compute changes for BasicFeature fields
	return s.guardEngine.BuildChangeDiff(&oldBasic, &newBasic)
}

// computeFeatureParamsChanges computes changes for FeatureParams fields using the guard engine.
func (s *Service) computeFeatureParamsChanges(oldFeature, newFeature *domain.Feature) map[string]domain.ChangeValue {
	// Convert features to FeatureParams using the ConvertToFeatureParams method
	oldParams := oldFeature.ConvertToFeatureParams()
	newParams := newFeature.ConvertToFeatureParams()

	// Use guard engine to compute changes for FeatureParams fields
	return s.guardEngine.BuildChangeDiff(&oldParams, &newParams)
}

// updateFeatureWithChildrenDirect performs direct update without a guard engine.
//
//nolint:gocognit // This is a complex function
func (s *Service) updateFeatureWithChildrenDirect(
	ctx context.Context,
	feature domain.Feature,
	existing domain.Feature,
	env domain.Environment,
	variants []domain.FlagVariant,
	rules []domain.Rule,
	tags []domain.FeatureTags,
) (domain.FeatureExtended, domain.GuardedResult, error) {
	var result domain.FeatureExtended
	envKey := env.Key

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Ensure project ID is set

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

		if feature.Enabled != existing.Enabled {
			errNotif := s.featureNotificationsRepo.AddNotification(
				ctx,
				feature.ProjectID,
				feature.EnvironmentID,
				feature.ID,
				makeStateNotificationPayload(ctx, feature.Enabled),
			)
			if errNotif != nil {
				slog.Error("failed to add feature notification", "error", errNotif)
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
		if err := s.reconcileTags(ctx, feature.ID, tags); err != nil {
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
		err = fmt.Errorf("tx update feature with children: %w", err)

		return domain.FeatureExtended{}, domain.GuardedResult{}, err
	}

	return result, domain.GuardedResult{Pending: false}, nil
}

// computeTagChanges computes changes for tags by comparing existing and new tags.
func (s *Service) computeTagChanges(
	featureID domain.FeatureID,
	existingTags []domain.FeatureTags,
	newTags []domain.FeatureTags,
) []domain.EntityChange {
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
		// Set feature_id for new tags
		newTag.FeatureID = featureID

		if newTag.TagID != "" {
			// Check if this is an update to the existing tag
			if existing, exists := existingByID[newTag.TagID]; exists {
				// Use the guard engine to compute changes
				tagChanges := s.guardEngine.BuildChangeDiff(&existing, &newTag)
				// Only create change if there are differences
				if len(tagChanges) > 0 {
					changes = append(changes, domain.EntityChange{
						Entity:   string(domain.EntityFeatureTag),
						EntityID: string(newTag.TagID),
						Action:   domain.EntityActionUpdate,
						Changes:  tagChanges,
					})
				}
			} else {
				// This is a new tag - use guard engine for insert changes
				tagChanges := s.guardEngine.BuildInsertChanges(&newTag)
				changes = append(changes, domain.EntityChange{
					Entity:   string(domain.EntityFeatureTag),
					EntityID: string(newTag.TagID),
					Action:   domain.EntityActionInsert,
					Changes:  tagChanges,
				})
			}
		} else {
			// New tag without ID - will be created
			tagChanges := s.guardEngine.BuildInsertChanges(&newTag)
			changes = append(changes, domain.EntityChange{
				Entity:   string(domain.EntityFeatureTag),
				EntityID: "", // Will be generated
				Action:   domain.EntityActionInsert,
				Changes:  tagChanges,
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

	return changes
}

// reconcileTags reconciles feature tags by comparing existing and new tags.
func (s *Service) reconcileTags(ctx context.Context, featureID domain.FeatureID, newTags []domain.FeatureTags) error {
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

func (s *Service) InvalidateCache(key, envKey string) {
	cacheKey := makeFeatureCacheKey(key, envKey)
	s.cache.Delete(cacheKey)
}

func (s *Service) InvalidateProjectCache(projectID domain.ProjectID) {
	// For now, we'll clear all cache entries since we don't have project-specific keys
	// In a more sophisticated implementation, we could maintain project-to-keys mapping
	s.cache.Clear()
}

func makeFeatureCacheKey(key, envKey string) string {
	return key + ":" + envKey
}

func makeStateNotificationPayload(ctx context.Context, enabled bool) domain.FeatureNotificationPayload {
	return domain.FeatureNotificationPayload{
		State: &domain.FeatureNotificationStatePayload{
			Enabled:   enabled,
			ChangedBy: appcontext.Username(ctx),
		},
	}
}
