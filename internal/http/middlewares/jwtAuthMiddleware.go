package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spksupakorn/go-restful-authentication/internal/dto"
	"github.com/spksupakorn/go-restful-authentication/internal/utils/jwt"
	"go.uber.org/zap"
)

// JWTAuthMiddleware creates a middleware for JWT authentication
func JWTAuthMiddleware(jwtManager *jwt.JWTManager, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header")
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		// Check if header starts with "Bearer "
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format")
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid authorization header format. Expected: Bearer <token>",
			})
			ctx.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			logger.Warn("Invalid or expired token", zap.Error(err))
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		// Set user info in context
		ctx.Set("user_id", claims.UserID.Hex())
		ctx.Set("email", claims.Email)

		ctx.Next()
	}
}
