package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("product not found")

type ProductUsecase struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductUsecase(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
) *ProductUsecase {
	return &ProductUsecase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (u *ProductUsecase) Create(
	ctx context.Context,
	categoryID int64,
	name, description string,
	price int64,
) (*domain.Product, error) {
	name = strings.TrimSpace(name)

	if categoryID <= 0 || name == "" || len(name) > 200 || price < 0 {
		return nil, ErrInvalidCatalogInput
	}

	category, err := u.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || category == nil {
			return nil, ErrInvalidCatalogInput
		}

		return nil, err
	}

	product := &domain.Product{
		CategoryID:  categoryID,
		Name:        name,
		Description: strings.TrimSpace(description),
		Price:       price,
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (u *ProductUsecase) List(
	ctx context.Context,
	filter repository.ProductFilter,
) ([]domain.Product, int64, error) {
	return u.productRepo.FindAll(ctx, filter)
}

func (u *ProductUsecase) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := u.productRepo.FindByID(ctx, id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, ErrProductNotFound
	}

	return product, nil
}

func (u *ProductUsecase) Update(
	ctx context.Context,
	id, categoryID int64,
	name, description string,
	price int64,
) (*domain.Product, error) {
	name = strings.TrimSpace(name)

	if id <= 0 || categoryID <= 0 || name == "" ||
		len(name) > 200 || price < 0 {

		return nil, ErrInvalidCatalogInput
	}

	product, err := u.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = u.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCatalogInput
		}

		return nil, err
	}

	product.CategoryID = categoryID
	product.Name = name
	product.Description = strings.TrimSpace(description)
	product.Price = price

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (u *ProductUsecase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidCatalogInput
	}

	if _, err := u.GetByID(ctx, id); err != nil {
		return err
	}

	return u.productRepo.Delete(ctx, id)
}
