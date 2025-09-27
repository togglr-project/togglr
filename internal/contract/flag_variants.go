package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FlagVariantsUseCase interface {
	Create(ctx context.Context, variant domain.FlagVariant) (domain.FlagVariant, error)
	GetByID(ctx context.Context, id domain.FlagVariantID) (domain.FlagVariant, error)
	List(ctx context.Context) ([]domain.FlagVariant, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FlagVariant, error)
	Update(ctx context.Context, variant domain.FlagVariant) (domain.FlagVariant, error)
	Delete(ctx context.Context, id domain.FlagVariantID) error
}

type FlagVariantsRepository interface {
	Create(ctx context.Context, variant domain.FlagVariant) (domain.FlagVariant, error)
	GetByID(ctx context.Context, id domain.FlagVariantID) (domain.FlagVariant, error)
	List(ctx context.Context) ([]domain.FlagVariant, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FlagVariant, error)
	ListByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) ([]domain.FlagVariant, error)
	Update(ctx context.Context, variant domain.FlagVariant) (domain.FlagVariant, error)
	Delete(ctx context.Context, id domain.FlagVariantID) error
}
