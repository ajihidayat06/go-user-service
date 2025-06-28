package middleware

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"go-user-service/internal/pkg/errors"
	"go-user-service/internal/pkg/logger"
	"go-user-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(l logger.Logger) gin.HandlerFunc {
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

		l.Info(fmt.Sprintf("%s %s %d %v %s",
			method,
			path,
			statusCode,
			latency,
			clientIP,
		))
	}
}

// Recovery middleware untuk menangkap panic
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v\n%s", recovered, debug.Stack())

		appErr := errors.New(errors.ErrCodeInternal, "Internal server error")
		response.Error(c, appErr)
	})
}

// ErrorLogger middleware untuk log error
func ErrorLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log jika status code >= 400
		if param.StatusCode >= 400 {
			log.Printf("Error response: %d %s %s",
				param.StatusCode, param.Method, param.Path)
		}

		return ""
	})
}
