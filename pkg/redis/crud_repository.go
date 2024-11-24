package redis

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepositories struct {
	Client *redis.Client
}

type IRedisRepositories interface {
	Set(key string, data []byte, expiredTime time.Duration, ctx context.Context) error
	Hset(key string, data string, expireAt time.Time, ctx context.Context) error
	Get(key string, ctx context.Context) (string, error)
	Del(key string, ctx context.Context) error
	GetAllByField(ctx context.Context, modelType interface{}, filterFunc func(interface{}) bool) ([]interface{}, error)
	TTL(key string, ctx context.Context) (time.Duration, error)
}

func NewRedisRepositories(client *redis.Client) *RedisRepositories {
	log.Println("🚀 Initialized Repository : Redis")
	return &RedisRepositories{
		Client: client,
	}
}

func (r *RedisRepositories) Set(key string, data []byte, expiredTime time.Duration, ctx context.Context) error {
	err := r.Client.Set(ctx, key, data, expiredTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepositories) Hset(key string, data string, expireAt time.Time, ctx context.Context) error {
	err := r.Client.Set(ctx, key, data, time.Until(expireAt)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepositories) Get(key string, ctx context.Context) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("key does not exist")
	} else if err != nil {
		return "", err
	}
	return result, nil
}

func (r *RedisRepositories) Del(key string, ctx context.Context) error {
	_, err := r.Client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetAllByField fetches all records and filters them using a custom filter function
func (r *RedisRepositories) GetAllByField(ctx context.Context, modelType interface{}, filterFunc func(interface{}) bool) ([]interface{}, error) {
	var results []interface{}
	var cursor uint64

	for {
		// Use SCAN to fetch keys from Redis
		keys, nextCursor, err := r.Client.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			// Get the value for each key
			data, err := r.Client.Get(ctx, key).Result()
			if err == redis.Nil {
				continue // Skip non-existent keys
			} else if err != nil {
				return nil, err
			}

			// Create a new instance of the model type
			model := reflect.New(reflect.TypeOf(modelType)).Interface()

			// Unmarshal JSON into the model struct
			err = json.Unmarshal([]byte(data), &model)
			if err != nil {
				continue // Skip malformed data
			}

			// Apply the filter function
			if filterFunc(model) {
				results = append(results, model)
			}
		}

		// Break if SCAN iteration is complete
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return results, nil
}

func (r *RedisRepositories) TTL(key string, ctx context.Context) (time.Duration, error) {
	duration, err := r.Client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return duration, nil
}
