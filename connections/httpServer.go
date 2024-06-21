package connections

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"go.uber.org/zap"

	"keryx/utils"
)

func handleHome(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

func handleWS(c *gin.Context) {
	var r *http.Request = c.Request
	var w http.ResponseWriter = c.Writer
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	userID := c.Query("userId")
	if err := h.AddConn(userID, conn); err != nil {
		h.logger.Error("failed to add connection", zap.Error(err))
		conn.Close()
	}
}

func buildRoutes(r *gin.Engine) {
	r.GET("/", handleHome)
	r.GET("/ws", handleWS)
}

func newHTTPServer(h *Hub) *gin.Engine {
	if h.config.Env == utils.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(ginzap.Ginzap(h.logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(h.logger, true))

	buildRoutes(r)
	return r
}
