package repository

import (
	"context"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventory *domain.Inventory) error
	FindByProductID(ctx context.Context, productID int64) (*domain.Inventory, error)
	AdjustStock(ctx context.Context, productID int64, delta int) (bool, error)
}
