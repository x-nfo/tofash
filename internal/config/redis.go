package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func (cfg Config) NewRedisClient() *redis.Client {
	connect := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr: connect,
	})

	// Memastikan Redis terhubung
	_, err := client.Ping(Ctx).Result()
	if err != nil {
		fmt.Printf("[Config] Redis connection failed (running without Redis): %v\n", err)
		return nil
	}

	return client
}
