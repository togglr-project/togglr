package contract

import (
	"context"

	"github.com/rom8726/etoggle/internal/domain"
)

type RuleAttributesUseCase interface {
	Create(
		ctx context.Context,
		name domain.RuleAttribute,
		description *string,
	) (domain.RuleAttributeEntity, error)
	Delete(ctx context.Context, name domain.RuleAttribute) error
	List(ctx context.Context) ([]domain.RuleAttributeEntity, error)
}

type RuleAttributesRepository interface {
	Create(
		ctx context.Context,
		name domain.RuleAttribute,
		description *string,
	) (domain.RuleAttributeEntity, error)
	Delete(ctx context.Context, name domain.RuleAttribute) error
	List(ctx context.Context) ([]domain.RuleAttributeEntity, error)
}
