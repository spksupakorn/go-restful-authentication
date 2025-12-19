package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spksupakorn/go-restful-authentication/internal/http/controllers"
	"github.com/spksupakorn/go-restful-authentication/internal/http/middlewares"
	"github.com/spksupakorn/go-restful-authentication/internal/utils/jwt"
	"go.uber.org/zap"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(
	router *gin.Engine,
	userController *controllers.UserController,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
) {
	// Health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"service": "user-management-api",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
			auth.POST("/refresh", userController.RefreshToken)
		}

		// Protected routes (authentication required)
		users := v1.Group("/users")
		users.Use(middlewares.JWTAuthMiddleware(jwtManager, logger))
		{
			users.GET("", userController.GetAllUsers)
			users.GET("/:id", userController.GetUserByID)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
	}
}
