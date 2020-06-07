package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler - Customize http error
func (h *Handler) CustomHTTPErrorHandler(err error, context echo.Context) {
	// Get http error code
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	h.Logger.Error("Try to access path = " + context.Path())
	h.Logger.Error(err)

	response := make(map[string]interface{})
	response["message"] = "Try to access path = " + context.Path()
	response["success"] = false
	response["error"] = err
	context.JSON(code, response)
}
