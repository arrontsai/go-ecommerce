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
	"github.com/yourusername/ecommerce/pkg/config"
	"github.com/yourusername/ecommerce/pkg/database"
	"github.com/yourusername/ecommerce/pkg/logger"
	"github.com/yourusername/ecommerce/pkg/middleware"
	"github.com/yourusername/ecommerce/services/auth/handler"
	"github.com/yourusername/ecommerce/services/auth/repository"
	"github.com/yourusername/ecommerce/services/auth/service"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("auth-service")
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
	userRepo := repository.NewMongoUserRepository(mongoClient.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Create JWT auth middleware
	jwtMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecret)

	// Register routes
	authHandler.RegisterRoutes(router, jwtMiddleware)

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
