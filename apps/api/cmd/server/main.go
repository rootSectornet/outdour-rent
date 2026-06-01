package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/router"
	"github.com/rentoutdoor/api/internal/infrastructure/config"
	"github.com/rentoutdoor/api/internal/infrastructure/container"
	"github.com/rentoutdoor/api/internal/infrastructure/persistence/mysql"
	"github.com/rentoutdoor/api/pkg/logger"
	"go.uber.org/zap"
)

// @title Rent Outdoor API
// @version 1.0
// @description Outdoor Equipment Rental Marketplace API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	zapLogger := logger.New(cfg.App.Env)
	defer zapLogger.Sync()

	// Connect to database
	db, err := mysql.NewConnection(&cfg.DB)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer mysql.Close(db)

	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize DI container
	c := container.New(cfg, db, zapLogger)

	// Setup router
	r := gin.New()
	router.Setup(r, c)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		zapLogger.Info("Starting server",
			zap.String("port", cfg.App.Port),
			zap.String("env", cfg.App.Env),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server exited gracefully")
}
