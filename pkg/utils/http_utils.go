package utils

import "github.com/labstack/echo/v4"

// NewHTTPError example
func NewHTTPError(c echo.Context, status int, err error) {
	httpErr := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	c.JSON(status, httpErr)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
