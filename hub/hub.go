package hub

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/saidmithilesh/keryx/config"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func isUnexpectedClose(err error) bool {
	return websocket.IsUnexpectedCloseError(
		err,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}

type Hub struct {
	cfg       *config.Config
	logger    *zap.Logger
	connStore *ConnStore
	upgrader  websocket.Upgrader

	onMessage func([]byte, string)
}

func (hub *Hub) Register(r *http.Request, w http.ResponseWriter) {
	userID := r.URL.Query().Get("userID")
	// add logic to extract user id from jwt token
	// and validate the user id before proceeding

	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	connID := uuid.New().String()
	go hub.readPump(userID, connID, conn)
	hub.connStore.AddConnection(userID, connID, conn)
}

func (hub *Hub) SetMessageListener(listener func([]byte, string)) {
	hub.onMessage = listener
}

func (hub *Hub) readPump(userID string, connID string, conn *websocket.Conn) {
	defer hub.Unregister(userID, connID)

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			errMsg := "failed to read message"
			if isUnexpectedClose(err) {
				errMsg = "connection closed unexpectedly"
			}
			hub.logger.Error(errMsg, zap.Error(err))
			break
		}
		hub.logger.Info("Received message", zap.String("message", string(message)))
		go hub.onMessage(message, userID)
	}
}

func (hub *Hub) Unregister(userID string, connID string) {
	connection, exists := hub.connStore.GetConnection(userID, connID)
	if !exists {
		return
	}
	connection.Close()
	hub.connStore.RemoveConnection(userID, connID)
}

func (hub *Hub) SendMessage(conn *websocket.Conn, message []byte) {
	conn.SetWriteDeadline(time.Now().Add(writeWait))

	writer, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	writer.Write(message)
	writer.Close()
}

func (hub *Hub) Broadcast(message []byte) {
	for _, conn := range hub.connStore.GetAllConnections() {
		go hub.SendMessage(conn, message)
	}
}

func (hub *Hub) Multicast(userIDs []string, message []byte) {
	for _, userID := range userIDs {
		for _, conn := range hub.connStore.GetUserConnections(userID) {
			go hub.SendMessage(conn, message)
		}
	}
}

func (hub *Hub) Unicast(userID, connID string, message []byte) {
	conn, exists := hub.connStore.GetConnection(userID, connID)
	if !exists {
		return
	}
	go hub.SendMessage(conn, message)
}

func New(cfg *config.Config, logger *zap.Logger) *Hub {
	return &Hub{
		cfg:       cfg,
		logger:    logger,
		connStore: NewConnStore(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}
