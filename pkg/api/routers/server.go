package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers"
)

func landingPage(e *echo.Echo, h *handlers.Handler) {
	e.GET("/", h.LandingPage)
	e.GET("/health-check", h.HealthCheck)
	e.GET("/health_check", h.HealthCheck)
}

func users(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/v1/users/signup", h.Signup)
	e.POST("/api/v1/users/signin", h.Signin)
	e.GET("/api/v1/users/profile", h.Profile)
	e.GET("/api/v1/users/:id", h.GetUser)
	e.GET("/api/v1/users", h.GetAllUsers)
	e.PUT("/api/v1/users", h.UpdateUser)
	e.DELETE("/api/v1/users/:id", h.DeleteUser)
}

func files(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/v1/files", h.CreateFile)
}

func Routing(e *echo.Echo, h *handlers.Handler) {
	landingPage(e, h)
	users(e, h)
	files(e, h)
}
