package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Get(ctx context.Context, key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return value, nil
}

func Set(ctx context.Context, key string, value string, ttl int) error {
	if ttl > 0 {
		return redisClient.Set(ctx, key, value, 0).Err()
	}
	return redisClient.Set(ctx, key, value, 0).Err()
}

func Del(ctx context.Context, key string) error {
	return redisClient.Del(ctx, key).Err()
}

func DelWithPattern(ctx context.Context, pattern string) error {
	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := redisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			log.Printf("Failed to delete cache key: %s: %v", iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
