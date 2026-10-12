package repository

import (
	"context"
	"errors"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	repositoryContract "github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateWithItems(
	ctx context.Context,
	userID int64,
	requestItems []repositoryContract.CreateOrderItem,
) (*domain.Order, error) {
	var createdOrder domain.Order

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order := domain.Order{
			UserID: userID,
			Status: "CONFIRMED",
		}

		var total int64

		for _, requested := range requestItems {
			var product domain.Product
			if err := tx.First(&product, requested.ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return repositoryContract.ErrOrderProductNotFound
				}
				return err
			}

			var inventory domain.Inventory
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("product_id = ?", requested.ProductID).
				First(&inventory).
				Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return repositoryContract.ErrOrderInsufficientStock
			}
			if err != nil {
				return err
			}

			result := tx.Model(&domain.Inventory{}).
				Where("id = ? AND quantity >= ?", inventory.ID, requested.Quantity).
				UpdateColumn("quantity", gorm.Expr("quantity - ?", requested.Quantity))

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return repositoryContract.ErrOrderInsufficientStock
			}

			subtotal := product.Price * int64(requested.Quantity)
			total += subtotal

			order.Items = append(order.Items, domain.OrderItem{
				ProductID: requested.ProductID,
				Quantity:  requested.Quantity,
				UnitPrice: product.Price,
				Subtotal:  subtotal,
			})
		}

		order.TotalAmount = total

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		createdOrder = order

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, createdOrder.ID)
}

func (r *OrderRepository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	var order domain.Order

	err := r.db.WithContext(ctx).
		Preload("Items").
		//Preload("Items.Product").
		First(&order, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repositoryContract.ErrOrderNotFound
	}

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) FindByUserID(
	ctx context.Context,
	userID int64,
	limit int,
	offset int,
) ([]domain.Order, int64, error) {
	query := r.db.WithContext(ctx).Model(&domain.Order{}).Where("user_id = ?", userID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []domain.Order
	err := query.
		Preload("Items").
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).
		Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) Cancel(
	ctx context.Context,
	orderID int64,
	userID int64,
) (*domain.Order, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order domain.Order
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", orderID, userID).
			Preload("Items").
			First(&order).
			Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repositoryContract.ErrOrderNotFound
		}
		if err != nil {
			return err
		}

		if order.Status != "CONFIRMED" {
			return repositoryContract.ErrOrderCannotBeCancelled
		}

		for _, item := range order.Items {
			result := tx.Model(&domain.Inventory{}).
				Where("product_id = ?", item.ProductID).
				UpdateColumn("quantity", gorm.Expr("quantity + ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return repositoryContract.ErrOrderProductNotFound
			}
		}

		if err := tx.Model(&order).Update("status", "CANCELLED").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, orderID)
}
