package middlewares

import (
	"os"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/trinhdaiphuc/Source-code-marking/configs"
)

func configureMiddlewareLogger(e *configs.EchoServer) {
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
