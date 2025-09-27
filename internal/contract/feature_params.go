package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureParamsUseCase interface {
	GetByFeatureAndEnvironment(ctx context.Context, featureID domain.FeatureID, environmentID domain.EnvironmentID) (domain.FeatureParams, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error)
	Update(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error)
	Delete(ctx context.Context, projectID domain.ProjectID, featureID domain.FeatureID, environmentID domain.EnvironmentID) error
}

type FeatureParamsRepository interface {
	GetByFeatureAndEnvironment(ctx context.Context, featureID domain.FeatureID, environmentID domain.EnvironmentID) (domain.FeatureParams, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error)
	Create(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error)
	Update(ctx context.Context, projectID domain.ProjectID, params domain.FeatureParams) (domain.FeatureParams, error)
	Delete(ctx context.Context, projectID domain.ProjectID, featureID domain.FeatureID, environmentID domain.EnvironmentID) error
}
