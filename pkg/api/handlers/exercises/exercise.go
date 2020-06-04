package exercises

import (
	"github.com/go-redis/redis/v8"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseHandler struct {
	Logger      *internal.AppLog
	DB          *mongo.Client
	RedisClient *redis.Client
}

func NewExerciseHandler(logger *internal.AppLog, db *mongo.Client, redisClient *redis.Client) (u *ExerciseHandler) {
	u = &ExerciseHandler{
		Logger:      logger,
		DB:          db,
		RedisClient: redisClient,
	}
	return
}
