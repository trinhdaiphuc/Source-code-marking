package handlers

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler - Customize http error
func (h *Handler) CustomHTTPErrorHandler(err error, context echo.Context) {
	// Get http error code
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	h.Logger.Error("Try to access path=" + context.Path())
	h.Logger.Error(err)

	// Send error to Sentry
	go func() {
		sentry.CaptureException(err)
		sentry.Flush(time.Second * 5)
	}()

	response := make(map[string]interface{})
	response["message"] = "Try to access path=" + context.Path()
	response["success"] = false
	response["error"] = err.Error()
	context.JSON(code, response)
}
