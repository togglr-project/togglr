package features

import (
	"context"
	"fmt"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

type Service struct {
	txManager       db.TxManager
	repo            contract.FeaturesRepository
	flagVariantsRep contract.FlagVariantsRepository
	rulesRep        contract.RulesRepository
	schedulesRep    contract.FeatureSchedulesRepository
}

func New(
	txManager db.TxManager,
	repo contract.FeaturesRepository,
	flagVariantsRep contract.FlagVariantsRepository,
	rulesRep contract.RulesRepository,
	schedulesRep contract.FeatureSchedulesRepository,
) *Service {
	return &Service{
		txManager:       txManager,
		repo:            repo,
		flagVariantsRep: flagVariantsRep,
		rulesRep:        rulesRep,
		schedulesRep:    schedulesRep,
	}
}

func (s *Service) Create(ctx context.Context, feature domain.Feature) (domain.Feature, error) {
	var created domain.Feature
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.repo.Create(ctx, feature)
		if err != nil {
			return fmt.Errorf("create feature: %w", err)
		}
		return nil
	}); err != nil {
		return domain.Feature{}, fmt.Errorf("tx create feature: %w", err)
	}
	return created, nil
}

func (s *Service) CreateWithChildren(
	ctx context.Context,
	feature domain.Feature,
	variants []domain.FlagVariant,
	rules []domain.Rule,
) (domain.FeatureExtended, error) {
	var result domain.FeatureExtended

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Create feature first
		createdFeature, err := s.repo.Create(ctx, feature)
		if err != nil {
			return fmt.Errorf("create feature: %w", err)
		}
		result.Feature = createdFeature

		// Create variants
		createdVariants := make([]domain.FlagVariant, 0, len(variants))
		for _, variant := range variants {
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
		createdRules := make([]domain.Rule, 0, len(rules))
		for _, rule := range rules {
			rule.FeatureID = createdFeature.ID
			rule.ProjectID = createdFeature.ProjectID
			cr, err := s.rulesRep.Create(ctx, rule)
			if err != nil {
				return fmt.Errorf("create rule: %w", err)
			}
			createdRules = append(createdRules, cr)
		}
		result.Rules = createdRules

		return nil
	}); err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("tx create feature with children: %w", err)
	}

	return result, nil
}

func (s *Service) GetByID(ctx context.Context, id domain.FeatureID) (domain.Feature, error) {
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("get feature by id: %w", err)
	}
	return f, nil
}

func (s *Service) GetExtendedByID(ctx context.Context, id domain.FeatureID) (domain.FeatureExtended, error) {
	feature, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("get feature by id: %w", err)
	}

	variants, err := s.flagVariantsRep.ListByFeatureID(ctx, feature.ID)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("list flag variants: %w", err)
	}

	rules, err := s.rulesRep.ListByFeatureID(ctx, feature.ID)
	if err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("list rules: %w", err)
	}

	schedules, err := s.schedulesRep.ListByFeatureID(ctx, feature.ID)
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

func (s *Service) GetByKey(ctx context.Context, key string) (domain.Feature, error) {
	f, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return domain.Feature{}, fmt.Errorf("get feature by key: %w", err)
	}
	return f, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Feature, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list features: %w", err)
	}
	return items, nil
}

func (s *Service) ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Feature, error) {
	items, err := s.repo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list features by projectID: %w", err)
	}
	return items, nil
}

