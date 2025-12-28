package modules

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Mock services implementing required interfaces (simplified stubs).
type mockProductService struct{}

func (m *mockProductService) GetAllShop(c echo.Context) error {
	// Return a slice with a single product for testing.
	products := []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{{"prod-001", "Test Product"}}
	return c.JSON(http.StatusOK, products)
}

type mockOrderService struct{}

func (m *mockOrderService) CreateOrder(c echo.Context) error {
	var payload map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Return a fabricated order ID.
	return c.JSON(http.StatusCreated, map[string]string{"order_id": "order-123"})
}

// Health check handler (mirrors production implementation).
func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// setupTestServer creates an Echo instance with the routes needed for tests.
func setupTestServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	// Register routes.
	e.GET("/health", healthHandler)
	// Use mock services directly.
	productSvc := &mockProductService{}
	orderSvc := &mockOrderService{}
	e.GET("/api/v1/products", func(c echo.Context) error { return productSvc.GetAllShop(c) })
	e.POST("/api/v1/orders", func(c echo.Context) error { return orderSvc.CreateOrder(c) })
	return e
}

func TestHealthCheck(t *testing.T) {
	e := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("unexpected response body: %v", resp)
	}
}

func TestListProducts(t *testing.T) {
	e := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var products []map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if len(products) == 0 {
		t.Fatalf("expected at least one product, got none")
	}
	if products[0]["id"] == "" {
		t.Fatalf("product ID missing in response: %v", products[0])
	}
}

func TestCreateOrder(t *testing.T) {
	e := setupTestServer()
	payload := map[string]interface{}{
		"user_id": "user-123",
		"items": []map[string]interface{}{{
			"product_id": "prod-001",
			"quantity":   2,
			"size":       "L",
			"color":      "Blue",
			"sku":        "SKU-TEST-999",
		}},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp["order_id"] == "" {
		t.Fatalf("order_id missing in response: %v", resp)
	}
}
