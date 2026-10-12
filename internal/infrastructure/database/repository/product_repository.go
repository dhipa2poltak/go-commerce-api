package repository

import (
	"context"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).
		Create(product).
		Error
}

func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	var product domain.Product

	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&product, id).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindAll(ctx context.Context, filter repository.ProductFilter) ([]domain.Product, int64, error) {
	query := r.db.Model(&domain.Product{})
	//var products []domain.Product
	//query := r.db.WithContext(ctx).Preload("Category")

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []domain.Product

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.
		Order("id ASC").
		Find(&products).
		Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).
		Save(product).
		Error
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&domain.Product{}, id).
		Error
}
