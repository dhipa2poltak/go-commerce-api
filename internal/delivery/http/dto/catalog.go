package dto

import (
	"time"
)

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"gte=0"`
}

type UpdateProductRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"gte=0"`
}

type ProductResponse struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ProductListResponse struct {
	Items      []ProductResponse  `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}
