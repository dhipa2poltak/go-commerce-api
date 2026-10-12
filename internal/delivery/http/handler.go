package http

import (
	"net/http"

	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/handler"
	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/middleware"
	"github.com/dhipa2poltak/go-commerce-api/internal/infrastructure/auth"
	"github.com/dhipa2poltak/go-commerce-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authUsecase *usecase.AuthUsecase,
	userUsecase *usecase.UserUsecase,
	categoryUsecase *usecase.CategoryUsecase,
	productUsecase *usecase.ProductUsecase,
	inventoryUsecase *usecase.InventoryUsecase,
	orderUsecase *usecase.OrderUsecase,
	jwtService *auth.JWTService,
) *gin.Engine {
	router := gin.Default()

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	productHandler := handler.NewProductHandler(productUsecase)
	inventoryHandler := handler.NewInventoryHandler(inventoryUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "go-commerce-api is running",
		})
	})

	api := router.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			protected.GET("/users/me", userHandler.Me)

			categories := protected.Group("/categories")
			{
				categories.POST("", categoryHandler.Create)
				categories.GET("", categoryHandler.List)
				categories.GET("/:id", categoryHandler.GetByID)
				categories.PUT("/:id", categoryHandler.Update)
				categories.DELETE("/:id", categoryHandler.Delete)
			}

			products := protected.Group("/products")
			{
				products.POST("", productHandler.Create)
				products.GET("", productHandler.List)
				products.GET("/:id", productHandler.GetByID)
				products.PUT("/:id", productHandler.Update)
				products.DELETE("/:id", productHandler.Delete)

				products.POST("/:id/inventory", inventoryHandler.CreateInitialStock)
				products.GET("/:id/inventory", inventoryHandler.GetByProductID)
				products.PATCH("/:id/inventory", inventoryHandler.ChangeStock)
			}

			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.Create)
				orders.GET("", orderHandler.List)
				orders.GET("/:id", orderHandler.GetByID)
				orders.PATCH("/:id/cancel", orderHandler.Cancel)
			}
		}
	}

	return router
}
