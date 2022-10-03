package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Live godoc
// @Summary Checks the API liveness
// @Description Checks the API liveness, can be used with Kubernetes Probes
// @Tags health
// @Produce json
// @Success 200 {object} string
// @Router /health/live/ [get]
func (e *Endpoints) Live(c echo.Context) error {
	c.JSON(http.StatusOK, "OK")
	return nil
}

// Ready godoc
// @Summary Checks the API readiness
// @Description Checks the API readiness, can be used with Kubernetes Probes
// @Tags health
// @Produce json
// @Success 200 {object} string
// @Router /health/ready/ [get]
func (e *Endpoints) Ready(c echo.Context) error {
	var err error
	if err = e.Config.DB.Ping(); err == nil {
		c.JSON(http.StatusOK, "READY")
		return nil
	}
	c.JSON(http.StatusNotFound, "YDAER")
	return err
}
