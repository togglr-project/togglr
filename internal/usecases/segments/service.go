package segments

import (
	"context"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Service struct {
	txManager db.TxManager
	repo      contract.SegmentsRepository
	rulesRepo contract.RulesRepository
}

func New(
	txManager db.TxManager,
	repo contract.SegmentsRepository,
	rulesRepo contract.RulesRepository,
) *Service {
	return &Service{txManager: txManager, repo: repo, rulesRepo: rulesRepo}
}

func (s *Service) Create(ctx context.Context, segment domain.Segment) (domain.Segment, error) {
	var created domain.Segment

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.repo.Create(ctx, segment)
		if err != nil {
			return fmt.Errorf("create segment: %w", err)
		}

		return nil
	}); err != nil {
		return domain.Segment{}, fmt.Errorf("tx create segment: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id domain.SegmentID) (domain.Segment, error) {
	seg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Segment{}, fmt.Errorf("get segment by id: %w", err)
	}

	return seg, nil
}

func (s *Service) ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Segment, error) {
	items, err := s.repo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list segments by projectID: %w", err)
	}

	return items, nil
}

func (s *Service) ListByProjectIDFiltered(
	ctx context.Context,
	projectID domain.ProjectID,
	filter contract.SegmentsListFilter,
) ([]domain.Segment, int, error) {
	items, total, err := s.repo.ListByProjectIDFiltered(ctx, projectID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list segments by projectID filtered: %w", err)
	}

	return items, total, nil
}

func (s *Service) Update(ctx context.Context, segment domain.Segment) (domain.Segment, error) {
	var updated domain.Segment

	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.Update(ctx, segment)
		if err != nil {
			return fmt.Errorf("update segment: %w", err)
		}

		rules, err := s.rulesRepo.ListNotCustomizedRulesBySegment(ctx, segment.ID)
		if err != nil {
			return fmt.Errorf("list rules by segment: %w", err)
		}

		for _, rule := range rules {
			rule.Conditions = updated.Conditions
			if _, err := s.rulesRepo.Update(ctx, rule); err != nil {
				return fmt.Errorf("update rule: %w", err)
			}
		}

		return nil
	}); err != nil {
		return domain.Segment{}, fmt.Errorf("tx update segment: %w", err)
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id domain.SegmentID) error {
	if err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("delete segment: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("tx delete segment: %w", err)
	}

	return nil
}

func (s *Service) ListDesyncFeatureIDs(
	ctx context.Context,
	segmentID domain.SegmentID,
) ([]domain.FeatureID, error) {
	ids, err := s.rulesRepo.ListCustomizedFeatureIDsBySegment(ctx, segmentID)
	if err != nil {
		return nil, fmt.Errorf("list desync feature ids: %w", err)
	}

	return ids, nil
}
