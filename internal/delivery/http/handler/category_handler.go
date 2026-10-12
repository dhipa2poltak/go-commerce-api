package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/dto"
	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/usecase"
)

type CategoryHandler struct {
	usecase *usecase.CategoryUsecase
}

func NewCategoryHandler(uc *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: uc}
}

func categoryResponse(c *domain.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
		return
	}

	category, err := h.usecase.Create(c.Request.Context(), req.Name)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    categoryResponse(category),
	})
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.usecase.List(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list categories")
		return
	}

	items := make([]dto.CategoryResponse, 0, len(categories))
	for i := range categories {
		items = append(items, categoryResponse(&categories[i]))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	category, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categoryResponse(category),
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid request body")
		return
	}

	category, err := h.usecase.Update(c.Request.Context(), id, req.Name)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categoryResponse(category),
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		handleCategoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleCategoryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidCatalogInput):
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	case errors.Is(err, usecase.ErrCategoryNotFound):
		writeError(c, http.StatusNotFound, "CATEGORY_NOT_FOUND", err.Error())
	case errors.Is(err, usecase.ErrCategoryNameExists):
		writeError(c, http.StatusConflict, "CATEGORY_NAME_EXISTS", err.Error())
	case errors.Is(err, usecase.ErrCategoryInUse):
		writeError(c, http.StatusConflict, "CATEGORY_IN_USE", err.Error())
	default:
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Category operation failed")
	}
}
