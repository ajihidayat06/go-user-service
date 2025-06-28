package middleware

import "github.com/gin-gonic/gin"

// RequestID middleware untuk tracking request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-ID")
		if requestId == "" {
			requestId = generateRequestID()
		}

		c.Header("X-Request-ID", requestId)
		c.Set("request_id", requestId)
		c.Next()
	}
}

func generateRequestID() string {
	// Simple implementation, bisa pakai UUID library
	return "req-" + randomString(10)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[len(charset)%len(charset)]
	}
	return string(b)
}
