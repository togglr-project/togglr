package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureParamsUseCase interface {
	GetByFeatureWithEnv(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) (domain.FeatureParams, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error)
	Update(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error)
	Delete(
		ctx context.Context,
		projectID domain.ProjectID,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) error
}

type FeatureParamsRepository interface {
	GetByFeatureWithEnv(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) (domain.FeatureParams, error)
	GetForUpdate(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) (domain.FeatureParams, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error)
	Create(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error)
	Update(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error)
	Delete(
		ctx context.Context,
		projectID domain.ProjectID,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) error
}
