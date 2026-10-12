package repository

import (
	"context"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
)

type CreateOrderItem struct {
	ProductID int64
	Quantity  int
}

type OrderRepository interface {
	CreateWithItems(
		ctx context.Context,
		userID int64,
		items []CreateOrderItem,
	) (*domain.Order, error)

	FindByID(ctx context.Context, id int64) (*domain.Order, error)

	FindByUserID(
		ctx context.Context,
		userID int64,
		limit int,
		offset int,
	) ([]domain.Order, int64, error)

	Cancel(ctx context.Context, orderID int64, userID int64) (*domain.Order, error)
}