func (s *Service) ListExtendedByProjectID(
	ctx context.Context,
	projectID domain.ProjectID,
) ([]domain.FeatureExtended, error) {
	features, err := s.repo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list features by projectID: %w", err)
	}

	result := make([]domain.FeatureExtended, 0, len(features))
	for _, feature := range features {
		variants, err := s.flagVariantsRep.ListByFeatureID(ctx, feature.ID)
		if err != nil {
			return nil, fmt.Errorf("list flag variants: %w", err)
		}

		rules, err := s.rulesRep.ListByFeatureID(ctx, feature.ID)
		if err != nil {
			return nil, fmt.Errorf("list rules: %w", err)
		}

		schedules, err := s.schedulesRep.ListByFeatureID(ctx, feature.ID)
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

func (s *Service) Update(ctx context.Context, feature domain.Feature) (domain.Feature, error) {
	var updated domain.Feature
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.Update(ctx, feature)
		if err != nil {
			return fmt.Errorf("update feature: %w", err)
		}
		return nil
	}); err != nil {
		return domain.Feature{}, fmt.Errorf("tx update feature: %w", err)
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id domain.FeatureID) error {
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("delete feature: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("tx delete feature: %w", err)
	}
	return nil
}

func (s *Service) Toggle(ctx context.Context, id domain.FeatureID, enabled bool) (domain.Feature, error) {
	var updated domain.Feature
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Load existing feature to ensure it exists and to keep other fields unchanged
		existing, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("get feature by id: %w", err)
		}

		existing.Enabled = enabled

		updated, err = s.repo.Update(ctx, existing)
		if err != nil {
			return fmt.Errorf("update feature: %w", err)
		}
		return nil
	}); err != nil {
		return domain.Feature{}, fmt.Errorf("tx toggle feature: %w", err)
	}
	return updated, nil
}

// UpdateWithChildren updates feature and reconciles its child entities (variants and rules).
func (s *Service) UpdateWithChildren(
	ctx context.Context,
	feature domain.Feature,
	variants []domain.FlagVariant,
	rules []domain.Rule,
) (domain.FeatureExtended, error) {
	var result domain.FeatureExtended

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Load to preserve immutable fields like ProjectID
		existing, err := s.repo.GetByID(ctx, feature.ID)
		if err != nil {
			return fmt.Errorf("get feature by id: %w", err)
		}

		feature.ProjectID = existing.ProjectID

		updated, err := s.repo.Update(ctx, feature)
		if err != nil {
			return fmt.Errorf("update feature: %w", err)
		}
		result.Feature = updated

		// Reconcile variants
		existingVariants, err := s.flagVariantsRep.ListByFeatureID(ctx, feature.ID)
		if err != nil {
			return fmt.Errorf("list flag variants: %w", err)
		}

		existingVMap := make(map[domain.FlagVariantID]domain.FlagVariant, len(existingVariants))
		for _, v := range existingVariants {
			existingVMap[v.ID] = v
		}

		requestedVMap := make(map[domain.FlagVariantID]domain.FlagVariant, len(variants))
		updatedVariants := make([]domain.FlagVariant, 0, len(variants))
		for _, variant := range variants {
			variant.ProjectID = feature.ProjectID
			variant.FeatureID = feature.ID
			if variant.ID != "" {
				requestedVMap[variant.ID] = variant
			}

			if variant.ID != "" {
				if _, ok := existingVMap[variant.ID]; ok {
					uv, uErr := s.flagVariantsRep.Update(ctx, variant)
					if uErr != nil {
						return fmt.Errorf("update flag variant: %w", uErr)
					}
					updatedVariants = append(updatedVariants, uv)
					continue
				}
			}

			cv, cErr := s.flagVariantsRep.Create(ctx, variant)
			if cErr != nil {
				return fmt.Errorf("create flag variant: %w", cErr)
			}
			updatedVariants = append(updatedVariants, cv)
		}

		// Delete variants not present in request
		for id := range existingVMap {
			if _, ok := requestedVMap[id]; !ok {
				if dErr := s.flagVariantsRep.Delete(ctx, id); dErr != nil {
					return fmt.Errorf("delete flag variant: %w", dErr)
				}
			}
		}

		result.FlagVariants = updatedVariants

		// Reconcile rules
		existingRules, err := s.rulesRep.ListByFeatureID(ctx, feature.ID)
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
			rule.ProjectID = feature.ProjectID
			rule.FeatureID = feature.ID
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

		return nil
	}); err != nil {
		return domain.FeatureExtended{}, fmt.Errorf("tx update feature with children: %w", err)
	}

	return result, nil
}
