package contract

import (
	"context"
	"time"

	"github.com/rom8726/etoggle/internal/domain"
)

// AuditLogRepository provides read access to audit log entries for change tracking.
// Services like features-processor can use it to detect project/feature changes.
//
// Note: Entries are expected to include feature_id for feature-related entities
// (feature, rule, flag_variant, feature_schedule).
// Implementations should return rows ordered by created_at ASC.
//
// Keep this contract minimal to avoid leaking storage details.
// If necessary, expand with filtering parameters later.
type AuditLogRepository interface {
	ListSince(ctx context.Context, since time.Time) ([]domain.AuditLog, error)
	ListChanges(ctx context.Context, filter domain.ChangesListFilter) (domain.ChangesListResult, error)
}
