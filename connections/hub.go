package connections

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"keryx/utils"
)

type Hub struct {
	ID uuid.UUID

	config *utils.Config
	logger *zap.Logger

	httpServer *gin.Engine

	stop chan struct{}
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

func NewHub(config *utils.Config, logger *zap.Logger) *Hub {
	h := &Hub{
		ID:     uuid.New(),
		config: config,
		logger: logger,
		stop:   make(chan struct{}),
	}

	h.httpServer = newHTTPServer(h)
	return h
}
