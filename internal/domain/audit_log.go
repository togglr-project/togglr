package domain

import (
	"encoding/json"
	"time"
)

type AuditAction string

const (
	AuditActionCreate AuditAction = "create"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
)

type EntityType string

const (
	EntityFeature     EntityType = "feature"
	EntityRule        EntityType = "rule"
	EntityFlagVariant EntityType = "flag_variant"
)

type AuditLogID uint64

type AuditLog struct {
	ID        AuditLogID
	FeatureID FeatureID
	Entity    EntityType
	Actor     string
	Action    AuditAction
	OldValue  json.RawMessage
	NewValue  json.RawMessage
	CreatedAt time.Time
}

func (id AuditLogID) Uint64() uint64 {
	return uint64(id)
}
