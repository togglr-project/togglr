package contract

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type FeatureTagsUseCase interface {
	ListFeatureTags(ctx context.Context, featureID domain.FeatureID) ([]domain.Tag, error)
	AddFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error
	RemoveFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error
}

type FeatureTagsRepository interface {
	ListFeatureTags(ctx context.Context, featureID domain.FeatureID) ([]domain.Tag, error)
	AddFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error
	RemoveFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error
	HasFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) (bool, error)
}
