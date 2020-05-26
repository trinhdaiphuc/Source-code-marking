package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheckResponse response a message
type HealthCheckResponse struct {
	Message string `json:"message"`
}

// HealthCheck - Health Check Handler
func (h *Handler) HealthCheck(c echo.Context) error {
	resp := HealthCheckResponse{
		Message: "Everything is good!",
	}
	h.Logger.Debug("Heck-check service ", resp)
	return c.JSON(http.StatusOK, resp)
}
