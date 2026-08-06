package redisconnect

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	return client
}
