package entity

type CartItem struct {
	ProductID int64  `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	Size      string `json:"size,omitempty"`
	Color     string `json:"color,omitempty"`
	SKU       string `json:"sku,omitempty"`
}
