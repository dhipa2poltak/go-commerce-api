package usecase

import (
	"context"
	"errors"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrInventoryNotFound      = errors.New("inventory not found")
	ErrInventoryAlreadyExists = errors.New("inventory already exists")
	ErrInvalidStockChange     = errors.New("invalid stock change")
	ErrInsufficientStock      = errors.New("insufficient stock")
)

type InventoryUsecase struct {
	inventoryRepo repository.InventoryRepository
	productRepo   repository.ProductRepository
}

func NewInventoryUsecase(
	inventoryRepo repository.InventoryRepository,
	productRepo repository.ProductRepository,
) *InventoryUsecase {
	return &InventoryUsecase{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

func (u *InventoryUsecase) GetByProductID(ctx context.Context, productID int64) (*domain.Inventory, error) {
	if productID <= 0 {
		return nil, ErrInvalidStockChange
	}

	inventory, err := u.inventoryRepo.FindByProductID(ctx, productID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInventoryNotFound
	}
	if err != nil {
		return nil, err
	}
	if inventory == nil {
		return nil, ErrInventoryNotFound
	}

	return inventory, nil
}

func (u *InventoryUsecase) CreateInitialStock(
	ctx context.Context,
	productID int64,
	quantity int,
) (*domain.Inventory, error) {
	if productID <= 0 || quantity < 0 {
		return nil, ErrInvalidStockChange
	}

	_, err := u.productRepo.FindByID(ctx, productID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	_, err = u.inventoryRepo.FindByProductID(ctx, productID)
	if err == nil {
		return nil, ErrInventoryAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	inventory := &domain.Inventory{
		ProductID: productID,
		Quantity:  quantity,
	}

	if err := u.inventoryRepo.Create(ctx, inventory); err != nil {
		if isDuplicateError(err) {
			return nil, ErrInventoryAlreadyExists
		}

		return nil, err
	}

	return inventory, nil
}

func (u *InventoryUsecase) ChangeStock(
	ctx context.Context,
	productID int64,
	delta int,
) (*domain.Inventory, error) {
	if productID <= 0 || delta == 0 {
		return nil, ErrInvalidStockChange
	}

	_, err := u.productRepo.FindByID(ctx, productID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	updated, err := u.inventoryRepo.AdjustStock(ctx, productID, delta)
	if err != nil {
		return nil, err
	}
	if !updated {
		_, err := u.inventoryRepo.FindByProductID(ctx, productID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInventoryNotFound
		}
		if err != nil {
			return nil, err
		}
		return nil, ErrInsufficientStock
	}

	return u.GetByProductID(ctx, productID)
}
