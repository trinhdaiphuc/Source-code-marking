package middlewares

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/trinhdaiphuc/Source-code-marking/configs"
)

func configureJWT(echoServer *configs.EchoServer) {
	// JWT middleware
	echoServer.EchoContext.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if (c.Request().Method == "POST" && (c.Path() == "/api/v1/users/signin" || c.Path() == "/api/v1/users/signup")) ||
				(c.Request().Method == "GET" &&
					(c.Path() == "/api/v1/users/password" || c.Path() == "/api/v1/users/confirmation" ||
						c.Path() == "/" || c.Path() == "/metrics" ||
						c.Path() == "/health_check" || c.Path() == "/health-check")) {
				return true
			}
			return false
		},
	}))
}
