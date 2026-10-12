package repository

import (
	"context"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	FindByID(ctx context.Context, id int64) (*domain.Category, error)
	FindAll(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64) error
}
