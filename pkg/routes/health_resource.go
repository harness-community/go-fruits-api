package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Live godoc
// @Summary Checks the API liveness
// @Description Checks the API liveness, can be used with Kubernetes Probes
// @Tags health
// @Produce json
// @Success 200 {object} string
// @Router /health/live [get]
func (e *Endpoints) Live(c *gin.Context) {
	c.JSON(http.StatusOK, "live")
}

// Ready godoc
// @Summary Checks the API readiness
// @Description Checks the API readiness, can be used with Kubernetes Probes
// @Tags health
// @Produce json
// @Success 200 {object} string
// @Router /health/ready [get]
func (e *Endpoints) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, "ready")
}
