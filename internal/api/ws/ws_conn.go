package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	conn *websocket.Conn
	send chan []byte
}

func newWSConn(conn *websocket.Conn) *wsConn {
	wsConnection := &wsConn{conn: conn, send: make(chan []byte, 16)}
	go wsConnection.writer()

	return wsConnection
}

func (w *wsConn) writer() {
	for msg := range w.send {
		_ = w.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := w.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			_ = w.conn.Close()

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
