package contract

import (
	"context"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

// ErrorReportRepository provides access to error reports storage.
type ErrorReportRepository interface {
	Insert(ctx context.Context, report domain.ErrorReport) error
	CountRecent(ctx context.Context, featureID domain.FeatureID, envID domain.EnvironmentID, window time.Duration) (int, error)
	GetHealth(ctx context.Context, featureID domain.FeatureID, envID domain.EnvironmentID, window time.Duration) (domain.FeatureHealth, error)
}

// ErrorReportsUseCase encapsulates business logic around error reports and feature health.
type ErrorReportsUseCase interface {
	ReportError(
		ctx context.Context,
		projectID domain.ProjectID,
		featureKey string,
		envKey string,
		reqCtx map[domain.RuleAttribute]any,
		reportType string,
		reportMsg string,
	) (accepted bool, err error)

	GetFeatureHealth(
		ctx context.Context,
		projectID domain.ProjectID,
		featureKey string,
		envKey string,
	) (domain.FeatureHealth, error)
}
