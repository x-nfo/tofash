package stubs

import (
	"context"
	"fmt"
	"log"
	"sync"

	"errors"
	orderEntity "tofash/internal/modules/order/entity"
	paymentEntity "tofash/internal/modules/payment/entity"
	productEntity "tofash/internal/modules/product/entity"
)

// ==========================================
// 3. IN-MEMORY SESSION (Redis Replacement for User/Auth)
// ==========================================
var sessionStore = make(map[string]string)
var sessionMu sync.Mutex

func SaveSession(token string, data []byte) error {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	sessionStore[token] = string(data)
	log.Printf("[STUB-SESSION] Saved session for token: %s...", token[:10])
	return nil
}

func GetSession(token string) (string, error) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	val, ok := sessionStore[token]
	if !ok {
		return "", errors.New("session not found")
	}
	return val, nil
}

// ==========================================
// 1. IN-MEMORY CART (Redis Replacement)
// ==========================================
type InMemoryCartRepository struct {
	store map[string][]productEntity.CartItem
	mu    sync.Mutex
}

func NewInMemoryCartRepository() *InMemoryCartRepository {
	return &InMemoryCartRepository{
		store: make(map[string][]productEntity.CartItem),
	}
}

// AddToCart matches CartRedisRepositoryInterface
func (r *InMemoryCartRepository) AddToCart(ctx context.Context, userID string, items []productEntity.CartItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Overwrite or Merge? Redis impl uses Set, so it overwrites.
	// But usually AddToCart implies merging in service layer?
	// Checking redis repo: it does `c.Client.Set(..., data, 0)`. It overwrites the key.
	// So we stick to overwrite behavior.
	r.store[userID] = items

	log.Printf("[STUB-REDIS] Saved Cart for User %s. Items: %d", userID, len(items))
	return nil
}

// GetCart matches CartRedisRepositoryInterface
func (r *InMemoryCartRepository) GetCart(ctx context.Context, userID string) ([]productEntity.CartItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.store[userID], nil
}

// RemoveFromCart matches CartRedisRepositoryInterface
func (r *InMemoryCartRepository) RemoveFromCart(ctx context.Context, userID int64, productID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Simplified key logic for Stub: just use the ID string.
	// The original repo had complex prefixing "cart:%d" etc.
	// Here we ensure internal consistency of the map.
	key := fmt.Sprintf("%d", userID)

	items, ok := r.store[key]
	if !ok {
		return nil
	}

	newItems := []productEntity.CartItem{}
	for _, item := range items {
		if item.ProductID != productID {
			newItems = append(newItems, item)
		}
	}
	r.store[key] = newItems
	return nil
}

// RemoveAllCart matches CartRedisRepositoryInterface
func (r *InMemoryCartRepository) RemoveAllCart(ctx context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%d", userID)
	delete(r.store, key)
	return nil
}

// ==========================================
// 2. NO-OP PUBLISHER (RabbitMQ Replacement)
// ==========================================
type NoOpPublisher struct{}

func NewNoOpPublisher() *NoOpPublisher {
	return &NoOpPublisher{}
}

// Product Publisher Interface
func (n *NoOpPublisher) PublishUpdateStock(productID int64, quantity int64) {
	log.Printf("[STUB-RABBITMQ] Update Stock Skipped: Product %d, Qty %d", productID, quantity)
}
func (n *NoOpPublisher) PublishProductToQueue(product productEntity.ProductEntity) error { return nil }
func (n *NoOpPublisher) DeleteProductFromQueue(productID int64) error                    { return nil }

// func (n *NoOpPublisher) PublishProductToOrder(msg interface{}) error { return nil } // Removed if unused or incorrect interface

// Order Publisher Interface
func (n *NoOpPublisher) PublishOrderToQueue(order orderEntity.OrderEntity) error {
	log.Printf("[STUB-RABBITMQ] Order Event Skipped: OrderID %d", order.ID)
	return nil
}
func (n *NoOpPublisher) PublishSendEmailUpdateStatus(email, msg, name string, id int64) error {
	return nil
}
func (n *NoOpPublisher) PublishSendPushNotifUpdateStatus(msg, topic string, id int64) error {
	return nil
}
func (n *NoOpPublisher) PublishUpdateStatus(name string, id int64, status string) error { return nil }
func (n *NoOpPublisher) PublishDeleteOrderFromQueue(id int64) error                     { return nil }

// Payment Publisher Interface
func (n *NoOpPublisher) PublishPaymentSuccess(payment paymentEntity.PaymentEntity) error {
	log.Printf("[STUB-RABBITMQ] Payment Success Event Skipped: OrderID %d", payment.OrderID)
	return nil
}
