package main

import (
	"log"

	"github.com/spksupakorn/go-restful-authentication/internal/config"
	"github.com/spksupakorn/go-restful-authentication/internal/http/server"
	"github.com/spksupakorn/go-restful-authentication/internal/infrastructure/database"
	"github.com/spksupakorn/go-restful-authentication/internal/pkg/logger"

	"go.uber.org/zap"

	_ "github.com/spksupakorn/go-restful-authentication/docs" // Swagger docs
)

// @title User Management API
// @version 1.0
// @description RESTful API for user management with JWT authentication and MongoDB
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	zapLogger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting User Management API",
		zap.String("version", "1.0.0"),
		zap.String("environment", cfg.Server.Env),
	)

	// Connect to MongoDB
	db, err := database.NewMongoDB(cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	// Create and start server
	srv := server.NewServer(cfg, zapLogger, db)
	if err := srv.Start(); err != nil {
		zapLogger.Fatal("Server error", zap.Error(err))
	}
}
