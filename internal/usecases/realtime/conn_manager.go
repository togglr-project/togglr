package realtime

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
