package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//HealthResource builds and handles HealthResource URI
// via /health
func HealthResource(rg *gin.RouterGroup) {
	health := rg.Group("/health")

	health.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ready")
	})

	health.GET("/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, "live")
	})
}
