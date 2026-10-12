package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/middleware"
)

func authenticatedUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(middleware.UserIDContextKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)

	return userID, ok && userID > 0
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(c, http.StatusBadRequest, "INVALID_ID", "ID must be a positive integer")
		return 0, false
	}

	return id, true
}

func queryInt(
	c *gin.Context,
	key string,
	defaultValue, minValue, maxValue int,
) (int, error) {
	value := c.Query(key)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < minValue || parsed > maxValue {
		return 0, errors.New("Invalid query parameter")
	}

	return parsed, nil
}

func writeError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
