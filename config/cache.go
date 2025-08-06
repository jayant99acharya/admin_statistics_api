package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			redisDB = db
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Warning: Failed to connect to Redis: %v\n", err)
		fmt.Println("Continuing without Redis cache...")
		RedisClient = nil
		return
	}

	fmt.Println("Connected to Redis successfully!")
}

func DisconnectRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}

func SetCache(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("redis client not available")
	}
	
	ctx := context.Background()
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func GetCache(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("redis client not available")
	}
	
	ctx := context.Background()
	return RedisClient.Get(ctx, key).Result()
}

func DeleteCache(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("redis client not available")
	}
	
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}