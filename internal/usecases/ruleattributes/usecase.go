package ruleattributes

import (
	"context"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Service struct {
	txManager db.TxManager
	repo      contract.RuleAttributesRepository
}

func New(txManager db.TxManager, repo contract.RuleAttributesRepository) *Service {
	return &Service{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	name domain.RuleAttribute,
	description *string,
) (domain.RuleAttributeEntity, error) {
	var created domain.RuleAttributeEntity
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.repo.Create(ctx, name, description)

		return err
	})

	return created, err
}

func (s *Service) Delete(ctx context.Context, name domain.RuleAttribute) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.repo.Delete(ctx, name)
	})
}

func (s *Service) List(ctx context.Context) ([]domain.RuleAttributeEntity, error) {
	return s.repo.List(ctx)
}
