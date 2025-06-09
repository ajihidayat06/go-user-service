package main

import (
	// "context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-user-service/internal/pkg/config"
	"go-user-service/internal/pkg/database"
	"go-user-service/internal/pkg/logger"
	// "go-user-service/pkg/events"

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

	// Initialize event processor
	// eventProcessor := events.NewProcessor(redis, loggerInstance)

	// Create context for graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Start background workers
	// go startEmailWorker(ctx, eventProcessor, loggerInstance)
	// go startNotificationWorker(ctx, eventProcessor, loggerInstance)
	// go startUserEventWorker(ctx, eventProcessor, loggerInstance)

	loggerInstance.Info("Worker started successfully")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	loggerInstance.Info("Shutting down worker...")

	// Cancel context to stop all workers
	// cancel()

	// Give workers time to finish
	time.Sleep(5 * time.Second)

	// Close database connections
	sqlDB, _ := db.DB()
	sqlDB.Close()
	redis.Close()

	loggerInstance.Info("Worker exited")
}

// startEmailWorker handles email sending events
// func startEmailWorker(ctx context.Context, processor *events.Processor, logger *logger.Logger) {
	// logger.Info("Starting email worker...")
	
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		logger.Info("Email worker stopped")
	// 		return
	// 	default:
	// 		// Process email events from Redis queue
	// 		err := processor.ProcessEmailEvents(ctx)
	// 		if err != nil {
	// 			logger.Error("Error processing email events: ", err)
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }
// }

// startNotificationWorker handles push notification events  
// func startNotificationWorker(ctx context.Context, processor *events.Processor, logger *logger.Logger) {
// 	logger.Info("Starting notification worker...")
	
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			logger.Info("Notification worker stopped")
// 			return
// 		default:
// 			// Process notification events from Redis queue
// 			err := processor.ProcessNotificationEvents(ctx)
// 			if err != nil {
// 				logger.Error("Error processing notification events: ", err)
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}
// }

// // startUserEventWorker handles user-related events
// func startUserEventWorker(ctx context.Context, processor *events.Processor, logger *logger.Logger) {
// 	logger.Info("Starting user event worker...")
	
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			logger.Info("User event worker stopped")
// 			return
// 		default:
// 			// Process user events from Redis queue
// 			err := processor.ProcessUserEvents(ctx)
// 			if err != nil {
// 				logger.Error("Error processing user events: ", err)
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}
// }