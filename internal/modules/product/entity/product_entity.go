package entity

import "time"

type ProductEntity struct {
	ID           int64    `json:"id"`
	CategorySlug string   `json:"category_slug"`
	ParentID     *int64   `json:"parent_id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`  // Keep for backward compatibility
	Images       []string `json:"images"` // NEW: Multiple images support
	Description  string   `json:"description"`
	RegulerPrice float64  `json:"reguler_price"`
	SalePrice    float64  `json:"sale_price"`
	Unit         string   `json:"unit"`
	Weight       int      `json:"weight"`
	Stock        int      `json:"stock"`
	Variant      int      `json:"variant"`

	// Fashion-specific attributes
	SKU      string `json:"sku"`
	Size     string `json:"size"`
	Color    string `json:"color"`
	Material string `json:"material"`

	Status       string          `json:"status"`
	CategoryName string          `json:"category_name"`
	Child        []ProductEntity `json:"child"`
	CreatedAt    time.Time       `json:"created_at"`
}

type QueryStringProduct struct {
	Search       string
	Page         int
	Limit        int
	OrderBy      string
	OrderType    string
	CategorySlug string
	StartPrice   int64
	EndPrice     int64
	Status       string
}

type PublishOrderItemEntity struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}
