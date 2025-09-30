package contract

import (
	"context"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type RealtimeEventsRepository interface {
	// FetchAfter returns events created strictly after given timestamp, ordered by created_at ASC.
	FetchAfter(ctx context.Context, after time.Time) ([]domain.RealtimeEvent, error)
}

type RealtimeBroadcaster interface {
	// Add registers a connection interested in a specific project and environment.
	Add(projectID domain.ProjectID, envID int64, c WSConnection)
	// Remove unregisters the connection.
	Remove(projectID domain.ProjectID, envID int64, c WSConnection)
}

// WSConnection is a minimal interface to decouple WS transport from broadcaster.
// API layer should implement this for WebSocket connections.
// It should be safe for concurrent use; Send should be non-blocking or bounded.
// Close should be idempotent.
// Keep it simple for MVP.
//
// Note: We define a transport-agnostic contract to keep layers separated.
// The API layer adapts a websocket connection to this interface.
// The use case layer only works with this lightweight abstraction.
//
// 120-char limit respected across lines.
//
//nolint:revive // keep a simple name
type WSConnection interface {
	Send(msg []byte) bool
	Close()
}
