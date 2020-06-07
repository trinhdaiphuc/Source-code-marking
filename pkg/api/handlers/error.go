package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler - Customize http error
func (h *Handler) CustomHTTPErrorHandler(err error, c echo.Context) {
	// Get http error code
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	h.Logger.Error("Try to access path = " + c.Path() + ", method = " + c.Request().Method)
	h.Logger.Error(err)

	response := make(map[string]interface{})
	response["message"] = "Try to access path = " + c.Path() + ", method = " + c.Request().Method
	response["success"] = false
	response["error"] = err
	c.JSON(code, response)
}
