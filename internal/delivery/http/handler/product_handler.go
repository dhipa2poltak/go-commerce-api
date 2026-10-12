package handler

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/dto"
	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"github.com/dhipa2poltak/go-commerce-api/internal/usecase"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(uc *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: uc}
}

func productResponse(p *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
		return
	}

	product, err := h.usecase.Create(
		c.Request.Context(), req.CategoryID, req.Name, req.Description, req.Price,
	)
	if err != nil {
		handleProductError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    productResponse(product),
	})
}

func (h *ProductHandler) List(c *gin.Context) {
	page, err := queryInt(c, "page", 1, 1, 1000000)
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
		return
	}

	pageSize, err := queryInt(c, "page_size", 10, 1, 100)
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_PAGE_SIZE", "page_size must be between 1 and 100")
		return
	}

	filter := repository.ProductFilter{
		Search: strings.TrimSpace(c.Query("search")),
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	categoryValue := c.Query("category_id")
	if categoryValue != "" {
		categoryID, parseErr := strconv.ParseInt(categoryValue, 10, 64)
		if parseErr != nil || categoryID <= 0 {
			writeError(c, http.StatusBadRequest, "INVALID_CATEGORY_ID", "category_id must be a positive integer")
			return
		}

		filter.CategoryID = &categoryID
	}

	products, total, err := h.usecase.List(c.Request.Context(), filter)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list products")
		return
	}

	items := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		items = append(items, productResponse(&products[i]))
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.ProductListResponse{
			Items: items,
			Pagination: dto.PaginationResponse{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	product, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		handleProductError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    productResponse(product),
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
		return
	}

	product, err := h.usecase.Update(
		c.Request.Context(),
		id,
		req.CategoryID,
		req.Name,
		req.Description,
		req.Price,
	)
	if err != nil {
		handleProductError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    productResponse(product),
	})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		handleProductError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleProductError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidCatalogInput):
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	case errors.Is(err, usecase.ErrProductNotFound):
		writeError(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
	default:
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Product operation failed")
	}
}
