package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"tofash/internal/modules/product/entity"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type CartRedisRepositoryInterface interface {
	AddToCart(ctx context.Context, userID string, items []entity.CartItem) error
	GetCart(ctx context.Context, userID string) ([]entity.CartItem, error)
	RemoveFromCart(ctx context.Context, userID int64, productID int64) error
	RemoveAllCart(ctx context.Context, userID int64) error
}

type CartRedisRepository struct {
	Client *redis.Client
}

// RemoveAllCart implements CartRedisRepositoryInterface.
func (c *CartRedisRepository) RemoveAllCart(ctx context.Context, userID int64) error {
	// DIAGNOSTIC LOG: Check if Redis client is nil
	if c.Client == nil {
		log.Errorf("[CartRedisRepository-DIAG] RemoveAllCart: Redis client is nil - Redis server is not running or not connected")
		return fmt.Errorf("redis client is not initialized - Redis server may not be running")
	}

	return c.Client.Del(ctx, fmt.Sprintf("cart:cart:%d", userID)).Err()
}

// AddToCart implements CartRedisRepositoryInterface.
func (c *CartRedisRepository) AddToCart(ctx context.Context, userID string, items []entity.CartItem) error {
	// DIAGNOSTIC LOG: Check if Redis client is nil
	if c.Client == nil {
		log.Errorf("[CartRedisRepository-DIAG] AddToCart: Redis client is nil - Redis server is not running or not connected")
		return fmt.Errorf("redis client is not initialized - Redis server may not be running")
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Errorf("[CartRedisRepository-1] AddToCart: %v", err)
		return err
	}
	return c.Client.Set(ctx, fmt.Sprintf("cart:%s", userID), data, 0).Err()
}

// GetCart implements CartRedisRepositoryInterface.
func (c *CartRedisRepository) GetCart(ctx context.Context, userID string) ([]entity.CartItem, error) {
	// DIAGNOSTIC LOG: Check if Redis client is nil
	if c.Client == nil {
		log.Errorf("[CartRedisRepository-DIAG] GetCart: Redis client is nil - Redis server is not running or not connected")
		return nil, fmt.Errorf("redis client is not initialized - Redis server may not be running")
	}

	val, err := c.Client.Get(ctx, fmt.Sprintf("cart:%s", userID)).Result()
	if err == redis.Nil {
		log.Infof("[CartRedisRepository-1] GetCart: Cart not found")
		return nil, nil
	}
	if err != nil {
		log.Errorf("[CartRedisRepository-2] GetCart: %v", err)
		return nil, err
	}
	var items []entity.CartItem
	err = json.Unmarshal([]byte(val), &items)
	if err != nil {
		log.Errorf("[CartRedisRepository-3] GetCart: %v", err)
		return nil, err
	}
	return items, nil
}

// RemoveFromCart implements CartRedisRepositoryInterface.
func (c *CartRedisRepository) RemoveFromCart(ctx context.Context, userID int64, productID int64) error {
	// DIAGNOSTIC LOG: Check if Redis client is nil
	if c.Client == nil {
		log.Errorf("[CartRedisRepository-DIAG] RemoveFromCart: Redis client is nil - Redis server is not running or not connected")
		return fmt.Errorf("redis client is not initialized - Redis server may not be running")
	}

	cart, err := c.GetCart(ctx, fmt.Sprintf("cart:%d", userID))
	if err != nil {
		log.Errorf("[CartRedisRepository-1] RemoveFromCart: %v", err)
		return err
	}

	newCart := []entity.CartItem{}
	for _, item := range cart {
		if item.ProductID != productID {
			newCart = append(newCart, item)
		}
	}

	err = c.Client.Del(ctx, fmt.Sprintf("cart:cart:%d", userID)).Err()
	if err != nil {
		log.Errorf("[CartRedisRepository-2] RemoveFromCart: %v", err)
		return err
	}

	return c.AddToCart(ctx, fmt.Sprintf("cart:%d", userID), newCart)
}

func NewCartRedisRepository(client *redis.Client) CartRedisRepositoryInterface {
	return &CartRedisRepository{
		Client: client,
	}
}
