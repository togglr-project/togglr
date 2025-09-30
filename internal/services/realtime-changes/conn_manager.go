package realtime_changes

import (
	"log/slog"
	"strconv"
	"sync"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// connManager implements contract.RealtimeBroadcaster
type connManager struct {
	mu    sync.RWMutex
	conns map[string]map[contract.WSConnection]struct{}
}

func newConnManager() *connManager {
	return &connManager{conns: make(map[string]map[contract.WSConnection]struct{})}
}

func makeKey(projectID domain.ProjectID, envID int64) string {
	return string(projectID) + ":" + fmtInt(envID)
}

func (m *connManager) Add(projectID domain.ProjectID, envID int64, c contract.WSConnection) {
	key := makeKey(projectID, envID)
	m.mu.Lock()
	defer m.mu.Unlock()
	connsSet, ok := m.conns[key]
	if !ok {
		connsSet = make(map[contract.WSConnection]struct{})
		m.conns[key] = connsSet
	}
	connsSet[c] = struct{}{}
}

func (m *connManager) Remove(projectID domain.ProjectID, envID int64, c contract.WSConnection) {
	key := makeKey(projectID, envID)
	m.mu.Lock()
	defer m.mu.Unlock()
	if set, ok := m.conns[key]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(m.conns, key)
		}
	}
}

func (m *connManager) broadcastMsg(projectID domain.ProjectID, envID int64, msg []byte) {
	key := makeKey(projectID, envID)
	m.mu.RLock()
	connsSet := m.conns[key]
	m.mu.RUnlock()

	slog.Info("realtime: broadcasting to connections",
		"project_id", projectID,
		"environment_id", envID,
		"connection_count", len(connsSet))

	for connection := range connsSet {
		if !connection.Send(msg) {
			// client closed or buffer full; drop it
			slog.Info("realtime: removing closed connection")
			connection.Close()
			m.Remove(projectID, envID, connection)
		}
	}
}

func fmtInt(v int64) string {
	// fast int64 to string without fmt package to avoid allocations in hot path
	return strconv.FormatInt(v, 10)
}
