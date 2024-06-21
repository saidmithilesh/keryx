package connections

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"

	"keryx/utils"
)

type Hub struct {
	ID uuid.UUID

	config *utils.Config
	logger *zap.Logger

	httpServer *gin.Engine

	stop chan struct{}

	fd          int
	connections map[int]net.Conn
	lock        *sync.RWMutex
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

func (h *Hub) AddConn(conn net.Conn) error {
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
	if err != nil {
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
	h.httpServer.Run(fmt.Sprintf(":%s", h.config.Port))
}

var h *Hub

func NewHub(config *utils.Config, logger *zap.Logger) *Hub {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		logger.Fatal("failed to initialise hub", zap.Error(err))
	}

	h = &Hub{
		ID:          uuid.New(),
		config:      config,
		logger:      logger,
		stop:        make(chan struct{}),
		fd:          fd,
		connections: make(map[int]net.Conn),
		lock:        &sync.RWMutex{},
	}

	h.httpServer = newHTTPServer(h)
	go readPump()
	return h
}

func getWSFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
