package repository

import (
	"errors"
)

var (
	ErrOrderProductNotFound   = errors.New("order product not found")
	ErrOrderInsufficientStock = errors.New("insufficient stock for order")
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled")
)
