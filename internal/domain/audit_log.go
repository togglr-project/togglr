package domain

import (
	"encoding/json"
	"time"
)

type AuditActor string

const (
	AuditActorUser   AuditActor = "user"
	AuditActorSystem AuditActor = "system"
	AuditActorSDK    AuditActor = "sdk"
)

type AuditAction string

const (
	AuditActionCreate AuditAction = "create"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
)

type EntityType string

const (
	EntityFeature         EntityType = "feature"
	EntityRule            EntityType = "rule"
	EntityFlagVariant     EntityType = "flag_variant"
	EntityFeatureSchedule EntityType = "feature_schedule"
)

type AuditLogID uint64

type AuditLog struct {
	ID        AuditLogID
	ProjectID ProjectID
	FeatureID FeatureID
	EntityID  string
	RequestID string
	Entity    EntityType
	Actor     string
	Username  string
	Action    AuditAction
	OldValue  json.RawMessage
	NewValue  json.RawMessage
	CreatedAt time.Time
}

func (id AuditLogID) Uint64() uint64 {
	return uint64(id)
}

// Change represents a single change in the audit log
type Change struct {
	ID       AuditLogID
	Entity   EntityType
	EntityID string // UUID of the changed entity
	Action   AuditAction
	OldValue *json.RawMessage
	NewValue *json.RawMessage
}

// ChangeGroup represents a group of changes made in a single request
type ChangeGroup struct {
	RequestID string
	Actor     string
	Username  string
	CreatedAt time.Time
	Changes   []Change
}

// ChangesListFilter represents filter parameters for listing changes
type ChangesListFilter struct {
	ProjectID ProjectID
	Page      int
	PerPage   int
	SortBy    string
	SortDesc  bool
	Actor     *string
	Entity    *EntityType
	Action    *AuditAction
	FeatureID *FeatureID
	From      *time.Time
	To        *time.Time
}

// ChangesListResult represents paginated result of changes
type ChangesListResult struct {
	ProjectID ProjectID
	Items     []ChangeGroup
	Total     int
}
