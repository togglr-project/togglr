package featureschedules

import (
	"context"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Service struct {
	txManager db.TxManager
	repo      contract.FeatureSchedulesRepository
}

func New(
	txManager db.TxManager,
	repo contract.FeatureSchedulesRepository,
) *Service {
	return &Service{txManager: txManager, repo: repo}
}

func (s *Service) Create(ctx context.Context, sch domain.FeatureSchedule) (domain.FeatureSchedule, error) {
	var created domain.FeatureSchedule
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.repo.Create(ctx, sch)
		if err != nil {
			return fmt.Errorf("create feature_schedule: %w", err)
		}

		return nil
	}); err != nil {
		return domain.FeatureSchedule{}, fmt.Errorf("tx create feature_schedule: %w", err)
	}
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id domain.FeatureScheduleID) (domain.FeatureSchedule, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.FeatureSchedule{}, fmt.Errorf("get feature_schedule by id: %w", err)
	}
	return item, nil
}

func (s *Service) List(ctx context.Context) ([]domain.FeatureSchedule, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list feature_schedules: %w", err)
	}
	return items, nil
}

func (s *Service) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureSchedule, error) {
	items, err := s.repo.ListByFeatureID(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("list feature_schedules by featureID: %w", err)
	}
	return items, nil
}

func (s *Service) Update(ctx context.Context, sch domain.FeatureSchedule) (domain.FeatureSchedule, error) {
	var updated domain.FeatureSchedule
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.Update(ctx, sch)
		if err != nil {
			return fmt.Errorf("update feature_schedule: %w", err)
		}
		return nil
	}); err != nil {
		return domain.FeatureSchedule{}, fmt.Errorf("tx update feature_schedule: %w", err)
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id domain.FeatureScheduleID) error {
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("delete feature_schedule: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("tx delete feature_schedule: %w", err)
	}
	return nil
}
