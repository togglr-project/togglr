package domain

import (
	"encoding/json"
	"time"
)

type AuditLogID uint64

type AuditLog struct {
	ID        AuditLogID
	FeatureID FeatureID
	Actor     string
	Action    string
	OldValue  json.RawMessage
	NewValue  json.RawMessage
	CreatedAt time.Time
}

func (id AuditLogID) Uint64() uint64 {
	return uint64(id)
}
