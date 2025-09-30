package domain

import "time"

type RealtimeEvent struct {
	Source         string // "pending" | "audit"
	EventID        string
	ProjectID      ProjectID
	EnvironmentID  int64
	EnvironmentKey string
	Entity         string
	EntityID       string
	Action         string
	CreatedAt      time.Time
}
