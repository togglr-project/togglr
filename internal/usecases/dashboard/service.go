package dashboard

import (
	"context"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
	repo contract.DashboardRepository
}

func New(repo contract.DashboardRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Overview(
	ctx context.Context,
	envKey string,
	projectID *string,
	limit uint,
) (domain.DashboardOverview, error) {
	projects, err := s.repo.ProjectHealth(ctx, envKey, projectID)
	if err != nil {
		return domain.DashboardOverview{}, fmt.Errorf("load project health: %w", err)
	}

	categories, err := s.repo.CategoryHealth(ctx, envKey, projectID)
	if err != nil {
		return domain.DashboardOverview{}, fmt.Errorf("load category health: %w", err)
	}

	recentActivity, err := s.repo.RecentActivity(ctx, envKey, projectID, limit)
	if err != nil {
		return domain.DashboardOverview{}, fmt.Errorf("load recent activity: %w", err)
	}

	risky, err := s.repo.RiskyFeatures(ctx, envKey, projectID, limit)
	if err != nil {
		return domain.DashboardOverview{}, fmt.Errorf("load risky features: %w", err)
	}

	pending, err := s.repo.PendingSummary(ctx, envKey, projectID)
	if err != nil {
		return domain.DashboardOverview{}, fmt.Errorf("load pending summary: %w", err)
	}

	return domain.DashboardOverview{
		Projects:       projects,
		Categories:     categories,
		RecentActivity: recentActivity,
		RiskyFeatures:  risky,
		PendingSummary: pending,
		Upcoming:       nil, // TODO: integrate features-processor for upcoming
		Recent:         nil, // TODO: integrate audit rules for feature-level recent
	}, nil
}

var _ contract.DashboardUseCase = (*Service)(nil)
