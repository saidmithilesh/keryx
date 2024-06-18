package connections

import (
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
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
}

func (h *Hub) Start() {
	h.logger.Info("starting hub", zap.String("hubId", h.ID.String()), zap.String("port", h.config.Port))
	h.httpServer.Run(fmt.Sprintf(":%s", h.config.Port))
}

func (h *Hub) buildHTTPServer() {
	if h.config.Env == utils.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(ginzap.Ginzap(h.logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(h.logger, true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	h.httpServer = r
}

func NewHub(config *utils.Config, logger *zap.Logger) *Hub {
	hub := &Hub{
		ID:     uuid.New(),
		config: config,
		logger: logger,
	}

	hub.buildHTTPServer()

	return hub
}
