package connections

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"keryx/utils"
)

func handleHome(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

func buildRoutes(r *gin.Engine) {
	r.GET("/", handleHome)
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
