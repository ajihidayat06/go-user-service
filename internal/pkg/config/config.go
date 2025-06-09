package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for our application
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Server   ServerConfig
	App      AppConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
}

// RedisConfig holds redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Address  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshExpiresIn time.Duration
}

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	Google   OAuthProviderConfig
	Facebook OAuthProviderConfig
}

// OAuthProviderConfig holds OAuth provider configuration
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	APIPort    string
	GRPCPort   string
	WorkerPort string
}

// AppConfig holds application configuration
type AppConfig struct {
	Name     string
	Version  string
	AppEnv   string
	LogLevel string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "user_service"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn:        getEnvAsDuration("JWT_EXPIRES_IN", "24h"),
			RefreshExpiresIn: getEnvAsDuration("JWT_REFRESH_EXPIRES_IN", "168h"), // 7 days
		},
		OAuth: OAuthConfig{
			Google: OAuthProviderConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
			},
			Facebook: OAuthProviderConfig{
				ClientID:     getEnv("FACEBOOK_CLIENT_ID", ""),
				ClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("FACEBOOK_REDIRECT_URL", ""),
			},
		},
		Server: ServerConfig{
			APIPort:    getEnv("API_PORT", "8080"),
			GRPCPort:   getEnv("GRPC_PORT", "9090"),
			WorkerPort: getEnv("WORKER_PORT", "8081"),
		},
		App: AppConfig{
			Name:     getEnv("APP_NAME", "user-service"),
			Version:  getEnv("APP_VERSION", "1.0.0"),
			AppEnv:   getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// BuildDSN builds database DSN from config
func (d *DatabaseConfig) BuildDSN() string {
	if d.DSN != "" {
		return d.DSN
	}
	return "host=" + d.Host + " port=" + d.Port + " user=" + d.User + " password=" + d.Password + " dbname=" + d.DBName + " sslmode=" + d.SSLMode
}

// BuildAddress builds redis address from config
func (r *RedisConfig) BuildAddress() string {
	if r.Address != "" {
		return r.Address
	}
	return r.Host + ":" + r.Port
}

// IsDevelopment checks if app is in development mode
func (a *AppConfig) IsDevelopment() bool {
	return a.AppEnv == "development" || a.AppEnv == "dev"
}

// IsProduction checks if app is in production mode
func (a *AppConfig) IsProduction() bool {
	return a.AppEnv == "production" || a.AppEnv == "prod"
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 24 * time.Hour // fallback
}