package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers"
)

func LandingPage(e *echo.Echo, h *handlers.Handler) {
	e.GET("/", h.LandingPage)
	e.GET("/health-check", h.HealthCheck)
	e.GET("/health_check", h.HealthCheck)
}

func Users(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/signup", h.Signup)
	e.POST("/api/signin", h.Signin)
	e.GET("/api/profile", h.Profile)
}
