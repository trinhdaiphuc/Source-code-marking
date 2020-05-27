package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LandingPage - LandingPage handler
func (h *Handler) LandingPage(e echo.Context) error {
	h.Logger.Debug("Landing page")
	return e.String(http.StatusOK, "When come to Source code marking server!")
}
