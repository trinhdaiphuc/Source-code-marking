package notifications

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationHandler struct {
	Logger           *internal.AppLog
	DB               *mongo.Client
	JWTKey           string
	RedisClient      *redis.Client
	WebsocketClients map[*websocket.Conn]string
}

func NewNotificationHandler(logger *internal.AppLog, db *mongo.Client, redisClient *redis.Client, websocketClient map[*websocket.Conn]string) (u *NotificationHandler) {
	u = &NotificationHandler{
		Logger:           logger,
		DB:               db,
		RedisClient:      redisClient,
		WebsocketClients: websocketClient,
	}
	return
}
