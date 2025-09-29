package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeaturesListFilter struct {
	Kind         *domain.FeatureKind
	Enabled      *bool
	TextSelector *string
	TagIDs       []domain.TagID
	SortBy       string // name, key, kind, enabled, created_at, updated_at
	SortDesc     bool
	Page         uint
	PerPage      uint
}

type FeaturesUseCase interface {
	// CreateWithChildren creates feature and its related variants and rules in a single transaction.
	CreateWithChildren(
		ctx context.Context,
		feature domain.Feature,
		variants []domain.FlagVariant,
		rules []domain.Rule,
		tagsIDs []domain.TagID,
	) (domain.FeatureExtended, error)
	// UpdateWithChildren updates feature and reconciles its variants and rules in a single transaction.
	UpdateWithChildren(
		ctx context.Context,
		envKey string,
		feature domain.Feature,
		variants []domain.FlagVariant,
		rules []domain.Rule,
		tags []domain.FeatureTags,
	) (domain.FeatureExtended, domain.GuardedResult, error)
	GetByIDWithEnv(ctx context.Context, id domain.FeatureID, env string) (domain.Feature, error)
	GetExtendedByID(ctx context.Context, id domain.FeatureID, env string) (domain.FeatureExtended, error)
	GetByKeyWithEnv(ctx context.Context, key, env string) (domain.Feature, error)
	List(ctx context.Context, env string) ([]domain.Feature, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID, env string) ([]domain.Feature, error)
	ListByProjectIDFiltered(
		ctx context.Context,
		projectID domain.ProjectID,
		env string,
		filter FeaturesListFilter,
	) ([]domain.Feature, int, error)
	ListExtendedByProjectIDFiltered(
		ctx context.Context,
		projectID domain.ProjectID,
		env string,
		filter FeaturesListFilter,
	) ([]domain.FeatureExtended, int, error)
	ListExtendedByProjectID(
		ctx context.Context,
		projectID domain.ProjectID,
		env string,
	) ([]domain.FeatureExtended, error)
	// Toggle enables or disables a feature flag by its ID and returns updated entity.
	Toggle(ctx context.Context, id domain.FeatureID, enabled bool, env string) (domain.Feature, domain.GuardedResult, error)
	Delete(ctx context.Context, id domain.FeatureID, env string) (domain.GuardedResult, error)
	// GetFeatureParams returns feature parameters for all environments
	GetFeatureParams(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error)
}

type FeaturesRepository interface {
	Create(ctx context.Context, envID domain.EnvironmentID, feature domain.BasicFeature) (domain.BasicFeature, error)
	GetByID(ctx context.Context, id domain.FeatureID) (domain.BasicFeature, error)
	GetByIDWithEnv(ctx context.Context, id domain.FeatureID, env string) (domain.Feature, error)
	GetByKeyWithEnv(ctx context.Context, key, env string) (domain.Feature, error)
	List(ctx context.Context, env string) ([]domain.Feature, error)
	ListByProjectID(ctx context.Context, projectID domain.ProjectID, env string) ([]domain.Feature, error)
	ListByProjectIDFiltered(
		ctx context.Context,
		projectID domain.ProjectID,
		env string,
		filter FeaturesListFilter,
	) ([]domain.Feature, int, error)
	Update(ctx context.Context, envID domain.EnvironmentID, feature domain.BasicFeature) (domain.BasicFeature, error)
	Delete(ctx context.Context, envID domain.EnvironmentID, id domain.FeatureID) error
}
