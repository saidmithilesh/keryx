package connections

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"

	"keryx/comms"
	"keryx/utils"
)

type Hub struct {
	ID uuid.UUID

	config *utils.Config
	logger *zap.Logger

	httpServer *gin.Engine

	stop   chan struct{}
	router *comms.Router

	fd          int
	connections map[int]net.Conn
	lock        *sync.RWMutex

	userFDMap map[string]int
	fdUserMap map[int]string
}

func (h *Hub) startRegistry() {
	h.logger.Info(
		"starting hub registry heartbeat",
		zap.String("hubId", h.ID.String()),
		zap.Int("registryInterval", int(h.config.RegistryInterval.Seconds())),
	)

	ticker := time.NewTicker(h.config.RegistryInterval)
	for {
		select {
		case <-ticker.C:
			h.logger.Info(
				"hub registry heartbeat",
				zap.String("hubId", h.ID.String()),
			)
		case <-h.stop:
			ticker.Stop()
			h.logger.Info(
				"stopping hub registry heartbeat",
				zap.String("hubId", h.ID.String()),
			)
			return
		}
	}
}

func (h *Hub) AddConn(userID string, conn net.Conn) error {
	connFD := getWSFD(conn)
	err := unix.EpollCtl(
		h.fd,
		syscall.EPOLL_CTL_ADD,
		connFD,
		&unix.EpollEvent{
			Events: unix.POLLIN | unix.POLLHUP,
			Fd:     int32(connFD),
		},
	)
	h.userFDMap[userID] = connFD
	h.fdUserMap[connFD] = userID
	if err != nil {
		delete(h.userFDMap, userID)
		delete(h.fdUserMap, connFD)
		return err
	}

	h.lock.Lock()
	defer h.lock.Unlock()
	h.connections[connFD] = conn
	h.logger.Info(
		"added new connection",
		zap.Int("numConnections", len(h.connections)),
	)
	return nil
}

func (h *Hub) RemoveConn(conn net.Conn) error {
	connFD := getWSFD(conn)
	err := unix.EpollCtl(h.fd, syscall.EPOLL_CTL_DEL, connFD, nil)
	if err != nil {
		return err
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.connections, connFD)
	h.logger.Info(
		"removed connection",
		zap.Int("numConnections", len(h.connections)),
	)
	return nil
}

func (h *Hub) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(h.fd, events, 100)
	if err != nil {
		return nil, err
	}
	h.lock.RLock()
	defer h.lock.RUnlock()
	var connections []net.Conn
	for i := 0; i < n; i++ {
		conn := h.connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func (h *Hub) readPump() {
	for {
		connections, err := h.Wait()
		if err != nil {
			h.logger.Error("failed to epoll wait", zap.Error(err))
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			if msg, _, err := wsutil.ReadClientData(conn); err != nil {
				if err := h.RemoveConn(conn); err != nil {
					h.logger.Error(
						"failed to remove connection",
						zap.Error(err),
					)
				}
				conn.Close()
			} else {
				// This is commented out since in demo usage, stdout is showing
				// messages sent from > 1M connections at very high rate
				h.logger.Info(
					"new message received",
					zap.String("message", string(msg)),
				)
				go h.router.Route(msg)
			}
		}
	}
}

func (h *Hub) writePump() {
	for {
		resp := <-h.router.Output
		userID := resp.SenderID.String()
		fd := h.userFDMap[userID]
		message, err := json.Marshal(resp)
		if err != nil {
			h.logger.Error(
				"failed to convert router response to bytes",
				zap.Error(err),
			)
			continue
		}
		err = wsutil.WriteServerMessage(h.connections[fd], ws.OpBinary, message)
		if err != nil {
			h.logger.Error(
				"failed to send message to client",
				zap.Error(err),
			)
			continue
		}
	}
}

func (h *Hub) Stop() {
	close(h.stop)
}

func (h *Hub) Start() {
	h.logger.Info(
		"starting hub",
		zap.String("hubId", h.ID.String()),
		zap.String("port", h.config.Port),
	)

	go h.startRegistry()
	go h.readPump()
	go h.writePump()
	h.httpServer.Run(fmt.Sprintf(":%s", h.config.Port))
}

var h *Hub

func NewHub(
	config *utils.Config,
	logger *zap.Logger,
	router *comms.Router,
) *Hub {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		logger.Fatal("failed to initialise hub", zap.Error(err))
	}

	h = &Hub{
		ID:          uuid.New(),
		config:      config,
		logger:      logger,
		stop:        make(chan struct{}),
		router:      router,
		fd:          fd,
		connections: make(map[int]net.Conn),
		lock:        &sync.RWMutex{},
		userFDMap:   make(map[string]int),
		fdUserMap:   make(map[int]string),
	}

	h.httpServer = newHTTPServer(h)
	return h
}

func getWSFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
