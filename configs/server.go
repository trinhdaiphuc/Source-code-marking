package configs

import (
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
)

type EchoServer struct {
	Logger      *internal.AppLog
	EchoContext *echo.Echo
}

func ConfigureMaxProcess() {
	ServerMaxProcess, err := strconv.Atoi(os.Getenv("SERVER_MAX_PROCESS"))
	if err == nil {
		runtime.GOMAXPROCS(ServerMaxProcess)
	}
	log.Info("Server is running with max process: ", runtime.GOMAXPROCS(0))
}

func NewEchoServer() *EchoServer {
	echoServer := new(EchoServer)
	echoServer.EchoContext = echo.New()
	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(echoServer.EchoContext)
	return echoServer
}

func LoggerConfig(e *EchoServer) {
	appLog := internal.NewAppLog(os.Getenv("ENV"), os.Getenv("LOG_LEVEL"), os.Getenv("ACCESS_LOG_FILE_PATH"))

	e.Logger = appLog

	if l, ok := e.EchoContext.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level} ${file} ${long_file} ${line}")
	}

	e.EchoContext.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// set log level
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		e.EchoContext.Logger.SetLevel(log.DEBUG)
	case "INFO":
		e.EchoContext.Logger.SetLevel(log.INFO)
	case "WARNING":
		e.EchoContext.Logger.SetLevel(log.WARN)
	case "ERROR":
		e.EchoContext.Logger.SetLevel(log.ERROR)
	default:
		e.EchoContext.Logger.SetLevel(log.WARN)
	}

	// set log output
	if len(os.Getenv("ACCESS_LOG_FILE_PATH")) > 0 {
		accessLogFileHandler, err := os.OpenFile(os.Getenv("ACCESS_LOG_FILE_PATH"), os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		e.EchoContext.Logger.SetOutput(accessLogFileHandler)
	}
}

func ConfigureMiddleware(echoServer *EchoServer) (err error) {
	// remove trailing slash
	echoServer.EchoContext.Pre(middleware.RemoveTrailingSlash())
	// limit body size
	echoServer.EchoContext.Use(middleware.BodyLimit("64M"))

	// add CORS header in response middleware
	allowOrigins := []string{}
	if os.Getenv("ENV") == "production" {
		allowOrigins = append(allowOrigins, os.Getenv("FRONT_END_SERVER_HOST"))
	} else {
		allowOrigins = append(allowOrigins, "*")
	}
	echoServer.EchoContext.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{"*"},
		ExposeHeaders: []string{
			"X-Requested-With", "X-Xss-Protection", "X-Frame-Options",
			"Content-Length", "X-Content-Type-Options", "Origin",
			"Content-Type", "Accept", "Authorization", "Access-Token", "Refresh-Token",
		},
	}))

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

	// JWT middleware
	echoServer.EchoContext.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Path() == "/api/v1/users/signin" || c.Path() == "/api/v1/users/signup" ||
				c.Path() == "/" || c.Path() == "/metrics" || c.Path() == "/swagger/*" ||
				c.Path() == "/health_check" || c.Path() == "/health-check" {
				return true
			}
			return false
		},
	}))

	return err
}
