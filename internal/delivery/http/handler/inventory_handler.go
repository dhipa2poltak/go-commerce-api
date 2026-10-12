package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/dto"
	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/usecase"
)

type InventoryHandler struct {
	usecase *usecase.InventoryUsecase
}

func NewInventoryHandler(uc *usecase.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{usecase: uc}
}

func inventoryResponse(inv *domain.Inventory) dto.InventoryResponse {
	return dto.InventoryResponse{
		ID:        inv.ID,
		ProductID: inv.ProductID,
		Quantity:  inv.Quantity,
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
	}
}

func (h *InventoryHandler) GetByProductID(c *gin.Context) {
	productID, ok := parseID(c)
	if !ok {
		return
	}

	inventory, err := h.usecase.GetByProductID(c.Request.Context(), productID)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    inventoryResponse(inventory),
	})
}

func (h *InventoryHandler) CreateInitialStock(c *gin.Context) {
	productID, ok := parseID(c)
	if !ok {
		return
	}

	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", "quantity must be zero or greater")
		return
	}

	inventory, err := h.usecase.CreateInitialStock(c.Request.Context(), productID, req.Quantity)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    inventoryResponse(inventory),
	})
}

func (h *InventoryHandler) ChangeStock(c *gin.Context) {
	productID, ok := parseID(c)
	if !ok {
		return
	}

	var req dto.ChangeStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", "delta must be a non-zero integer")
		return
	}

	inventory, err := h.usecase.ChangeStock(c.Request.Context(), productID, req.Delta)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    inventoryResponse(inventory),
	})
}

func handleInventoryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidStockChange):
		writeError(c, http.StatusBadRequest, "INVALID_STOCK_CHANGE", err.Error())
	case errors.Is(err, usecase.ErrProductNotFound):
		writeError(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
	case errors.Is(err, usecase.ErrInventoryNotFound):
		writeError(c, http.StatusNotFound, "INVENTORY_NOT_FOUND", err.Error())
	case errors.Is(err, usecase.ErrInventoryAlreadyExists):
		writeError(c, http.StatusConflict, "INVENTORY_ALREADY_EXISTS", err.Error())
	case errors.Is(err, usecase.ErrInsufficientStock):
		writeError(c, http.StatusConflict, "INSUFFICIENT_STOCK", err.Error())
	default:
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Inventory operation failed")
	}
}
