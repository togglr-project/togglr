package ws

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Handler struct {
	broadcaster contract.RealtimeBroadcaster
}

func New(broadcaster contract.RealtimeBroadcaster) *Handler {
	return &Handler{broadcaster: broadcaster}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	// Log connection attempt
	projectID := domain.ProjectID(req.URL.Query().Get("project_id"))
	_ = req.Header.Get("Sec-WebSocket-Protocol")

	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		http.Error(writer, "upgrade failed", http.StatusBadRequest)
		return
	}

	if projectID == "" {
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "project_id required"),
			time.Now().Add(time.Second))
		_ = conn.Close()
		return
	}

	var envID int64
	if v := req.URL.Query().Get("env_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			envID = id
		}
	}
	if v := req.URL.Query().Get("environment_id"); v != "" && envID == 0 {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			envID = id
		}
	}
	if envID == 0 {
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "env_id required"),
			time.Now().Add(time.Second))
		_ = conn.Close()
		return
	}

	c := newWSConn(conn)
	h.broadcaster.Add(projectID, envID, c)
	defer func() {
		h.broadcaster.Remove(projectID, envID, c)
		c.Close()
	}()

	// Read loop to detect a client close; ignore messages
	for {
		if _, _, err := conn.NextReader(); err != nil {
			return
		}
	}
}

type wsConn struct {
	c    *websocket.Conn
	send chan []byte
}

func newWSConn(c *websocket.Conn) *wsConn {
	w := &wsConn{c: c, send: make(chan []byte, 16)}
	go w.writer()
	return w
}

func (w *wsConn) writer() {
	for msg := range w.send {
		_ = w.c.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := w.c.WriteMessage(websocket.TextMessage, msg); err != nil {
			_ = w.c.Close()
			return
		}
	}
}

func (w *wsConn) Send(msg []byte) bool {
	select {
	case w.send <- msg:
		return true
	default:
		return false
	}
}

func (w *wsConn) Close() { close(w.send) }
