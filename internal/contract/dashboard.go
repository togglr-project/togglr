package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type DashboardRepository interface {
	ProjectHealth(ctx context.Context, envKey string, projectID *string) ([]domain.ProjectHealth, error)
	CategoryHealth(ctx context.Context, envKey string, projectID *string) ([]domain.CategoryHealth, error)
	RecentActivity(ctx context.Context, envKey string, projectID *string, limit uint) ([]domain.RecentActivity, error)
	RiskyFeatures(ctx context.Context, envKey string, projectID *string, limit uint) ([]domain.RiskyFeature, error)
	PendingSummary(ctx context.Context, envKey string, projectID *string) ([]domain.PendingSummary, error)
}

type DashboardUseCase interface {
	Overview(
		ctx context.Context,
		envKey string,
		projectID *string,
		limit uint,
	) (domain.DashboardOverview, error)
}
