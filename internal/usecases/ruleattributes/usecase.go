package ruleattributes

import (
	"context"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
)

type Service struct {
	repo contract.RuleAttributesRepository
}

func New(repo contract.RuleAttributesRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	name domain.RuleAttribute,
	description *string,
) (domain.RuleAttributeEntity, error) {
	return s.repo.Create(ctx, name, description)
}

func (s *Service) Delete(ctx context.Context, name domain.RuleAttribute) error {
	return s.repo.Delete(ctx, name)
}

func (s *Service) List(ctx context.Context) ([]domain.RuleAttributeEntity, error) {
	return s.repo.List(ctx)
}
