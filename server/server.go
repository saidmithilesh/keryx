package server

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/saidmithilesh/keryx/config"
	"go.uber.org/zap"
)

type Server struct {
	cfg    *config.Config
	logger *zap.Logger
	engine *gin.Engine
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) buildRoutes() {
	s.engine.GET("/health", s.healthCheckHandler)
}

func (s *Server) Run() {
	s.engine.Run(":" + s.cfg.AppPort)
}

func New(cfg *config.Config, logger *zap.Logger) *Server {
	isProduction := cfg.IsProduction()
	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, isProduction))
	engine.Use(ginzap.RecoveryWithZap(logger, isProduction))

	server := &Server{
		cfg:    cfg,
		logger: logger,
		engine: engine,
	}

	server.buildRoutes()
	return server
}
