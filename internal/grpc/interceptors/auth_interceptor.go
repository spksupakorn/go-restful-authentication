package interceptors

import (
	"context"
	"strings"

	"github.com/spksupakorn/go-restful-authentication/internal/http/custom"
	"github.com/spksupakorn/go-restful-authentication/internal/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a gRPC interceptor for JWT authentication
type AuthInterceptor struct {
	jwtManager *jwt.JWTManager
	logger     *zap.Logger
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(jwtManager *jwt.JWTManager, logger *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Unary returns a server interceptor function for unary calls
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		i.logger.Debug("gRPC unary call",
			zap.String("method", info.FullMethod),
		)

		// Skip authentication for login and create user methods
		if i.isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract and validate token
		newCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// Stream returns a server interceptor function for stream calls
func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		i.logger.Debug("gRPC stream call",
			zap.String("method", info.FullMethod),
		)

		// Skip authentication for public methods
		if i.isPublicMethod(info.FullMethod) {
			return handler(srv, stream)
		}

		// Extract and validate token
		newCtx, err := i.authenticate(stream.Context())
		if err != nil {
			return err
		}

		wrapped := &wrappedStream{
			ServerStream: stream,
			ctx:          newCtx,
		}

		return handler(srv, wrapped)
	}
}

// authenticate extracts and validates JWT token from metadata
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		i.logger.Warn("Missing metadata in gRPC request")
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md["authorization"]
	if len(values) == 0 {
		i.logger.Warn("Missing authorization token in metadata")
		return nil, status.Error(codes.Unauthenticated, "missing authorization token")
	}

	accessToken := values[0]

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(accessToken, "Bearer ") {
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	}

	// Validate token
	claims, err := i.jwtManager.ValidateToken(accessToken)
	if err != nil {
		i.logger.Warn("Invalid access token",
			zap.Error(err),
		)
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	} // Add user info to context
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "user_email", claims.Email)

	i.logger.Debug("Successfully authenticated gRPC request",
		zap.String("user_id", claims.UserID.Hex()),
		zap.String("email", claims.Email),
	)

	return ctx, nil
}

// isPublicMethod checks if the method doesn't require authentication
func (i *AuthInterceptor) isPublicMethod(method string) bool {
	publicMethods := []string{
		"/user.UserService/Login",
		"/user.UserService/CreateUser",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}

	return false
}

// wrappedStream wraps a grpc.ServerStream with a custom context
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the custom context
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// ErrorToGRPCStatus converts custom errors to gRPC status
func ErrorToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*custom.AppError); ok {
		switch appErr.Code {
		case 400:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case 401:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case 403:
			return status.Error(codes.PermissionDenied, appErr.Message)
		case 404:
			return status.Error(codes.NotFound, appErr.Message)
		case 422:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case 500:
			return status.Error(codes.Internal, appErr.Message)
		default:
			return status.Error(codes.Unknown, appErr.Message)
		}
	}

	return status.Error(codes.Internal, err.Error())
}
