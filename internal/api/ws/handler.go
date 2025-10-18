package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Handler struct {
	eventsBroadcaster contract.RealtimeBroadcaster
}

func New(eventsBroadcaster contract.RealtimeBroadcaster) *Handler {
	return &Handler{
		eventsBroadcaster: eventsBroadcaster,
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	connectionType := req.URL.Query().Get("type") // "realtime"

	slog.Debug("WebSocket connection attempt",
		slog.String("type", connectionType),
		slog.String("remote_addr", req.RemoteAddr))

	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", slog.String("error", err.Error()))
		http.Error(writer, "upgrade failed", http.StatusBadRequest)

		return
	}

	slog.Debug("WebSocket connection established")

	h.handleRealtimeConnection(conn, req)
}

func (h *Handler) handleRealtimeConnection(conn *websocket.Conn, req *http.Request) {
	projectID := domain.ProjectID(req.URL.Query().Get("project_id"))
	if projectID == "" {
		h.closeConnection(conn, websocket.ClosePolicyViolation, "project_id required")

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
		h.closeConnection(conn, websocket.ClosePolicyViolation, "env_id required")

		return
	}

	wsConnection := newWSConn(conn)
	slog.Debug("adding WebSocket connection to broadcaster",
		"project_id", projectID,
		"environment_id", envID)
	h.eventsBroadcaster.Add(projectID, envID, wsConnection)
	defer func() {
		slog.Debug("removing WebSocket connection from broadcaster")
		h.eventsBroadcaster.Remove(projectID, envID, wsConnection)
		wsConnection.Close()
	}()

	h.setupPingPong(conn, "realtime")

	h.handleRealtimeMessages(conn)
}

func (h *Handler) handleRealtimeMessages(conn *websocket.Conn) {
	slog.Debug("starting WebSocket read loop")
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Debug("WebSocket connection closed normally", slog.String("error", err.Error()))

				return
			}

			if !strings.Contains(err.Error(), "close 1005 (no status)") {
				slog.Error("ws socket client error", "error", err)
			}

			return
		}

		slog.Debug("WebSocket message received",
			slog.Int("message_type", messageType),
			slog.String("message", string(message)))

		h.handlePingMessage(conn, message, "realtime")
	}
}

func (h *Handler) handlePingMessage(conn *websocket.Conn, message []byte, connectionType string) {
	if len(message) == 0 {
		return
	}

	var msg map[string]any
	if err := json.Unmarshal(message, &msg); err != nil {
		return
	}

	if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
		slog.Debug("WebSocket ping received, sending pong", slog.String("type", connectionType))

		pongMsg := map[string]any{"type": "pong", "timestamp": time.Now().Unix()}
		if pongData, err := json.Marshal(pongMsg); err == nil {
			if err := conn.WriteMessage(websocket.TextMessage, pongData); err != nil {
				slog.Error("WebSocket pong send failed",
					slog.String("type", connectionType),
					slog.String("error", err.Error()))
			} else {
				slog.Debug("WebSocket pong sent", slog.String("type", connectionType))
			}
		}
	}
}

func (h *Handler) setupPingPong(conn *websocket.Conn, connectionType string) {
	conn.SetPingHandler(func(message string) error {
		slog.Debug("WebSocket ping received",
			slog.String("type", connectionType),
			slog.String("message", message))

		err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err != nil {
			slog.Error("WebSocket pong failed",
				slog.String("type", connectionType),
				slog.String("error", err.Error()))
		} else {
			slog.Debug("WebSocket pong sent", slog.String("type", connectionType))
		}

		return err
	})
}

func (h *Handler) closeConnection(conn *websocket.Conn, code int, reason string) {
	_ = conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, reason),
		time.Now().Add(time.Second),
	)
	_ = conn.Close()
}
