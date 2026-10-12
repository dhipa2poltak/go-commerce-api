package usecase

import (
	"context"
	"errors"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
)

var (
	ErrInvalidOrderInput = errors.New("invalid order input")
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderStock        = errors.New("insufficient stock")
	ErrOrderCannotCancel = errors.New("order cannot be cancelled")
)

type OrderUsecase struct {
	repo repository.OrderRepository
}

func NewOrderUsecase(repo repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

func (u *OrderUsecase) Create(
	ctx context.Context,
	userID int64,
	items []repository.CreateOrderItem,
) (*domain.Order, error) {
	if userID <= 0 || len(items) == 0 {
		return nil, ErrInvalidOrderInput
	}

	for _, item := range items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return nil, ErrInvalidOrderInput
		}
	}

	return u.repo.CreateWithItems(ctx, userID, items)
}

func (u *OrderUsecase) GetByID(
	ctx context.Context,
	orderID, userID int64,
) (*domain.Order, error) {
	if orderID <= 0 || userID <= 0 {
		return nil, ErrInvalidOrderInput
	}

	order, err := u.repo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}

		return nil, err
	}

	if order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (u *OrderUsecase) ListByUserID(
	ctx context.Context,
	userID int64,
	page, pageSize int,
) ([]domain.Order, int64, error) {
	if userID <= 0 || page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, 0, ErrInvalidOrderInput
	}

	offset := (page - 1) * pageSize

	return u.repo.FindByUserID(ctx, userID, pageSize, offset)
}

func (u *OrderUsecase) Cancel(
	ctx context.Context,
	orderID, userID int64,
) (*domain.Order, error) {
	if orderID <= 0 || userID <= 0 {
		return nil, ErrInvalidOrderInput
	}

	order, err := u.repo.Cancel(ctx, orderID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrOrderNotFound):
			return nil, ErrOrderNotFound
		case errors.Is(err, repository.ErrOrderCannotBeCancelled):
			return nil, ErrOrderCannotCancel
		case errors.Is(err, repository.ErrOrderInsufficientStock):
			return nil, ErrOrderStock
		default:
			return nil, err
		}
	}

	return order, nil
}
