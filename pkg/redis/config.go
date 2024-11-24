package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func RedisClient(redisHost, redisPort, redisUsername, redisPassword string) (*redis.Client, error) {
	redisURL := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Only set Username & password if authorization enabled
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
		// Username: redisUsername,
		// Password: redisPassword,
		DB: 0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("✨ Connected to Redis.")
	return client, nil
}
