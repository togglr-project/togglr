package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureAlgorithmsRepository interface {
	Create(
		ctx context.Context,
		featureAlgorithm domain.FeatureAlgorithmDTO,
	) error
	Update(
		ctx context.Context,
		featureAlgorithm domain.FeatureAlgorithm,
	) error
	Delete(
		ctx context.Context,
		id domain.FeatureAlgorithmID,
	) error
	DeleteByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) error
	GetByID(ctx context.Context, id domain.FeatureAlgorithmID) (domain.FeatureAlgorithm, error)
	GetByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) (domain.FeatureAlgorithm, error)
	ListByFeatureID(
		ctx context.Context,
		featureID domain.FeatureID,
	) ([]domain.FeatureAlgorithm, error)
	ListByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) ([]domain.FeatureAlgorithm, error)
	ListExtendedByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) ([]domain.FeatureAlgorithmExtended, error)
	ListEnabled(ctx context.Context) ([]domain.FeatureAlgorithm, error)
	ListAll(ctx context.Context) ([]domain.FeatureAlgorithm, error)
	ListAllExtended(ctx context.Context) ([]domain.FeatureAlgorithmExtended, error)
	ListByProjectIDWithEnvID(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
	) ([]domain.FeatureAlgorithm, error)
}

type FeatureAlgorithmsUseCase interface {
	Create(
		ctx context.Context,
		featureAlgorithm domain.FeatureAlgorithmDTO,
	) error
	Update(
		ctx context.Context,
		featureAlgorithm domain.FeatureAlgorithm,
	) error
	DeleteByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) error
	GetByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) (domain.FeatureAlgorithm, error)
	ListByProjectIDWithEnvID(
		ctx context.Context,
		projectID domain.ProjectID,
		envID domain.EnvironmentID,
	) ([]domain.FeatureAlgorithm, error)
}
