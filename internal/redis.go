package internal

import (
	"context"

	"github.com/fatih/structs"
	"github.com/go-redis/redis/v8"
)

func RedisSetCachedWithHash(key string, redisClient *redis.Client, values interface{}) (err error) {
	ctx := context.TODO()
	convertMap := structs.Map(values)
	err = redisClient.HMSet(ctx, key, convertMap).Err()
	return
}

func RedisGetCachedWithHash(key string, redisClient *redis.Client) (data map[string]string, err error) {
	data, err = redisClient.HGetAll(context.Background(), key).Result()
	return
}
