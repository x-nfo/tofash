package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Test koneksi ke localhost:6379
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Ping Redis
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Gagal terhubung ke Redis: %v", err)
	}

	fmt.Printf("Redis Connection Successful! PING: %s\n", pong)

	// Test set dan get
	err = client.Set(ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		log.Fatalf("Gagal set key: %v", err)
	}

	val, err := client.Get(ctx, "test_key").Result()
	if err != nil {
		log.Fatalf("Gagal get key: %v", err)
	}

	fmt.Printf("Set/Get test successful! Value: %s\n", val)

	// Cleanup
	client.Del(ctx, "test_key")
	fmt.Println("Redis connection test completed successfully!")
}
