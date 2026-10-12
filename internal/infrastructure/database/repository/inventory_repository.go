package repository

import (
	"context"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(ctx context.Context, inventory *domain.Inventory) error {
	return r.db.WithContext(ctx).
		Create(inventory).
		Error
}

func (r *InventoryRepository) FindByProductID(ctx context.Context, productID int64) (*domain.Inventory, error) {
	var inventory domain.Inventory

	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		First(&inventory).
		Error

	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *InventoryRepository) AdjustStock(ctx context.Context, productID int64, delta int) (bool, error) {
	result := r.db.WithContext(ctx).Model(&domain.Inventory{}).
		Where("product_id = ? AND quantity + ? >= 0", productID, delta).
		UpdateColumn("quantity", gorm.Expr("quantity + ?", delta))

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
