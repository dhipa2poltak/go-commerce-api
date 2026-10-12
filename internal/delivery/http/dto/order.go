package dto

import (
	"time"
)

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required,gt=0"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderItemResponse struct {
	ID        int64 `json:"id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	UnitPrice int64 `json:"unit_price"`
	Subtotal  int64 `json:"subtotal"`
}

type OrderResponse struct {
	ID          int64               `json:id`
	UserID      int64               `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount int64               `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type OrderListResponse struct {
	Itmes      []OrderResponse `json:"items"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
}
