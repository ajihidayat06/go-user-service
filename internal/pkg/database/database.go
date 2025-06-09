package database

import (
	"context"
	"fmt"
	"time"

	"go-user-service/internal/pkg/config"

	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.BuildDSN()

	// Configure GORM logger
	var gormLogger logger.Interface
	gormLogger = logger.Default.LogMode(logger.Info)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// NewRedisConnection creates a new Redis connection
func NewRedisConnection(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.BuildAddress(),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        20,
		MinIdleConns:    5,
		MaxIdleConns:    10,
		ConnMaxIdleTime: time.Minute * 5,
		ConnMaxLifetime: time.Hour,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		DialTimeout:     5 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

// Migration interface for database migrations
type Migration interface {
	Up() error
	Down() error
}

// Migrator handles database migrations
type Migrator struct {
	db *gorm.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// AutoMigrate runs auto migration for given models
func (m *Migrator) AutoMigrate(models ...interface{}) error {
	return m.db.AutoMigrate(models...)
}

// CreateSchema creates database schema if not exists
func (m *Migrator) CreateSchema(schemaName string) error {
	return m.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)).Error
}

// DropSchema drops database schema
func (m *Migrator) DropSchema(schemaName string) error {
	return m.db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName)).Error
}

// Health check functions
func (m *Migrator) HealthCheck() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close database connection
func (m *Migrator) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Transaction helper
func WithTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Redis helper functions
type RedisHelper struct {
	client *redis.Client
}

// NewRedisHelper creates a new Redis helper
func NewRedisHelper(client *redis.Client) *RedisHelper {
	return &RedisHelper{client: client}
}

// SetWithExpiration sets a key with expiration
func (r *RedisHelper) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (r *RedisHelper) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete deletes a key
func (r *RedisHelper) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists checks if key exists
func (r *RedisHelper) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// PushToList pushes value to list
func (r *RedisHelper) PushToList(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// PopFromList pops value from list
func (r *RedisHelper) PopFromList(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// ListLength gets list length
func (r *RedisHelper) ListLength(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

// SetHash sets hash field
func (r *RedisHelper) SetHash(ctx context.Context, key, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

// GetHash gets hash field
func (r *RedisHelper) GetHash(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// GetAllHash gets all hash fields
func (r *RedisHelper) GetAllHash(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// DeleteHashField deletes hash field
func (r *RedisHelper) DeleteHashField(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}
