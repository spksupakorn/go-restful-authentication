package main

import (
	"log"

	"github.com/spksupakorn/go-restful-authentication/internal/config"
	"github.com/spksupakorn/go-restful-authentication/internal/http/server"
	"github.com/spksupakorn/go-restful-authentication/internal/infrastructure/database"
	"github.com/spksupakorn/go-restful-authentication/internal/utils/logger"
	"go.uber.org/zap"
)

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
