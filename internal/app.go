package internal

import (
	"fmt"
	"go-user-service/internal/pkg/config"
	"go-user-service/internal/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App represents the main application structure
type App struct {
	Config      *config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	Logger      logger.Logger
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, db *gorm.DB, redis *redis.Client, logger logger.Logger) *App {
	
	return &App{
		Config:      cfg,
		DB:          db,
		Redis:       redis,
		Logger:      logger,
	}
}

// Health check handler
func (a *App) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "user-service",
	})
}

// CORS middleware
func (a *App) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// Logger middleware
func (a *App) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		a.Logger.Info(fmt.Sprintf("%s %s %d %v %s",
			method,
			path,
			statusCode,
			latency,
			clientIP,
		))
	}
}

func (a *App) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(middleware.CORSMiddleware())
	// router.Use(middleware.RequestIDMiddleware())
	// router.Use(middleware.LoggerMiddleware(a.logger))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
			"env":    a.Config.App.AppEnv,
		})
	})

	// API versioning
	// v1 := router.Group("/api/v1")

	// Auth routes
	

	// Other service routes can be added here...

	return router
}