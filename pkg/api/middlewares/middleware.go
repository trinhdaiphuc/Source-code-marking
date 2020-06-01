package middlewares

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/trinhdaiphuc/Source-code-marking/configs"
)

func ConfigureMiddleware(echoServer *configs.EchoServer) (err error) {
	// remove trailing slash
	echoServer.EchoContext.Pre(middleware.RemoveTrailingSlash())
	// limit body size
	echoServer.EchoContext.Use(middleware.BodyLimit("64M"))

	// CORS
	configureCORS(echoServer)

	// Logger
	configureMiddlewareLogger(echoServer)

	// JWT
	configureJWT(echoServer)

	// add recover middleware
	echoServer.EchoContext.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}))

	// enable prevent XSS, XFrame attack
	echoServer.EchoContext.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}))

	return err
}
