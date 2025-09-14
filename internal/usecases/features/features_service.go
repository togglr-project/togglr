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
}

func New(
	txManager db.TxManager,
	repo contract.FeaturesRepository,
	flagVariantsRep contract.FlagVariantsRepository,
	rulesRep contract.RulesRepository,
) *Service {
	return &Service{
		txManager:       txManager,
		repo:            repo,
		flagVariantsRep: flagVariantsRep,
		rulesRep:        rulesRep,
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
		for _, v := range variants {
			v.FeatureID = createdFeature.ID
			cv, err := s.flagVariantsRep.Create(ctx, v)
			if err != nil {
				return fmt.Errorf("create flag variant: %w", err)
			}
			createdVariants = append(createdVariants, cv)
		}
		result.FlagVariants = createdVariants

		// Create rules
		createdRules := make([]domain.Rule, 0, len(rules))
		for _, r := range rules {
			r.FeatureID = createdFeature.ID
			cr, err := s.rulesRep.Create(ctx, r)
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
