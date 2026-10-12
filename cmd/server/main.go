package main

import (
	"log"

	"github.com/dhipa2poltak/go-commerce-api/internal/config"
	httpdelivery "github.com/dhipa2poltak/go-commerce-api/internal/delivery/http"
	"github.com/dhipa2poltak/go-commerce-api/internal/infrastructure/auth"
	"github.com/dhipa2poltak/go-commerce-api/internal/infrastructure/database"
	repositoryImpl "github.com/dhipa2poltak/go-commerce-api/internal/infrastructure/database/repository"
	"github.com/dhipa2poltak/go-commerce-api/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatalf("failed to get SQL database instance: %v", err)
	// }
	// defer sqlDB.Close()

	log.Println("PostgreSQL connection established")

	userRepository := repositoryImpl.NewUserRepository(db)
	categoryRepository := repositoryImpl.NewCategoryRepository(db)
	productRespository := repositoryImpl.NewProductRepository(db)
	inventoryRepository := repositoryImpl.NewInventoryRepository(db)
	orderRepository := repositoryImpl.NewOrderRepository(db)

	jwtService := auth.NewJWTService(cfg.JWTSecret)

	authUsecase := usecase.NewAuthUsecase(userRepository, jwtService)
	userUsecase := usecase.NewUserUsecase(userRepository)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)
	productUsecase := usecase.NewProductUsecase(productRespository, categoryRepository)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepository, productRespository)
	orderUsecase := usecase.NewOrderUsecase(orderRepository)

	router := httpdelivery.SetupRouter(
		authUsecase,
		userUsecase,
		categoryUsecase,
		productUsecase,
		inventoryUsecase,
		orderUsecase,
		jwtService,
	)

	log.Printf("HTTP server running on port %s", cfg.AppPort)

	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
