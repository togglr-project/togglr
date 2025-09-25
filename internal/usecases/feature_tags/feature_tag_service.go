package feature_tags

import (
	"context"
	"errors"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Service struct {
	txManager      db.TxManager
	featureTagRepo contract.FeatureTagsRepository
	tagRepo        contract.TagsRepository
	featureRepo    contract.FeaturesRepository
}

func New(
	txManager db.TxManager,
	featureTagRepo contract.FeatureTagsRepository,
	tagRepo contract.TagsRepository,
	featureRepo contract.FeaturesRepository,
) *Service {
	return &Service{
		txManager:      txManager,
		featureTagRepo: featureTagRepo,
		tagRepo:        tagRepo,
		featureRepo:    featureRepo,
	}
}

func (s *Service) ListFeatureTags(ctx context.Context, featureID domain.FeatureID) ([]domain.Tag, error) {
	// Check if the feature exists
	_, err := s.featureRepo.GetByID(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("get feature: %w", err)
	}

	// Get feature tags
	tags, err := s.featureTagRepo.ListFeatureTags(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("list feature tags: %w", err)
	}

	return tags, nil
}

func (s *Service) AddFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error {
	// Check if the feature exists
	_, err := s.featureRepo.GetByID(ctx, featureID)
	if err != nil {
		return fmt.Errorf("get feature: %w", err)
	}

	// Check if tag exists
	_, err = s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

	// Check if the association already exists
	exists, err := s.featureTagRepo.HasFeatureTag(ctx, featureID, tagID)
	if err != nil {
		return fmt.Errorf("check feature tag: %w", err)
	}

	if exists {
		return errors.New("tag already associated with feature")
	}

	// Add association
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.featureTagRepo.AddFeatureTag(ctx, featureID, tagID)
	})
	if err != nil {
		return fmt.Errorf("add feature tag: %w", err)
	}

	return nil
}

func (s *Service) RemoveFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error {
	// Check if the feature exists
	_, err := s.featureRepo.GetByID(ctx, featureID)
	if err != nil {
		return fmt.Errorf("get feature: %w", err)
	}

	// Check if a tag exists
	_, err = s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

	// Remove association
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.featureTagRepo.RemoveFeatureTag(ctx, featureID, tagID)
	})
	if err != nil {
		return fmt.Errorf("remove feature tag: %w", err)
	}

	return nil
}
