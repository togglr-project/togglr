package ws

import (
	"encoding/json"
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
	eventsBroadcaster contract.RealtimeBroadcaster
}

func New(eventsBroadcaster contract.RealtimeBroadcaster) *Handler {
	return &Handler{
		eventsBroadcaster: eventsBroadcaster,
	}
}

//nolint:nestif,gocognit // fix it
func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	// Log connection attempt
	projectID := domain.ProjectID(req.URL.Query().Get("project_id"))

	slog.Debug("WebSocket connection attempt",
		slog.String("project_id", string(projectID)),
		slog.String("remote_addr", req.RemoteAddr))

	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", slog.String("error", err.Error()))
		http.Error(writer, "upgrade failed", http.StatusBadRequest)

		return
	}

	slog.Debug("WebSocket connection established")

	if projectID == "" {
		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "project_id required"),
			time.Now().Add(time.Second),
		)
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
	slog.Debug("adding WebSocket connection to broadcaster",
		"project_id", projectID,
		"environment_id", envID)
	h.eventsBroadcaster.Add(projectID, envID, wsConnection)
	defer func() {
		slog.Debug("removing WebSocket connection from broadcaster")
		h.eventsBroadcaster.Remove(projectID, envID, wsConnection)
		wsConnection.Close()
	}()

	// Set up ping/pong to keep the connection alive
	conn.SetPingHandler(func(message string) error {
		slog.Debug("WebSocket ping received", slog.String("message", message))
		// Respond to ping with pong
		err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err != nil {
			slog.Error("WebSocket pong failed", slog.String("error", err.Error()))
		} else {
			slog.Debug("WebSocket pong sent")
		}

		return err
	})

	// Read loop to detect a client close - no deadline to avoid premature disconnections
	slog.Debug("starting WebSocket read loop")
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Check if it's a normal close or an error
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Normal close, no need to log as an error
				slog.Debug("WebSocket connection closed normally", slog.String("error", err.Error()))

				return
			}

			// Log other errors
			slog.Error("ws socket client error", "error", err)

			return
		}

		// Handle the message
		slog.Debug("WebSocket message received",
			slog.Int("message_type", messageType),
			slog.String("message", string(message)))

		// Handle ping messages from a client
		if messageType == websocket.TextMessage {
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
					slog.Debug("WebSocket ping received, sending pong")
					// Send pong response
					pongMsg := map[string]interface{}{"type": "pong", "timestamp": time.Now().Unix()}
					if pongData, err := json.Marshal(pongMsg); err == nil {
						if err := conn.WriteMessage(websocket.TextMessage, pongData); err != nil {
							slog.Error("WebSocket pong send failed", "error", err)
						} else {
							slog.Debug("WebSocket pong sent")
						}
					}
				}
			}
		}
	}
}
