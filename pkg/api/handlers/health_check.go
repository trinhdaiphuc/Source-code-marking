package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck - Health Check Handler
func (h *Handler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Everything is good!")
}
