package main

import (
	"fmt"
	"tofash/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println("=== Configuration Check ===")
	fmt.Printf("App Port: %s\n", cfg.App.AppPort)
	fmt.Printf("App Env: %s\n", cfg.App.AppEnv)
	fmt.Println()
	fmt.Println("=== Database Config ===")
	fmt.Printf("Host: '%s'\n", cfg.Psql.Host)
	fmt.Printf("Port: '%s'\n", cfg.Psql.Port)
	fmt.Printf("User: '%s'\n", cfg.Psql.User)
	fmt.Printf("Password: '%s'\n", cfg.Psql.Password)
	fmt.Printf("DB Name: '%s'\n", cfg.Psql.DBName)
	fmt.Println()
	fmt.Println("=== Redis Config ===")
	fmt.Printf("Redis Host: '%s'\n", cfg.Redis.Host)
	fmt.Printf("Redis Port: '%s'\n", cfg.Redis.Port)
	fmt.Println()
	fmt.Println("=== Testing Redis Connection ===")
	redisClient := cfg.NewRedisClient()
	if redisClient == nil {
		fmt.Println("Redis client is NULL - connection failed")
	} else {
		fmt.Println("Redis client initialized successfully!")
	}
}
