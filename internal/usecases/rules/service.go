package rules

import (
	"context"
	"fmt"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

type Service struct {
	txManager db.TxManager
	repo      contract.RulesRepository
}

func New(
	txManager db.TxManager,
	repo contract.RulesRepository,
) *Service {
	return &Service{txManager: txManager, repo: repo}
}

func (s *Service) Create(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	var created domain.Rule
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.repo.Create(ctx, rule)
		if err != nil {
			return fmt.Errorf("create rule: %w", err)
		}
		return nil
	}); err != nil {
		return domain.Rule{}, fmt.Errorf("tx create rule: %w", err)
	}
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id domain.RuleID) (domain.Rule, error) {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Rule{}, fmt.Errorf("get rule by id: %w", err)
	}
	return r, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Rule, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}
	return items, nil
}

func (s *Service) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.Rule, error) {
	items, err := s.repo.ListByFeatureID(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("list rules by featureID: %w", err)
	}
	return items, nil
}

func (s *Service) Update(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	var updated domain.Rule
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.Update(ctx, rule)
		if err != nil {
			return fmt.Errorf("update rule: %w", err)
		}
		return nil
	}); err != nil {
		return domain.Rule{}, fmt.Errorf("tx update rule: %w", err)
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id domain.RuleID) error {
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("delete rule: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("tx delete rule: %w", err)
	}
	return nil
}
