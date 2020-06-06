package configs

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
)

type EchoServer struct {
	Logger      *internal.AppLog
	EchoContext *echo.Echo
	RedisClient *redis.Client
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
	echoServer.RedisClient = newRedisClient()
	return echoServer
}

func LoggerConfig(e *EchoServer) {
	appLog := internal.NewAppLog(os.Getenv("ENV"), os.Getenv("LOG_LEVEL"), os.Getenv("ACCESS_LOG_FILE_PATH"))
	e.Logger = appLog
}

func newRedisClient() (client *redis.Client) {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       db, // use default DB
	})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Error("Error when connecting to redis: ", err)
	} else {
		log.Info("Connected to redis ", pong)
	}
	return
}
