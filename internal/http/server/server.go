package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spksupakorn/go-restful-authentication/internal/config"
	"github.com/spksupakorn/go-restful-authentication/internal/repositories"

	"github.com/spksupakorn/go-restful-authentication/internal/http/controllers"
	"github.com/spksupakorn/go-restful-authentication/internal/http/middlewares"
	"github.com/spksupakorn/go-restful-authentication/internal/http/routes"
	"github.com/spksupakorn/go-restful-authentication/internal/infrastructure/database"
	"github.com/spksupakorn/go-restful-authentication/internal/services"
	"github.com/spksupakorn/go-restful-authentication/internal/usecases"
	"github.com/spksupakorn/go-restful-authentication/internal/utils/jwt"
	"github.com/spksupakorn/go-restful-authentication/internal/utils/validator"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	config            *config.Config
	logger            *zap.Logger
	router            *gin.Engine
	db                *database.MongoDB
	backgroundService *services.BackgroundService
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, logger *zap.Logger, db *database.MongoDB) *Server {
	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Recovery middleware
	router.Use(gin.Recovery())

	// Logging middleware
	router.Use(middlewares.LoggingMiddleware(logger))

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Server{
		config: cfg,
		logger: logger,
		router: router,
		db:     db,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Initialize dependencies
	userRepo := repositories.NewMongoUserRepository(s.db, s.logger)
	jwtManager := jwt.NewJWTManager(s.config)
	validatorInstance := validator.NewValidator()

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo, jwtManager, s.logger)

	// Initialize controllers
	userController := controllers.NewUserController(userUseCase, validatorInstance, s.logger)

	// Setup routes
	routes.SetupRoutes(s.router, userController, jwtManager, s.logger)

	// Start background service
	s.backgroundService = services.NewBackgroundService(userRepo, s.logger)
	ctx := context.Background()
	s.backgroundService.Start(ctx)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		s.logger.Info("Starting HTTP server",
			zap.String("address", addr),
			zap.String("environment", s.config.Server.Env),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down server...")

	// Stop background service
	if s.backgroundService != nil {
		s.backgroundService.Stop()
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	// Close database connection
	if err := s.db.Close(ctx); err != nil {
		s.logger.Error("Failed to close database connection", zap.Error(err))
		return err
	}

	s.logger.Info("Server exited gracefully")
	return nil
}
