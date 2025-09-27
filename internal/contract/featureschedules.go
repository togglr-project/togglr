package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeatureSchedulesUseCase interface {
	Create(ctx context.Context, s domain.FeatureSchedule) (domain.FeatureSchedule, error)
	GetByID(ctx context.Context, id domain.FeatureScheduleID) (domain.FeatureSchedule, error)
	List(ctx context.Context) ([]domain.FeatureSchedule, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureSchedule, error)
	Update(ctx context.Context, s domain.FeatureSchedule) (domain.FeatureSchedule, error)
	Delete(ctx context.Context, id domain.FeatureScheduleID) error
}

type FeatureSchedulesRepository interface {
	Create(ctx context.Context, s domain.FeatureSchedule) (domain.FeatureSchedule, error)
	GetByID(ctx context.Context, id domain.FeatureScheduleID) (domain.FeatureSchedule, error)
	List(ctx context.Context) ([]domain.FeatureSchedule, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureSchedule, error)
	ListByFeatureIDWithEnvID(
		ctx context.Context,
		featureID domain.FeatureID,
		envID domain.EnvironmentID,
	) ([]domain.FeatureSchedule, error)
	Update(ctx context.Context, s domain.FeatureSchedule) (domain.FeatureSchedule, error)
	Delete(ctx context.Context, id domain.FeatureScheduleID) error
}
