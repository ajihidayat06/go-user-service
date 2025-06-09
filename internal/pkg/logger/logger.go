package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// Fields type for structured logging
type Fields map[string]interface{}

// New creates a new logger instance
func New(level, env string) *Logger {
	log := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

	// Set formatter based on environment
	if env == "production" || env == "prod" {
		// JSON formatter for production
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	} else {
		// Text formatter for development
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set output to stdout
	log.SetOutput(os.Stdout)

	// Add caller information
	log.SetReportCaller(true)

	return &Logger{Logger: log}
}

// WithFields adds fields to logger
func (l *Logger) WithFields(fields Fields) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// WithField adds a single field to logger
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithError adds error field to logger
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithContext adds context to logger
func (l *Logger) WithContext(ctx interface{}) *logrus.Entry {
	return l.Logger.WithField("context", ctx)
}

// HTTP request logging helpers
func (l *Logger) LogHTTPRequest(method, path, userAgent, clientIP string, statusCode int, duration int64) {
	l.WithFields(Fields{
		"method":      method,
		"path":       path,
		"user_agent": userAgent,
		"client_ip":  clientIP,
		"status":     statusCode,
		"duration":   duration,
		"type":       "http_request",
	}).Info("HTTP Request")
}

// Database operation logging helpers
func (l *Logger) LogDBOperation(operation, table string, duration int64, err error) {
	entry := l.WithFields(Fields{
		"operation": operation,
		"table":     table,
		"duration":  duration,
		"type":      "db_operation",
	})

	if err != nil {
		entry.WithError(err).Error("Database operation failed")
	} else {
		entry.Info("Database operation completed")
	}
}

// Service operation logging helpers
func (l *Logger) LogServiceOperation(service, operation string, duration int64, err error) {
	entry := l.WithFields(Fields{
		"service":   service,
		"operation": operation,
		"duration":  duration,
		"type":      "service_operation",
	})

	if err != nil {
		entry.WithError(err).Error("Service operation failed")
	} else {
		entry.Info("Service operation completed")
	}
}

// Authentication logging helpers
func (l *Logger) LogAuthOperation(operation, userID, method string, success bool, err error) {
	entry := l.WithFields(Fields{
		"operation": operation,
		"user_id":   userID,
		"method":    method,
		"success":   success,
		"type":      "auth_operation",
	})

	if err != nil {
		entry.WithError(err).Warn("Authentication operation failed")
	} else {
		entry.Info("Authentication operation completed")
	}
}

// Security logging helpers
func (l *Logger) LogSecurityEvent(event, userID, clientIP, details string) {
	l.WithFields(Fields{
		"event":     event,
		"user_id":   userID,
		"client_ip": clientIP,
		"details":   details,
		"type":      "security_event",
	}).Warn("Security event detected")
}

// Business logic logging helpers
func (l *Logger) LogBusinessEvent(event, userID string, data map[string]interface{}) {
	fields := Fields{
		"event":   event,
		"user_id": userID,
		"type":    "business_event",
	}

	// Add data fields
	for k, v := range data {
		fields[k] = v
	}

	l.WithFields(fields).Info("Business event occurred")
}

// Performance logging helpers
func (l *Logger) LogPerformance(operation string, duration int64, metadata map[string]interface{}) {
	fields := Fields{
		"operation": operation,
		"duration":  duration,
		"type":      "performance",
	}

	// Add metadata fields
	for k, v := range metadata {
		fields[k] = v
	}

	level := l.Logger.Level
	entry := l.WithFields(fields)

	// Log as warning if duration is too long
	if duration > 5000 { // 5 seconds
		entry.Warn("Slow operation detected")
	} else if level == logrus.DebugLevel {
		entry.Debug("Performance metrics")
	}
}

// Error tracking helpers
func (l *Logger) LogError(err error, context string, data map[string]interface{}) {
	fields := Fields{
		"context": context,
		"type":    "error",
	}

	// Add data fields
	for k, v := range data {
		fields[k] = v
	}

	l.WithFields(fields).WithError(err).Error("Error occurred")
}

// Health check logging
func (l *Logger) LogHealthCheck(service string, healthy bool, err error) {
	entry := l.WithFields(Fields{
		"service": service,
		"healthy": healthy,
		"type":    "health_check",
	})

	if err != nil {
		entry.WithError(err).Error("Health check failed")
	} else if healthy {
		entry.Info("Health check passed")
	} else {
		entry.Warn("Health check failed")
	}
}

// Startup/Shutdown logging
func (l *Logger) LogStartup(service, version string, config map[string]interface{}) {
	fields := Fields{
		"service": service,
		"version": version,
		"type":    "startup",
	}

	// Add config fields (be careful not to log secrets)
	for k, v := range config {
		// Skip sensitive fields
		if !isSensitiveField(k) {
			fields[k] = v
		}
	}

	l.WithFields(fields).Info("Service starting up")
}

func (l *Logger) LogShutdown(service string, reason string) {
	l.WithFields(Fields{
		"service": service,
		"reason":  reason,
		"type":    "shutdown",
	}).Info("Service shutting down")
}

// Helper function to check if field is sensitive
func isSensitiveField(field string) bool {
	sensitiveFields := []string{
		"password", "secret", "key", "token", "auth",
		"credential", "private", "sensitive",
	}

	fieldLower := strings.ToLower(field)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldLower, sensitive) {
			return true
		}
	}
	return false
}

// Middleware logging helpers
func (l *Logger) LogMiddleware(middleware, path string, duration int64, data map[string]interface{}) {
	fields := Fields{
		"middleware": middleware,
		"path":       path,
		"duration":   duration,
		"type":       "middleware",
	}

	// Add data fields
	for k, v := range data {
		fields[k] = v
	}

	l.WithFields(fields).Debug("Middleware executed")
}

// External service logging
func (l *Logger) LogExternalService(service, operation, method, url string, statusCode int, duration int64, err error) {
	entry := l.WithFields(Fields{
		"service":     service,
		"operation":   operation,
		"method":      method,
		"url":         url,
		"status_code": statusCode,
		"duration":    duration,
		"type":        "external_service",
	})

	if err != nil {
		entry.WithError(err).Error("External service call failed")
	} else {
		entry.Info("External service call completed")
	}
}