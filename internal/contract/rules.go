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
	SyncCustomized(ctx context.Context, id domain.RuleID) (domain.Rule, error)
}

type RulesRepository interface {
	Create(ctx context.Context, rule domain.Rule) (domain.Rule, error)
	GetByID(ctx context.Context, id domain.RuleID) (domain.Rule, error)
	List(ctx context.Context) ([]domain.Rule, error)
	ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.Rule, error)
	// ListCustomizedFeatureIDsBySegment returns distinct feature IDs that have customized rules for the given segment.
	ListCustomizedFeatureIDsBySegment(ctx context.Context, segmentID domain.SegmentID) ([]domain.FeatureID, error)
	Update(ctx context.Context, rule domain.Rule) (domain.Rule, error)
	Delete(ctx context.Context, id domain.RuleID) error
}
