package services

import (
	"context"
	"time"

	"github.com/spksupakorn/go-restful-authentication/internal/domain/repositories"
	"go.uber.org/zap"
)

// BackgroundService handles background tasks
type BackgroundService struct {
	userRepo repositories.UserRepository
	logger   *zap.Logger
	ticker   *time.Ticker
	done     chan bool
}

// NewBackgroundService creates a new background service
func NewBackgroundService(userRepo repositories.UserRepository, logger *zap.Logger) *BackgroundService {
	return &BackgroundService{
		userRepo: userRepo,
		logger:   logger,
		ticker:   time.NewTicker(10 * time.Second),
		done:     make(chan bool),
	}
}

// Start starts the background service
func (s *BackgroundService) Start(ctx context.Context) {
	s.logger.Info("Starting background service - user count logger")

	go func() {
		// Log immediately on start
		s.logUserCount(ctx)

		for {
			select {
			case <-s.done:
				s.logger.Info("Stopping background service")
				return
			case <-s.ticker.C:
				s.logUserCount(ctx)
			case <-ctx.Done():
				s.logger.Info("Context cancelled, stopping background service")
				return
			}
		}
	}()
}

// Stop stops the background service
func (s *BackgroundService) Stop() {
	s.logger.Info("Received stop signal for background service")
	s.ticker.Stop()
	close(s.done)
}

// logUserCount logs the current number of users in the database
func (s *BackgroundService) logUserCount(ctx context.Context) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		s.logger.Error("Failed to count users in background service", zap.Error(err))
		return
	}

	s.logger.Info("User count status",
		zap.Int64("total_users", count),
		zap.Time("timestamp", time.Now()),
	)
}
