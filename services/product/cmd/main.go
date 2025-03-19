package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/arrontsai/ecommerce/pkg/config"
	"github.com/arrontsai/ecommerce/pkg/database"
	"github.com/arrontsai/ecommerce/pkg/logger"
	"github.com/arrontsai/ecommerce/pkg/middleware"
	"github.com/arrontsai/ecommerce/services/product/handler"
	"github.com/arrontsai/ecommerce/services/product/repository"
	"github.com/arrontsai/ecommerce/services/product/service"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("product-service")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.ServiceName, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	// Connect to MongoDB
	mongoClient, err := database.NewMongoClient(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		appLogger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Close()

	// Initialize repositories
	productRepo := repository.NewMongoProductRepository(mongoClient.DB)
	categoryRepo := repository.NewMongoCategoryRepository(mongoClient.DB)

	// Initialize services
	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo, productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Create JWT auth middleware
	jwtMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecret)

	// Register routes
	productHandler.RegisterRoutes(router, jwtMiddleware)
	categoryHandler.RegisterRoutes(router, jwtMiddleware)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		appLogger.Info("Starting server", zap.Int("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exiting")
}

