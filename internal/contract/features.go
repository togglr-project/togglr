package contract

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type FeaturesUseCase interface {
	Create(ctx context.Context, feature domain.Feature) (domain.Feature, error)
	// CreateWithChildren creates feature and its related variants and rules in a single transaction.
	CreateWithChildren(
		ctx context.Context,
		feature domain.Feature,
		variants []domain.FlagVariant,
		rules []domain.Rule,
	) (domain.FeatureExtended, error)
	GetByID(ctx context.Context, id domain.FeatureID) (domain.Feature, error)
	GetByKey(ctx context.Context, key string) (domain.Feature, error)
	List(ctx context.Context) ([]domain.Feature, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Feature, error)
	Update(ctx context.Context, feature domain.Feature) (domain.Feature, error)
	Delete(ctx context.Context, id domain.FeatureID) error
}

type FeaturesRepository interface {
	Create(ctx context.Context, feature domain.Feature) (domain.Feature, error)
	GetByID(ctx context.Context, id domain.FeatureID) (domain.Feature, error)
	GetByKey(ctx context.Context, key string) (domain.Feature, error)
	List(ctx context.Context) ([]domain.Feature, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Feature, error)
	Update(ctx context.Context, feature domain.Feature) (domain.Feature, error)
	Delete(ctx context.Context, id domain.FeatureID) error
}
