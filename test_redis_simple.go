package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Test koneksi ke Redis sesuai config di .env
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	connect := fmt.Sprintf("%s:%s", redisHost, redisPort)
	fmt.Printf("Mencoba terhubung ke Redis di: %s\n", connect)

	client := redis.NewClient(&redis.Options{
		Addr: connect,
	})

	// Ping Redis
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke Redis: %v\n", err)
	}

	fmt.Printf("✅ Redis Connection Successful! PING: %s\n", pong)

	// Test set dan get
	testKey := "test_connection_key"
	testValue := "test_connection_value"

	err = client.Set(ctx, testKey, testValue, 0).Err()
	if err != nil {
		log.Fatalf("❌ Gagal set key: %v\n", err)
	}

	val, err := client.Get(ctx, testKey).Result()
	if err != nil {
		log.Fatalf("❌ Gagal get key: %v\n", err)
	}

	fmt.Printf("✅ Set/Get test successful! Key: %s, Value: %s\n", testKey, val)

	// Cleanup
	client.Del(ctx, testKey)
	fmt.Println("✅ Redis connection test completed successfully!")
}
