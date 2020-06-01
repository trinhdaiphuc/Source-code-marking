package configs

import (
	"os"
	"runtime"
	"strconv"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
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
}
