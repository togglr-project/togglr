package dashboard

import (
	"context"
	"errors"
	"fmt"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
	repo             contract.DashboardRepository
	featuresUseCase  contract.FeaturesUseCase
	featureProcessor contract.FeatureProcessor
}

func New(
	repo contract.DashboardRepository,
	featuresUseCase contract.FeaturesUseCase,
	featureProcessor contract.FeatureProcessor,
) *Service {
	return &Service{repo: repo, featuresUseCase: featuresUseCase, featureProcessor: featureProcessor}
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

	// Compute Upcoming using top-active features IDs and feature processor
	ids, err := s.repo.TopActiveFeatureIDs(ctx, envKey, projectID, limit)
	if err != nil {
		return domain.DashboardOverview{}, fmt.Errorf("load top active feature ids: %w", err)
	}

	upcoming := make([]domain.FeatureUpcoming, 0, len(ids))
	for _, id := range ids {
		// Load extended feature for environment
		fe, err := s.featuresUseCase.GetExtendedByID(ctx, domain.FeatureID(id), envKey)
		if err != nil {
			// if feature not found, skip; otherwise return error
			if errors.Is(err, domain.ErrEntityNotFound) {
				continue
			}

			return domain.DashboardOverview{}, fmt.Errorf("get feature extended: %w", err)
		}

		enabledNext, at := s.featureProcessor.NextState(fe)
		if at.IsZero() {
			continue
		}

		next := "disabled"
		if enabledNext {
			next = "enabled"
		}
		upcoming = append(upcoming, domain.FeatureUpcoming{
			FeatureID:   fe.ID.String(),
			FeatureName: fe.Name,
			NextState:   next,
			At:          at,
		})
		if uint(len(upcoming)) >= limit {
			break
		}
	}

	return domain.DashboardOverview{
		Projects:       projects,
		Categories:     categories,
		RecentActivity: recentActivity,
		RiskyFeatures:  risky,
		PendingSummary: pending,
		Upcoming:       upcoming,
		Recent:         nil, // Feature-level recent changes not implemented yet
	}, nil
}

var _ contract.DashboardUseCase = (*Service)(nil)
