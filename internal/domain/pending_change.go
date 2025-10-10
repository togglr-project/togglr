package domain

import (
	"time"
)

type PendingChangeID string

type PendingChangeStatus string

const (
	PendingChangeStatusPending   PendingChangeStatus = "pending"
	PendingChangeStatusApproved  PendingChangeStatus = "approved"
	PendingChangeStatusRejected  PendingChangeStatus = "rejected"
	PendingChangeStatusCancelled PendingChangeStatus = "cancelled"
)

type EntityAction string

const (
	EntityActionInsert EntityAction = "insert"
	EntityActionUpdate EntityAction = "update"
	EntityActionDelete EntityAction = "delete"
)

type EntityChange struct {
	Entity   string                 `json:"entity"`    // e.g., "feature", "rule"
	EntityID string                 `json:"entity_id"` // UUID of the entity
	Action   EntityAction           `json:"action"`    // insert, update, delete
	Changes  map[string]ChangeValue `json:"changes"`   // field -> {old, new}
}

type ChangeValue struct {
	Old interface{} `json:"old"`
	New interface{} `json:"new"`
}

type PendingChangeMeta struct {
	Reason            string `json:"reason,omitempty"`
	Client            string `json:"client,omitempty"`              // e.g., "ui", "api"
	Origin            string `json:"origin,omitempty"`              // e.g., "project-settings"
	SingleUserProject bool   `json:"single_user_project,omitempty"` // true if project has only 1 active user
}

type PendingChangePayload struct {
	Entities []EntityChange    `json:"entities"`
	Meta     PendingChangeMeta `json:"meta"`
}

type PendingChange struct {
	ID              PendingChangeID
	ProjectID       ProjectID
	RequestedBy     string
	RequestUserID   *int
	Change          PendingChangePayload
	Status          PendingChangeStatus
	CreatedAt       time.Time
	ApprovedBy      *string
	ApprovedUserID  *int
	ApprovedAt      *time.Time
	RejectedBy      *string
	RejectedAt      *time.Time
	RejectionReason *string
	EnvironmentID   EnvironmentID
}

type PendingChangeEntity struct {
	ID              string
	PendingChangeID PendingChangeID
	Entity          string
	EntityID        string
	CreatedAt       time.Time
}

type ProjectApprover struct {
	ProjectID ProjectID
	UserID    UserID
	Role      string
	CreatedAt time.Time
}

type ProjectSetting struct {
	ID        int
	ProjectID ProjectID
	Name      string
	Value     interface{} // JSON value
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Helper methods.
func (id PendingChangeID) String() string {
	return string(id)
}

func (status PendingChangeStatus) String() string {
	return string(status)
}

func (status PendingChangeStatus) IsValid() bool {
	return status == PendingChangeStatusPending ||
		status == PendingChangeStatusApproved ||
		status == PendingChangeStatusRejected ||
		status == PendingChangeStatusCancelled
}

func (action EntityAction) String() string {
	return string(action)
}

func (action EntityAction) IsValid() bool {
	return action == EntityActionInsert ||
		action == EntityActionUpdate ||
		action == EntityActionDelete
}

func (payload *PendingChangePayload) FeatureEntityOrFirst() string {
	if len(payload.Entities) > 0 {
		for _, entity := range payload.Entities {
			if entity.Entity == "feature" {
				return "feature"
			}
		}

		return payload.Entities[0].Entity
	}

	return "unknown"
}
