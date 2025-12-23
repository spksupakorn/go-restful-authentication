package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spksupakorn/go-restful-authentication/internal/config"
	grpcServer "github.com/spksupakorn/go-restful-authentication/internal/grpc"
	"github.com/spksupakorn/go-restful-authentication/internal/grpc/interceptors"
	"github.com/spksupakorn/go-restful-authentication/internal/pkg/jwt"
	"github.com/spksupakorn/go-restful-authentication/internal/pkg/validator"
	"github.com/spksupakorn/go-restful-authentication/internal/repositories"
	pb "github.com/spksupakorn/go-restful-authentication/proto"

	"github.com/spksupakorn/go-restful-authentication/internal/http/controllers"
	"github.com/spksupakorn/go-restful-authentication/internal/http/middlewares"
	"github.com/spksupakorn/go-restful-authentication/internal/http/routes"
	"github.com/spksupakorn/go-restful-authentication/internal/infrastructure/database"
	"github.com/spksupakorn/go-restful-authentication/internal/services"
	"github.com/spksupakorn/go-restful-authentication/internal/usecases"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the HTTP and gRPC server
type Server struct {
	config            *config.Config
	logger            *zap.Logger
	router            *gin.Engine
	db                *database.MongoDB
	backgroundService *services.BackgroundService
	grpcServer        *grpc.Server
}

var (
	once           sync.Once
	serverInstance *Server
)

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

	once.Do(func() {
		serverInstance = &Server{
			config: cfg,
			logger: logger,
			router: router,
			db:     db,
		}
	})

	return serverInstance
}

// Start starts both HTTP and gRPC servers
func (s *Server) Start() error {
	// Initialize dependencies
	userRepo := repositories.NewMongoUserRepository(s.db, s.logger)
	jwtManager := jwt.NewJWTManager(s.config)
	validatorInstance := validator.NewValidator()

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo, jwtManager, s.logger)

	// Initialize HTTP controllers
	userController := controllers.NewUserController(userUseCase, validatorInstance, s.logger)

	// Setup HTTP routes
	routes.SetupRoutes(s.router, userController, jwtManager, s.logger)

	// Initialize gRPC server with interceptors
	authInterceptor := interceptors.NewAuthInterceptor(jwtManager, s.logger)
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	// Register gRPC services
	userGRPCServer := grpcServer.NewUserServer(userUseCase, s.logger)
	pb.RegisterUserServiceServer(s.grpcServer, userGRPCServer)

	// Enable gRPC reflection for tools like grpcurl
	reflection.Register(s.grpcServer)

	// Start background service
	s.backgroundService = services.NewBackgroundService(userRepo, s.logger)
	ctx := context.Background()
	s.backgroundService.Start(ctx)

	// Start gRPC server in a goroutine
	grpcPort := s.config.Server.Port + 1 // Use HTTP port + 1 for gRPC
	grpcAddr := fmt.Sprintf("%s:%d", s.config.Server.Host, grpcPort)

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		s.logger.Fatal("Failed to listen for gRPC", zap.Error(err))
	}

	go func() {
		s.logger.Info("Starting gRPC server",
			zap.String("address", grpcAddr),
			zap.String("environment", s.config.Server.Env),
		)

		if err := s.grpcServer.Serve(grpcListener); err != nil {
			s.logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Create HTTP server
	httpAddr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	srv := &http.Server{
		Addr:    httpAddr,
		Handler: s.router,
	}

	// Start HTTP server in a goroutine
	go func() {
		s.logger.Info("Starting HTTP server",
			zap.String("address", httpAddr),
			zap.String("environment", s.config.Server.Env),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down servers...")

	// Stop background service
	if s.backgroundService != nil {
		s.backgroundService.Stop()
	}

	// Graceful shutdown of gRPC server
	s.logger.Info("Stopping gRPC server...")
	s.grpcServer.GracefulStop()

	// Graceful shutdown of HTTP server with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("HTTP server forced to shutdown", zap.Error(err))
		return err
	}

	// Close database connection
	if err := s.db.Close(shutdownCtx); err != nil {
		s.logger.Error("Failed to close database connection", zap.Error(err))
		return err
	}

	s.logger.Info("Servers exited gracefully")
	return nil
}
