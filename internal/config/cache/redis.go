package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

const (
	USER_CACHE_KEY    = "users:"
	SUBJECT_CACHE_KEY = "subjects:"
)

var (
	USER_CACHE_KEY_TTL    = 3600
	SUBJECT_CACHE_KEY_TTL = 3600
)

var redisClient *redis.Client

func ConnectRedis() {
	redisURL := os.Getenv("REDIS_URL")

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse Redis URL: %v", err)
	}

	redisClient = redis.NewClient(opt)
	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")
}

func GetRedisClient() *redis.Client {
	return redisClient
}
