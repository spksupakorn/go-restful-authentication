package grpc

import (
	"context"

	"github.com/spksupakorn/go-restful-authentication/internal/dto"
	"github.com/spksupakorn/go-restful-authentication/internal/grpc/interceptors"
	"github.com/spksupakorn/go-restful-authentication/internal/usecases"
	pb "github.com/spksupakorn/go-restful-authentication/proto"
	"go.uber.org/zap"
)

// UserServer implements the gRPC UserServiceServer interface
type UserServer struct {
	pb.UnimplementedUserServiceServer
	userUseCase usecases.UserUseCase
	logger      *zap.Logger
}

// NewUserServer creates a new gRPC user server
func NewUserServer(userUseCase usecases.UserUseCase, logger *zap.Logger) *UserServer {
	return &UserServer{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// CreateUser creates a new user via gRPC
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.logger.Info("gRPC CreateUser called",
		zap.String("email", req.GetEmail()),
		zap.String("name", req.GetName()),
	)

	// Create register DTO
	registerDto := &dto.RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Name:     req.GetName(),
	}

	// Call use case
	authResp, err := s.userUseCase.Register(ctx, registerDto)
	if err != nil {
		s.logger.Error("Failed to create user via gRPC",
			zap.Error(err),
		)
		return nil, interceptors.ErrorToGRPCStatus(err)
	}

	// Convert to gRPC response
	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:        authResp.User.ID,
			Email:     authResp.User.Email,
			Name:      authResp.User.Name,
			CreatedAt: authResp.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: authResp.User.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
	}, nil
}

// GetUser retrieves a user by ID via gRPC
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.logger.Info("gRPC GetUser called",
		zap.String("user_id", req.GetId()),
	)

	// Call use case
	user, err := s.userUseCase.GetUserByID(ctx, req.GetId())
	if err != nil {
		s.logger.Error("Failed to get user via gRPC",
			zap.String("user_id", req.GetId()),
			zap.Error(err),
		)
		return nil, interceptors.ErrorToGRPCStatus(err)
	}

	// Convert to gRPC response
	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// Login authenticates a user and returns tokens via gRPC
func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Info("gRPC Login called",
		zap.String("email", req.GetEmail()),
	)

	// Create login DTO
	loginDto := &dto.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	// Call use case
	authResp, err := s.userUseCase.Login(ctx, loginDto)
	if err != nil {
		s.logger.Error("Failed to login via gRPC",
			zap.String("email", req.GetEmail()),
			zap.Error(err),
		)
		return nil, interceptors.ErrorToGRPCStatus(err)
	}

	// Convert to gRPC response
	return &pb.LoginResponse{
		User: &pb.User{
			Id:        authResp.User.ID,
			Email:     authResp.User.Email,
			Name:      authResp.User.Name,
			CreatedAt: authResp.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: authResp.User.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
	}, nil
}

// ListUsers retrieves all users with pagination via gRPC
func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.logger.Info("gRPC ListUsers called",
		zap.Int32("page", req.GetPage()),
		zap.Int32("page_size", req.GetPageSize()),
	)

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Call use case
	usersResp, err := s.userUseCase.GetAllUsers(ctx, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list users via gRPC",
			zap.Error(err),
		)
		return nil, interceptors.ErrorToGRPCStatus(err)
	}

	// Convert to gRPC response
	users := make([]*pb.User, len(usersResp.Users))
	for i, user := range usersResp.Users {
		users[i] = &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &pb.ListUsersResponse{
		Users:      users,
		Total:      int64(usersResp.Total),
		Page:       int32(usersResp.Page),
		PageSize:   int32(usersResp.PageSize),
		TotalPages: int32(usersResp.TotalPages),
	}, nil
}

// UpdateUser updates a user via gRPC
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	s.logger.Info("gRPC UpdateUser called",
		zap.String("user_id", req.GetId()),
	)

	// Create update DTO
	updateDto := &dto.UpdateUserRequest{
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}

	// Call use case
	user, err := s.userUseCase.UpdateUser(ctx, req.GetId(), updateDto)
	if err != nil {
		s.logger.Error("Failed to update user via gRPC",
			zap.String("user_id", req.GetId()),
			zap.Error(err),
		)
		return nil, interceptors.ErrorToGRPCStatus(err)
	}

	// Convert to gRPC response
	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// DeleteUser deletes a user via gRPC
func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.logger.Info("gRPC DeleteUser called",
		zap.String("user_id", req.GetId()),
	)

	// Call use case
	err := s.userUseCase.DeleteUser(ctx, req.GetId())
	if err != nil {
		s.logger.Error("Failed to delete user via gRPC",
			zap.String("user_id", req.GetId()),
			zap.Error(err),
		)
		return nil, interceptors.ErrorToGRPCStatus(err)
	}

	// Return success response
	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}
