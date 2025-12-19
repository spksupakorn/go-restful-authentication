package usecases

import (
	"context"
	"fmt"

	"github.com/spksupakorn/go-restful-authentication/internal/domain/entities"
	"github.com/spksupakorn/go-restful-authentication/internal/domain/repositories"
	"github.com/spksupakorn/go-restful-authentication/internal/dto"
	"github.com/spksupakorn/go-restful-authentication/internal/utils/jwt"
	"github.com/spksupakorn/go-restful-authentication/internal/utils/password"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context, page, pageSize int) (*dto.UsersListResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type userUseCase struct {
	userRepo   repositories.UserRepository
	jwtManager *jwt.JWTManager
	logger     *zap.Logger
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo repositories.UserRepository,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
) UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Register registers a new user
func (u *userUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	exists, err := u.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		u.logger.Error("Failed to check email existence", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		u.logger.Error("Failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to process password")
	}

	// Create user entity
	user := &entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save user to database
	if err := u.userRepo.Create(ctx, user); err != nil {
		u.logger.Error("Failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user")
	}

	// Generate tokens
	tokens, err := u.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		u.logger.Error("Failed to generate tokens", zap.Error(err))
		return nil, fmt.Errorf("failed to generate tokens")
	}

	u.logger.Info("User registered successfully",
		zap.String("user_id", user.ID.Hex()),
		zap.String("email", user.Email),
	)

	return &dto.AuthResponse{
		User:         dto.ConvertToUserResponse(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// Login authenticates a user
func (u *userUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Warn("Login attempt with non-existent email", zap.String("email", req.Email))
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare password
	if err := password.ComparePassword(user.Password, req.Password); err != nil {
		u.logger.Warn("Login attempt with incorrect password",
			zap.String("user_id", user.ID.Hex()),
			zap.String("email", user.Email),
		)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	tokens, err := u.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		u.logger.Error("Failed to generate tokens", zap.Error(err))
		return nil, fmt.Errorf("failed to generate tokens")
	}

	u.logger.Info("User logged in successfully",
		zap.String("user_id", user.ID.Hex()),
		zap.String("email", user.Email),
	)

	return &dto.AuthResponse{
		User:         dto.ConvertToUserResponse(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// GetUserByID retrieves a user by ID
func (u *userUseCase) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := u.userRepo.GetByID(ctx, objectID)
	if err != nil {
		u.logger.Error("Failed to get user", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("user not found")
	}

	response := dto.ConvertToUserResponse(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return &response, nil
}

// GetAllUsers retrieves all users with pagination
func (u *userUseCase) GetAllUsers(ctx context.Context, page, pageSize int) (*dto.UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	users, err := u.userRepo.GetAll(ctx, skip, limit)
	if err != nil {
		u.logger.Error("Failed to get all users", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve users")
	}

	total, err := u.userRepo.Count(ctx)
	if err != nil {
		u.logger.Error("Failed to count users", zap.Error(err))
		return nil, fmt.Errorf("failed to count users")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.ConvertToUserResponse(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &dto.UsersListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser updates a user
func (u *userUseCase) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get existing user
	user, err := u.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email already exists
		if req.Email != user.Email {
			exists, err := u.userRepo.EmailExists(ctx, req.Email)
			if err != nil {
				u.logger.Error("Failed to check email existence", zap.Error(err))
				return nil, fmt.Errorf("internal server error")
			}
			if exists {
				return nil, fmt.Errorf("email already exists")
			}
			user.Email = req.Email
		}
	}

	// Update user
	if err := u.userRepo.Update(ctx, user); err != nil {
		u.logger.Error("Failed to update user", zap.Error(err))
		return nil, fmt.Errorf("failed to update user")
	}

	u.logger.Info("User updated successfully", zap.String("user_id", user.ID.Hex()))

	response := dto.ConvertToUserResponse(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return &response, nil
}

// DeleteUser deletes a user
func (u *userUseCase) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	if err := u.userRepo.Delete(ctx, objectID); err != nil {
		u.logger.Error("Failed to delete user", zap.Error(err))
		return fmt.Errorf("failed to delete user")
	}

	u.logger.Info("User deleted successfully", zap.String("user_id", id))
	return nil
}

// RefreshToken generates a new access token from refresh token
func (u *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	accessToken, err := u.jwtManager.RefreshAccessToken(refreshToken)
	if err != nil {
		u.logger.Warn("Failed to refresh token", zap.Error(err))
		return "", fmt.Errorf("invalid refresh token")
	}

	return accessToken, nil
}
