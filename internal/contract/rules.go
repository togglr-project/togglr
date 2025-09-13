package contract

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type RulesUseCase interface {
	Create(ctx context.Context, rule domain.Rule) (domain.Rule, error)
	GetByID(ctx context.Context, id domain.RuleID) (domain.Rule, error)
	List(ctx context.Context) ([]domain.Rule, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.Rule, error)
	Update(ctx context.Context, rule domain.Rule) (domain.Rule, error)
	Delete(ctx context.Context, id domain.RuleID) error
}

type RulesRepository interface {
	Create(ctx context.Context, rule domain.Rule) (domain.Rule, error)
	GetByID(ctx context.Context, id domain.RuleID) (domain.Rule, error)
	List(ctx context.Context) ([]domain.Rule, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.Rule, error)
	Update(ctx context.Context, rule domain.Rule) (domain.Rule, error)
	Delete(ctx context.Context, id domain.RuleID) error
}
