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

	"go-user-service/internal"
	"go-user-service/internal/pkg/config"
	"go-user-service/internal/pkg/database"
	"go-user-service/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize logger
	loggerInstance := logger.New(cfg.App.LogLevel, cfg.App.AppEnv)

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		loggerInstance.Fatal("Failed to connect to database: ", err)
	}

	// Initialize Redis
	redis, err := database.NewRedisConnection(cfg.Redis)
	if err != nil {
		loggerInstance.Fatal("Failed to connect to Redis: ", err)
	}

	// Set Gin mode
	if cfg.App.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize application
	app := internal.NewApp(cfg, db, redis, *loggerInstance)

	// Setup routes
	router := app.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.APIPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		loggerInstance.Info(fmt.Sprintf("Starting API server on port %s", cfg.Server.APIPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loggerInstance.Fatal("Failed to start server: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	loggerInstance.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		loggerInstance.Fatal("Server forced to shutdown: ", err)
	}

	// Close database connections
	sqlDB, _ := db.DB()
	sqlDB.Close()
	redis.Close()

	loggerInstance.Info("Server exited")
}