package request

type CartRequest struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Quantity  int64  `json:"quantity" binding:"required"`
	Size      string `json:"size" binding:"required"`
	Color     string `json:"color" binding:"required"`
	SKU       string `json:"sku" binding:"required"`
}
