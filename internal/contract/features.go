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
	// UpdateWithChildren updates feature and reconciles its variants and rules in a single transaction.
	UpdateWithChildren(
		ctx context.Context,
		feature domain.Feature,
		variants []domain.FlagVariant,
		rules []domain.Rule,
	) (domain.FeatureExtended, error)
	GetByID(ctx context.Context, id domain.FeatureID) (domain.Feature, error)
	GetExtendedByID(ctx context.Context, id domain.FeatureID) (domain.FeatureExtended, error)
	GetByKey(ctx context.Context, key string) (domain.Feature, error)
	List(ctx context.Context) ([]domain.Feature, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Feature, error)
	ListExtendedByProjectID(
		ctx context.Context,
		projectID domain.ProjectID,
	) ([]domain.FeatureExtended, error)
	Update(ctx context.Context, feature domain.Feature) (domain.Feature, error)
	// Toggle enables or disables a feature flag by its ID and returns updated entity.
	Toggle(ctx context.Context, id domain.FeatureID, enabled bool) (domain.Feature, error)
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
