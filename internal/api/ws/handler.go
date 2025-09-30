package ws

import (
	"log/slog"
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

	slog.Info("WebSocket connection attempt",
		slog.String("project_id", string(projectID)),
		slog.String("remote_addr", req.RemoteAddr))

	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", slog.String("error", err.Error()))
		http.Error(writer, "upgrade failed", http.StatusBadRequest)

		return
	}

	slog.Info("WebSocket connection established")

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
		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "env_id required"),
			time.Now().Add(time.Second),
		)
		_ = conn.Close()

		return
	}

	wsConnection := newWSConn(conn)
	slog.Info("adding WebSocket connection to broadcaster")
	h.broadcaster.Add(projectID, envID, wsConnection)
	defer func() {
		slog.Info("removing WebSocket connection from broadcaster")
		h.broadcaster.Remove(projectID, envID, wsConnection)
		wsConnection.Close()
	}()

	// Set up ping/pong to keep connection alive
	conn.SetPingHandler(func(message string) error {
		slog.Info("WebSocket ping received", slog.String("message", message))
		// Respond to ping with pong
		err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err != nil {
			slog.Error("WebSocket pong failed", slog.String("error", err.Error()))
		} else {
			slog.Info("WebSocket pong sent")
		}
		return err
	})

	// Read loop to detect a client close - no deadline to avoid premature disconnections
	slog.Info("starting WebSocket read loop")
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Check if it's a normal close or an error
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Normal close, no need to log as an error
				slog.Info("WebSocket connection closed normally", slog.String("error", err.Error()))
				return
			}

			// Log other errors
			slog.Error("ws socket client error", "error", err)

			return
		}

		// Handle the message (for now, just log it)
		slog.Info("WebSocket message received",
			slog.Int("message_type", messageType),
			slog.String("message", string(message)))
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
