package repository

import (
	"context"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
)

type ProductFilter struct {
	CategoryID *int64
	Search     string
	Limit      int
	Offset     int
}

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindByID(ctx context.Context, id int64) (*domain.Product, error)
	FindAll(ctx context.Context, filter ProductFilter) ([]domain.Product, int64, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
}
