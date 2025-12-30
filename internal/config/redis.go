package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func (cfg Config) NewRedisClient() *redis.Client {
	connect := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	fmt.Printf("[Config] Attempting to connect to Redis at: %s\n", connect)
	client := redis.NewClient(&redis.Options{
		Addr: connect,
	})

	// Memastikan Redis terhubung
	_, err := client.Ping(Ctx).Result()
	if err != nil {
		fmt.Printf("[Config] Redis connection failed: %v\n", err)
		fmt.Printf("[Config] Redis Host: %s, Port: %s\n", cfg.Redis.Host, cfg.Redis.Port)
		fmt.Printf("[Config] ERROR: Redis is REQUIRED for cart functionality. Please ensure Redis is running.\n")
		return nil
	}

	fmt.Printf("[Config] Redis connection successful!\n")
	return client
}
