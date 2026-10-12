package handler

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/dto"
	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"github.com/dhipa2poltak/go-commerce-api/internal/usecase"
)

type OrderHandler struct {
	usecase *usecase.OrderUsecase
}

func NewOrderHandler(uc *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid order request")
		return
	}

	items := make([]repository.CreateOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, repository.CreateOrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := h.usecase.Create(c.Request.Context(), userID, items)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    orderResponse(order),
	})
}

func (h *OrderHandler) List(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	page, err := queryInt(c, "page", 1, 1, 1000000)
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_PAGE", "page must be positive")
		return
	}

	pageSize, err := queryInt(c, "page_size", 10, 1, 100)
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_PAGE_SIZE", "page_size must be between 1 and 100")
		return
	}

	orders, total, err := h.usecase.ListByUserID(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	items := make([]dto.OrderResponse, 0, len(orders))
	for i := range orders {
		items = append(items, orderResponse(&orders[i]))
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":       items,
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	orderID, ok := parseID(c)
	if !ok {
		return
	}

	order, err := h.usecase.GetByID(c.Request.Context(), orderID, userID)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orderResponse(order),
	})
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	orderID, ok := parseID(c)
	if !ok {
		return
	}

	order, err := h.usecase.Cancel(c.Request.Context(), orderID, userID)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orderResponse(order),
	})
}

func orderResponse(order *domain.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		})
	}

	return dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

func handleOrderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidOrderInput):
		writeError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	case errors.Is(err, usecase.ErrOrderNotFound):
		writeError(c, http.StatusNotFound, "ORDER_NOT_FOUND", "Order not found")
	case errors.Is(err, usecase.ErrOrderStock):
		writeError(c, http.StatusConflict, "INSUFFICIENT_STOCK", "Insufficient stock")
	case errors.Is(err, usecase.ErrOrderCannotCancel):
		writeError(c, http.StatusConflict, "ORDER_CANNOT_BE_CANCELLED", err.Error())
	case errors.Is(err, repository.ErrOrderProductNotFound):
		writeError(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
	case errors.Is(err, repository.ErrOrderInsufficientStock):
		writeError(c, http.StatusConflict, "INSUFFICIENT_STOCK", "Insufficient stock")
	default:
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Order operation failed")
	}
}
