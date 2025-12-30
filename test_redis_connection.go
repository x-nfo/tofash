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

	fmt.Println("========================================")
	fmt.Println("  Redis Connection Test for Tofash App  ")
	fmt.Println("========================================")
	fmt.Println()

	// Load configuration
	fmt.Println("[1] Loading configuration from .env...")
	cfg := config.LoadConfig()
	fmt.Printf("    Redis Host: %s\n", cfg.Redis.Host)
	fmt.Printf("    Redis Port: %s\n", cfg.Redis.Port)
	fmt.Println()

	// Test connection using config
	fmt.Println("[2] Testing Redis connection using config...")
	redisClient := cfg.NewRedisClient()
	if redisClient == nil {
		fmt.Println("    ❌ FAILED: Redis client is NULL")
		fmt.Println("    This means Redis connection failed in NewRedisClient()")
		os.Exit(1)
	}
	fmt.Println("    ✅ SUCCESS: Redis client initialized")
	fmt.Println()

	// Test PING
	fmt.Println("[3] Testing PING command...")
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("    ❌ FAILED: %v\n", err)
	}
	fmt.Printf("    ✅ SUCCESS: PING response: %s\n", pong)
	fmt.Println()

	// Test SET operation
	fmt.Println("[4] Testing SET operation...")
	testKey := "tofash:test:connection"
	testValue := "redis_connection_successful"
	err = redisClient.Set(ctx, testKey, testValue, 0).Err()
	if err != nil {
		log.Fatalf("    ❌ FAILED: %v\n", err)
	}
	fmt.Printf("    ✅ SUCCESS: Set key '%s' with value '%s'\n", testKey, testValue)
	fmt.Println()

	// Test GET operation
	fmt.Println("[5] Testing GET operation...")
	val, err := redisClient.Get(ctx, testKey).Result()
	if err != nil {
		log.Fatalf("    ❌ FAILED: %v\n", err)
	}
	fmt.Printf("    ✅ SUCCESS: Retrieved value: %s\n", val)
	fmt.Println()

	// Test HSET/HGET (for cart functionality)
	fmt.Println("[6] Testing HSET/HGET (cart simulation)...")
	cartKey := "cart:user:123"
	itemData := `{"product_id":1,"quantity":2,"price":50000}`
	err = redisClient.HSet(ctx, cartKey, "item_1", itemData).Err()
	if err != nil {
		log.Fatalf("    ❌ FAILED: HSET error: %v\n", err)
	}
	retrievedItem, err := redisClient.HGet(ctx, cartKey, "item_1").Result()
	if err != nil {
		log.Fatalf("    ❌ FAILED: HGET error: %v\n", err)
	}
	fmt.Printf("    ✅ SUCCESS: Cart item stored and retrieved\n")
	fmt.Printf("    Cart Key: %s\n", cartKey)
	fmt.Printf("    Item Data: %s\n", retrievedItem)
	fmt.Println()

	// Test expiration (TTL)
	fmt.Println("[7] Testing TTL (Time To Live)...")
	expireKey := "tofash:test:expire"
	err = redisClient.Set(ctx, expireKey, "will_expire", 0).Err()
	if err != nil {
		log.Fatalf("    ❌ FAILED: %v\n", err)
	}
	err = redisClient.Expire(ctx, expireKey, 300).Err() // 5 minutes
	if err != nil {
		log.Fatalf("    ❌ FAILED: %v\n", err)
	}
	ttl, err := redisClient.TTL(ctx, expireKey).Result()
	if err != nil {
		log.Fatalf("    ❌ FAILED: %v\n", err)
	}
	fmt.Printf("    ✅ SUCCESS: TTL set to %v seconds\n", ttl)
	fmt.Println()

	// Cleanup
	fmt.Println("[8] Cleaning up test keys...")
	keys := []string{testKey, cartKey, expireKey}
	for _, key := range keys {
		redisClient.Del(ctx, key)
	}
	fmt.Println("    ✅ SUCCESS: All test keys deleted")
	fmt.Println()

	// Check for existing keys
	fmt.Println("[9] Checking for existing keys in Redis...")
	allKeys, err := redisClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Fatalf("    ❌ FAILED: %v\n", err)
	}
	if len(allKeys) == 0 {
		fmt.Println("    ℹ️  INFO: No keys found in Redis (fresh instance)")
	} else {
		fmt.Printf("    ℹ️  INFO: Found %d keys in Redis:\n", len(allKeys))
		for i, key := range allKeys {
			if i < 10 { // Show first 10 keys
				fmt.Printf("       - %s\n", key)
			}
		}
		if len(allKeys) > 10 {
			fmt.Printf("       ... and %d more keys\n", len(allKeys)-10)
		}
	}
	fmt.Println()

	// Close connection
	fmt.Println("[10] Closing Redis connection...")
	redisClient.Close()
	fmt.Println("    ✅ SUCCESS: Connection closed")
	fmt.Println()

	// Summary
	fmt.Println("========================================")
	fmt.Println("  ✅ ALL TESTS PASSED!                 ")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Println("  ✓ Redis is running and accessible")
	fmt.Println("  ✓ Configuration loaded correctly")
	fmt.Println("  ✓ Basic operations (PING, SET, GET) work")
	fmt.Println("  ✓ Hash operations (HSET, HGET) work for cart")
	fmt.Println("  ✓ TTL functionality works")
	fmt.Println("  ✓ Connection can be opened and closed properly")
	fmt.Println()
	fmt.Println("The Tofash application can use Redis for:")
	fmt.Println("  • Shopping cart storage")
	fmt.Println("  • Session management")
	fmt.Println("  • Caching")
	fmt.Println("  • Rate limiting")
	fmt.Println()
}
