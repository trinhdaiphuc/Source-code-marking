package files

import (
	"github.com/go-redis/redis/v8"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileHandler struct {
	Logger      *internal.AppLog
	DB          *mongo.Client
	RedisClient *redis.Client
}

func NewFileHandler(logger *internal.AppLog, db *mongo.Client, redisClient *redis.Client) (u *FileHandler) {
	u = &FileHandler{
		Logger:      logger,
		DB:          db,
		RedisClient: redisClient,
	}
	return
}
