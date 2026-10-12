package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dhipa2poltak/go-commerce-api/internal/infrastructure/auth"
)

const UserIDContextKey = "authenticated_user_id"

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			unauthorized(c)
			return
		}

		c.Set(UserIDContextKey, claims.UserID)
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "Missing or invalid access token",
		},
	})
}
