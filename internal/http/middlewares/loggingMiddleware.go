package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware creates a middleware for logging HTTP requests
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		// Process request
		ctx.Next()

		// Calculate execution time
		latency := time.Since(start)
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()

		// Log request
		logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", ctx.Request.UserAgent()),
		)
	}
}
