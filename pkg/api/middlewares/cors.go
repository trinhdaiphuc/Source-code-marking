package middlewares

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4/middleware"
	"github.com/trinhdaiphuc/Source-code-marking/configs"
)

func configureCORS(echoServer *configs.EchoServer) {
	// add CORS header in response middleware
	allowOrigins := []string{}
	if os.Getenv("ENV") == "production" {
		allowOrigins = append(allowOrigins, os.Getenv("FRONT_END_SERVER_HOST"))
	} else {
		allowOrigins = append(allowOrigins, "*")
	}
	echoServer.EchoContext.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{"*"},
		ExposeHeaders: []string{
			"X-Requested-With", "X-Xss-Protection", "X-Frame-Options",
			"Content-Length", "X-Content-Type-Options", "Origin", "Upgrade",
			"Content-Type", "Accept", "Authorization", "Access-Token", "Refresh-Token",
		},
	}))
}
