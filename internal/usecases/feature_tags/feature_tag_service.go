package feature_tags

import (
	"context"
	"fmt"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
)

type Service struct {
	featureTagRepo contract.FeatureTagsRepository
	tagRepo        contract.TagsRepository
	featureRepo    contract.FeaturesRepository
}

func New(
	featureTagRepo contract.FeatureTagsRepository,
	tagRepo contract.TagsRepository,
	featureRepo contract.FeaturesRepository,
) *Service {
	return &Service{
		featureTagRepo: featureTagRepo,
		tagRepo:        tagRepo,
		featureRepo:    featureRepo,
	}
}

func (s *Service) ListFeatureTags(ctx context.Context, featureID domain.FeatureID) ([]domain.Tag, error) {
	// Check if feature exists
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
	// Check if feature exists
	_, err := s.featureRepo.GetByID(ctx, featureID)
	if err != nil {
		return fmt.Errorf("get feature: %w", err)
	}

	// Check if tag exists
	_, err = s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

	// Check if association already exists
	exists, err := s.featureTagRepo.HasFeatureTag(ctx, featureID, tagID)
	if err != nil {
		return fmt.Errorf("check feature tag: %w", err)
	}
	if exists {
		return fmt.Errorf("tag already associated with feature")
	}

	// Add association
	err = s.featureTagRepo.AddFeatureTag(ctx, featureID, tagID)
	if err != nil {
		return fmt.Errorf("add feature tag: %w", err)
	}

	return nil
}

func (s *Service) RemoveFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error {
	// Check if feature exists
	_, err := s.featureRepo.GetByID(ctx, featureID)
	if err != nil {
		return fmt.Errorf("get feature: %w", err)
	}

	// Check if tag exists
	_, err = s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

	// Remove association
	err = s.featureTagRepo.RemoveFeatureTag(ctx, featureID, tagID)
	if err != nil {
		return fmt.Errorf("remove feature tag: %w", err)
	}

	return nil
}
