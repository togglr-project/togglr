package contract

import (
	"context"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

// AuditLogListFilter describes filters and sorting for listing audit logs at contract level.
// This keeps repository implementations interchangeable and avoids leaking storage details.
// Mirrors the needs of API ListProjectAuditLogs.
type AuditLogListFilter struct {
	ProjectID      domain.ProjectID
	EnvironmentKey *string
	Entity         *domain.EntityType
	EntityID       *string
	Actor          *string
	From           *time.Time
	To             *time.Time
	SortBy         string
	SortDesc       bool
	Page           int
	PerPage        int
}

// AuditLogRepository provides read access to audit log entries for change tracking and querying.
// Services like features-processor can use it to detect project/feature changes.
// API layer may also use it directly for read-only endpoints.
//
// Note: Entries are expected to include feature_id for feature-related entities
// (feature, rule, flag_variant, feature_schedule).
// Implementations should return rows ordered by created_at ASC for streaming.
//
// Keep this contract minimal and storage-agnostic.
// Additional list filters and methods should be defined here rather than in repository packages.
type AuditLogRepository interface {
	ListSince(ctx context.Context, since time.Time) ([]domain.AuditLog, error)
	ListChanges(ctx context.Context, filter domain.ChangesListFilter) (domain.ChangesListResult, error)
	ListByProjectIDFiltered(
		ctx context.Context,
		filter AuditLogListFilter,
	) (items []domain.AuditLog, total int, err error)
	GetByID(ctx context.Context, id domain.AuditLogID) (domain.AuditLog, error)
}
