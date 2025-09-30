package realtime

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
	repo         contract.RealtimeEventsRepository
	pollInterval time.Duration

	manager *connManager
}

func New(repo contract.RealtimeEventsRepository) *Service {
	return &Service{
		repo:         repo,
		pollInterval: 3 * time.Second,
		manager:      newConnManager(),
	}
}

// Start launches a background worker that polls the repository and broadcasts events.
func (s *Service) Start(context.Context) error {
	go s.worker(context.Background())

	return nil
}

func (s *Service) Stop(context.Context) error {
	return nil
}

func (s *Service) Broadcaster() contract.RealtimeBroadcaster { // expose broadcaster to API layer
	return s.manager
}

func (s *Service) worker(ctx context.Context) {
	lastSeen := time.Now().Add(-1 * time.Minute)
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	slog.Info("realtime worker started", "poll_interval", s.pollInterval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("realtime worker stopped")

			return
		case <-ticker.C:
			events, err := s.repo.FetchAfter(ctx, lastSeen)
			if err != nil {
				slog.Error("realtime: fetch events", "err", err)

				continue
			}

			if len(events) > 0 {
				slog.Info("realtime: found events", "count", len(events))
			}

			for _, evt := range events {
				lastSeen = evt.CreatedAt
				payload := s.toJSON(evt)
				slog.Info("realtime: broadcasting event",
					"source", evt.Source,
					"entity", evt.Entity,
					"action", evt.Action,
					"project_id", evt.ProjectID,
					"environment_id", evt.EnvironmentID)
				s.manager.broadcastMsg(evt.ProjectID, evt.EnvironmentID, payload)
			}
		}
	}
}

func (s *Service) toJSON(evt domain.RealtimeEvent) []byte {
	type out struct {
		Source      string    `json:"source"`
		Type        string    `json:"type"`
		Timestamp   time.Time `json:"timestamp"`
		ProjectID   string    `json:"project_id"`
		Environment string    `json:"environment"`
		Entity      string    `json:"entity"`
		EntityID    string    `json:"entity_id"`
		Action      string    `json:"action"`
	}

	msg := out{
		Source:      evt.Source,
		Type:        mapType(evt),
		Timestamp:   evt.CreatedAt,
		ProjectID:   string(evt.ProjectID),
		Environment: evt.EnvironmentKey,
		Entity:      evt.Entity,
		EntityID:    evt.EntityID,
		Action:      evt.Action,
	}
	b, _ := json.Marshal(msg)

	return b
}

func mapType(evt domain.RealtimeEvent) string { // simple mapping for MVP
	switch evt.Entity {
	case "feature":
		return "feature_" + evt.Action
	case "pending_change":
		return "pending_change_" + evt.Action
	default:
		return evt.Entity + "_" + evt.Action
	}
}

// connManager implements contract.RealtimeBroadcaster

type connManager struct {
	mu    sync.RWMutex
	conns map[string]map[contract.WSConnection]struct{}
}

func newConnManager() *connManager {
	return &connManager{conns: make(map[string]map[contract.WSConnection]struct{})}
}

func key(projectID domain.ProjectID, envID int64) string {
	return string(projectID) + ":" + fmtInt(envID)
}

func (m *connManager) Add(projectID domain.ProjectID, envID int64, c contract.WSConnection) {
	k := key(projectID, envID)
	m.mu.Lock()
	defer m.mu.Unlock()
	set, ok := m.conns[k]
	if !ok {
		set = make(map[contract.WSConnection]struct{})
		m.conns[k] = set
	}
	set[c] = struct{}{}
}

func (m *connManager) Remove(projectID domain.ProjectID, envID int64, c contract.WSConnection) {
	k := key(projectID, envID)
	m.mu.Lock()
	defer m.mu.Unlock()
	if set, ok := m.conns[k]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(m.conns, k)
		}
	}
}

func (m *connManager) broadcastMsg(projectID domain.ProjectID, envID int64, msg []byte) {
	k := key(projectID, envID)
	m.mu.RLock()
	set := m.conns[k]
	m.mu.RUnlock()

	slog.Info("realtime: broadcasting to connections",
		"project_id", projectID,
		"environment_id", envID,
		"connection_count", len(set))

	for c := range set {
		if !c.Send(msg) {
			// client closed or buffer full; drop it
			slog.Info("realtime: removing closed connection")
			c.Close()
			m.Remove(projectID, envID, c)
		}
	}
}

func fmtInt(v int64) string {
	// fast int64 to string without fmt package to avoid allocations in hot path
	return strconv.FormatInt(v, 10)
}
