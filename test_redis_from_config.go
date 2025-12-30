package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"tofash/internal/config"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Redis Connection Test ===")
	fmt.Println("Testing Redis connection using application config...")
	fmt.Println()

	// Load configuration
	cfg := config.LoadConfig()

	fmt.Printf("Config loaded:\n")
	fmt.Printf("  Redis Host: %s\n", cfg.Redis.Host)
	fmt.Printf("  Redis Port: %s\n", cfg.Redis.Port)
	fmt.Println()

	// Test connection using config
	fmt.Println("Attempting to connect to Redis...")
	redisClient := cfg.NewRedisClient()

	if redisClient == nil {
		fmt.Println("❌ FAILED: Redis client is NULL")
		fmt.Println("This means NewRedisClient() failed to connect to Redis")
		fmt.Println("The application will use in-memory cart repository as fallback")
		os.Exit(1)
	}

	fmt.Println("✅ SUCCESS: Redis client initialized")
	fmt.Println()

	// Test PING
	fmt.Println("Testing PING command...")
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ FAILED: %v\n", err)
	}
	fmt.Printf("✅ PING response: %s\n", pong)
	fmt.Println()

	// Test SET/GET
	fmt.Println("Testing SET/GET operations...")
	testKey := "tofash:test:config"
	testValue := "config_test_ok"
	err = redisClient.Set(ctx, testKey, testValue, 0).Err()
	if err != nil {
		log.Fatalf("❌ FAILED: SET error: %v\n", err)
	}
	val, err := redisClient.Get(ctx, testKey).Result()
	if err != nil {
		log.Fatalf("❌ FAILED: GET error: %v\n", err)
	}
	fmt.Printf("✅ SET/GET successful: %s = %s\n", testKey, val)
	redisClient.Del(ctx, testKey)
	fmt.Println()

	// Test cart-like hash
	fmt.Println("Testing HSET/HGET (cart simulation)...")
	cartKey := "cart:user:test:123"
	cartData := `{"product_id":1,"quantity":2}`
	err = redisClient.HSet(ctx, cartKey, "item_1", cartData).Err()
	if err != nil {
		log.Fatalf("❌ FAILED: HSET error: %v\n", err)
	}
	retrievedItem, err := redisClient.HGet(ctx, cartKey, "item_1").Result()
	if err != nil {
		log.Fatalf("❌ FAILED: HGET error: %v\n", err)
	}
	fmt.Printf("✅ HSET/HGET successful\n")
	redisClient.Del(ctx, cartKey)
	fmt.Println()

	redisClient.Close()

	fmt.Println("=== ALL TESTS PASSED ===")
	fmt.Println()
	fmt.Println("Conclusion:")
	fmt.Println("  ✓ Redis is accessible from Go application")
	fmt.Println("  ✓ Configuration is correct")
	fmt.Println("  ✓ The application CAN use Redis for cart storage")
	fmt.Println()
	fmt.Println("Note: If the running application is using in-memory cart,")
	fmt.Println("      it's because the executable was built before Redis was")
	fmt.Println("      available or there was a connection issue at startup.")
	fmt.Println("      Rebuild the application with the updated redis.go code")
	fmt.Println("      to see the new Redis connection logs.")
}
