package flagvariants

import (
	"context"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Service struct {
	txManager db.TxManager
	repo      contract.FlagVariantsRepository
}

func New(
	txManager db.TxManager,
	repo contract.FlagVariantsRepository,
) *Service {
	return &Service{txManager: txManager, repo: repo}
}

func (s *Service) Create(ctx context.Context, v domain.FlagVariant) (domain.FlagVariant, error) {
	var created domain.FlagVariant

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.repo.Create(ctx, v)
		if err != nil {
			return fmt.Errorf("create flag variant: %w", err)
		}

		return nil
	}); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("tx create flag variant: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id domain.FlagVariantID) (domain.FlagVariant, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.FlagVariant{}, fmt.Errorf("get flag variant by id: %w", err)
	}

	return v, nil
}

func (s *Service) List(ctx context.Context) ([]domain.FlagVariant, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list flag variants: %w", err)
	}

	return items, nil
}

func (s *Service) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FlagVariant, error) {
	items, err := s.repo.ListByFeatureID(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("list flag variants by featureID: %w", err)
	}

	return items, nil
}

func (s *Service) Update(ctx context.Context, v domain.FlagVariant) (domain.FlagVariant, error) {
	var updated domain.FlagVariant

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.Update(ctx, v)
		if err != nil {
			return fmt.Errorf("update flag variant: %w", err)
		}

		return nil
	}); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("tx update flag variant: %w", err)
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id domain.FlagVariantID) error {
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("delete flag variant: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("tx delete flag variant: %w", err)
	}

	return nil
}
