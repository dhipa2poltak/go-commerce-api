package dto

import (
	"time"
)

type CreateInventoryRequest struct {
	Quantity int `json:"quantity" binding:"gte=0"`
}

type ChangeStockRequest struct {
	Delta int `json:"delta" binding:"required,ne=0"`
}

type InventoryResponse struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
