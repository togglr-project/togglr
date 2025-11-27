package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

// CustomAlgorithmsRepository handles CRUD operations for custom WASM algorithms.
type CustomAlgorithmsRepository interface {
	Create(ctx context.Context, dto domain.CustomAlgorithmDTO) (domain.CustomAlgorithm, error)
	Update(ctx context.Context, alg domain.CustomAlgorithm) error
	Delete(ctx context.Context, id domain.CustomAlgorithmID) error
	GetByID(ctx context.Context, id domain.CustomAlgorithmID) (domain.CustomAlgorithm, error)
	GetBySlug(ctx context.Context, slug string) (domain.CustomAlgorithm, error)
	List(ctx context.Context) ([]domain.CustomAlgorithm, error)
	ListByKind(ctx context.Context, kind domain.AlgorithmKind) ([]domain.CustomAlgorithm, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}

// CustomAlgorithmStatsRepository handles custom algorithm state persistence.
type CustomAlgorithmStatsRepository interface {
	Upsert(ctx context.Context, stats domain.CustomAlgorithmStats) error
	UpsertBatch(ctx context.Context, stats []domain.CustomAlgorithmStats) error
	Get(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
		algID domain.CustomAlgorithmID,
		variantKey string,
	) (domain.CustomAlgorithmStats, error)
	GetByFeature(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
		algID domain.CustomAlgorithmID,
	) ([]domain.CustomAlgorithmStats, error)
	Delete(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
		algID domain.CustomAlgorithmID,
	) error
	LoadAll(ctx context.Context) ([]domain.CustomAlgorithmStats, error)
	LoadByAlgorithmID(ctx context.Context, algID domain.CustomAlgorithmID) ([]domain.CustomAlgorithmStats, error)
}

// CustomAlgorithmsUseCase provides business logic for custom algorithms.
type CustomAlgorithmsUseCase interface {
	Create(ctx context.Context, dto domain.CustomAlgorithmDTO) (domain.CustomAlgorithm, error)
	Update(ctx context.Context, alg domain.CustomAlgorithm) error
	Delete(ctx context.Context, id domain.CustomAlgorithmID) error
	GetByID(ctx context.Context, id domain.CustomAlgorithmID) (domain.CustomAlgorithm, error)
	GetBySlug(ctx context.Context, slug string) (domain.CustomAlgorithm, error)
	List(ctx context.Context) ([]domain.CustomAlgorithm, error)
	ListByKind(ctx context.Context, kind domain.AlgorithmKind) ([]domain.CustomAlgorithm, error)
}
