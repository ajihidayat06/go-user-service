package internal

import (
	"go-user-service/internal/pkg/config"
	"go-user-service/internal/pkg/logger"
	"go-user-service/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App represents the main application structure
type App struct {
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Logger logger.Logger
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, db *gorm.DB, redis *redis.Client, logger logger.Logger) *App {

	return &App{
		Config: cfg,
		DB:     db,
		Redis:  redis,
		Logger: logger,
	}
}

// Health check handler
func (a *App) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "user-service",
	})
}

func (a *App) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware(a.Logger))
	router.Use(gin.Recovery())
	// router.Use(middleware.RequestIDMiddleware())

	// Health check endpoint
	router.GET("/health", a.healthCheck)

	// dependency injection for handlers
	userHandler := diUser(a.DB)

	// API versioning
	v1 := router.Group("/api/v1")

	// Auth routes

	// User routes
	userHandler.RegisRoutes(v1)

	return router
}
