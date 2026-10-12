package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidCatalogInput = errors.New("invalid input")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategoryNameExists  = errors.New("category name already exists")
	ErrCategoryInUse       = errors.New("category is referenced by products")
)

type CategoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(repo repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (u *CategoryUsecase) Create(ctx context.Context, name string) (*domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 100 {
		return nil, ErrInvalidCatalogInput
	}

	category := &domain.Category{Name: name}
	if err := u.repo.Create(ctx, category); err != nil {
		if isDuplicateError(err) {
			return nil, ErrCategoryNameExists
		}

		return nil, err
	}

	return category, nil
}

func (u *CategoryUsecase) List(ctx context.Context) ([]domain.Category, error) {
	return u.repo.FindAll(ctx)
}

func (u *CategoryUsecase) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	category, err := u.repo.FindByID(ctx, id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCategoryNotFound
	}

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}

func (u *CategoryUsecase) Update(ctx context.Context, id int64, name string) (*domain.Category, error) {
	name = strings.TrimSpace(name)

	if id <= 0 || name == "" || len(name) > 100 {
		return nil, ErrInvalidCatalogInput
	}

	category, err := u.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = name
	if err := u.repo.Update(ctx, category); err != nil {
		if isDuplicateError(err) {
			return nil, ErrCategoryNameExists
		}

		return nil, err
	}

	return category, nil
}

func (u *CategoryUsecase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidCatalogInput
	}

	if _, err := u.GetByID(ctx, id); err != nil {
		return err
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		if isForeignKeyError(err) {
			return ErrCategoryInUse
		}

		return err
	}

	return nil
}
