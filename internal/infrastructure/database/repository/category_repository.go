package repository

import (
	"context"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).
		Create(category).
		Error
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*domain.Category, error) {
	var category domain.Category

	err := r.db.WithContext(ctx).
		First(&category, id).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category

	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&categories).
		Error

	return categories, err
}

func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).
		Save(category).
		Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(domain.Category{}, id).
		Error
}
