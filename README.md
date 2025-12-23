# RESTful API with Golang, MongoDB, JWT & Clean Architecture

A production-ready RESTful API built with Go that manages users with MongoDB persistence, JWT authentication, and follows hexagonal (clean) architecture principles.

## 🚀 Features

- ✅ **Clean Architecture** (Hexagonal Architecture / Ports & Adapters)
- ✅ **MongoDB** for data persistence
- ✅ **JWT Authentication** (Access & Refresh Tokens)
- ✅ **Password Hashing** with bcrypt
- ✅ **Input Validation** using go-playground/validator
- ✅ **Structured JSON Logging** with Zap
- ✅ **Graceful Shutdown** with context
- ✅ **Docker & Docker Compose** support
- ✅ **Background Goroutine** for monitoring user count
- ✅ **Logging Middleware** with execution time tracking
- ✅ **CORS** enabled
- ✅ **Unit Tests** with mocks
- ✅ **Makefile** for common tasks
- ✅ **Swagger/OpenAPI Documentation** with interactive UI
- ✅ **gRPC Support** with token-based authentication

## 📋 Requirements

- Go 1.23+
- MongoDB 7.0+
- Docker & Docker Compose (optional)
- Swag CLI (for Swagger generation)

## 🏗️ Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/                        # Configuration management
│   │   ├── config.go
│   │   └── config.dev.yaml
│   ├── domain/                        # Domain layer (entities & interfaces)
│   │   ├── entities/
│   │   │   └── userEntity.go
│   │   └── repositories/
│   │       ├── userRepository.go      # Repository interface (port)
│   │       └── mock_user_repository.go
│   ├── dto/                           # Data Transfer Objects
│   │   └── userDto.go
│   ├── http/                          # HTTP layer (adapters)
│   │   ├── controllers/
│   │   │   └── userController.go
│   │   ├── middlewares/
│   │   │   ├── jwtAuthMiddleware.go
│   │   │   └── loggingMiddleware.go
│   │   ├── routes/
│   │   │   └── routes.go
│   │   └── server/
│   │       └── server.go
│   ├── infrastructure/                # Infrastructure layer (adapters)
│   │   └── database/
│   │       ├── database.go
│   │       └── mongo_user_repository.go
│   ├── services/                      # Background services
│   │   └── background_service.go
│   ├── usecases/                      # Application/Business logic
│   │   ├── userUseCase.go
│   │   └── userUseCase_test.go
│   └── utils/                         # Utilities
│       ├── jwt/
│       │   └── jwt.go
│       ├── logger/
│       │   └── logger.go
│       ├── password/
│       │   └── password.go
│       └── validator/
│           └── validator.go
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── .env.example
├── go.mod
└── README.md
```

## 🔧 Installation

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/spksupakorn/go-restful-authentication.git
cd go-restful-authentication
```

2. **Install dependencies**
```bash
make install-deps
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

5. **Run MongoDB locally** (if not using Docker)
```bash
# Using Docker for MongoDB only
docker run -d -p 27017:27017 --name mongodb mongo:7.0
```

6. **Run the application**
```bash
make run
```

7. **Access Swagger UI**
```
Open your browser and navigate to: http://localhost:8080/swagger/index.html
```

### Using Docker Compose

```bash
# Build and start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down

# Clean up volumes
### Public Endpoints (No Authentication Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/swagger/index.html` | Interactive Swagger API documentation |
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login and get JWT tokens |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/health` | Health check ||
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login and get JWT tokens |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/health` | Health check |

### Protected Endpoints (JWT Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | Get all users (paginated) |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |

## 📝 API Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "user": {
    "id": "6756c8a9f1234567890abcde",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-12-19T10:30:00Z",
    "updated_at": "2024-12-19T10:30:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get All Users (Protected)

```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Get User by ID

```bash
curl -X GET http://localhost:8080/api/v1/users/6756c8a9f1234567890abcde \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/6756c8a9f1234567890abcde \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com"
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/6756c8a9f1234567890abcde \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🧪 Testing

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
# Opens coverage.html in your browser
```

### Run Specific Test
```bash
go test -v ./internal/usecases -run TestRegister_Success
```

## 🛠️ Development

### Available Make Commands

```bash
make help              # Show all available commands
make install-deps      # Install Go dependencies
make build             # Build the application
make run               # Run the application
make test              # Run tests
make test-coverage     # Run tests with coverage report
make clean             # Clean build artifacts
make docker-build      # Build Docker image
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make docker-logs       # View Docker logs
make docker-clean      # Remove containers and volumes
make lint              # Run golangci-lint
make fmt               # Format code
make vet               # Run go vet
```

## ⚙️ Configuration

Configuration can be set via environment variables or YAML config file.

### Environment Variables

```bash
# Application
APP_ENV=development
PORT=8080
HOST=0.0.0.0

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=user_management
MONGODB_TIMEOUT=10

# JWT
JWT_SECRET_KEY=your-secret-key
JWT_ACCESS_TOKEN_DURATION=24      # hours
JWT_REFRESH_TOKEN_DURATION=168    # hours (7 days)

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🏛️ Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters):

### Layers:

1. **Domain Layer** (`internal/domain`)
   - Entities: Core business objects
   - Repository Interfaces: Ports for data access

2. **Application Layer** (`internal/usecases`)
   - Business logic
   - Use cases implementation

3. **Infrastructure Layer** (`internal/infrastructure`)
   - Database implementation (MongoDB adapter)
   - External services

4. **HTTP Layer** (`internal/http`)
   - Controllers
   - Routes
   - Middlewares

### Benefits:
- ✅ Testability (easy to mock)
- ✅ Independence from frameworks
- ✅ Flexibility to change databases
- ✅ Clear separation of concerns
- ✅ Maintainable codebase

## 📊 Background Services

The application runs a background goroutine that logs the total number of users every 10 seconds. This demonstrates:
- Goroutine management
- Graceful shutdown of background tasks
- Context cancellation

## 📚 Swagger/OpenAPI Documentation

This API includes interactive **Swagger UI** documentation for easy testing and exploration.

### Accessing Swagger UI

Once the application is running, open your browser and navigate to:
```
http://localhost:8080/swagger/index.html
```

### Features:
- ✅ Interactive API explorer
- ✅ Try out endpoints directly from the browser
- ✅ View request/response schemas
- ✅ Automatic JWT token management
- ✅ Complete API documentation

### Using Swagger UI:

1. **For Public Endpoints** (Register, Login):
   - Click on the endpoint
   - Click "Try it out"
   - Fill in the request body
   - Click "Execute"

2. **For Protected Endpoints** (Users):
   - First, register or login to get an access token
   - Click the "Authorize" button at the top
   - Enter: `Bearer YOUR_ACCESS_TOKEN`
   - Click "Authorize"
   - Now you can access protected endpoints

### Regenerating Documentation

If you modify API comments in the code:

```bash
# Regenerate Swagger docs
make swagger

# Or manually
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### API Documentation Comments

The Swagger docs are generated from special comments in the code:

```go
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
    // ... implementation
}
```

For more details, see [SWAGGER_SETUP.md](SWAGGER_SETUP.md).

## 🔌 gRPC Support

This API also supports **gRPC** for high-performance, language-agnostic communication. Both HTTP REST and gRPC run concurrently.

### gRPC Endpoints

The application runs two servers:
- **HTTP REST**: `localhost:8080`
- **gRPC**: `localhost:8081` (HTTP port + 1)

### Available gRPC Methods

```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);  // Public
  rpc GetUser(GetUserRequest) returns (GetUserResponse);          // Protected
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);    // Protected
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse); // Protected
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse); // Protected
  rpc Login(LoginRequest) returns (LoginResponse);                // Public
}
```

### Testing gRPC with grpcurl

**Install grpcurl:**
```bash
# macOS
brew install grpcurl

# Linux
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

**Example gRPC Calls:**

1. **Login** (get tokens):
```bash
grpcurl -plaintext \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }' \
  localhost:8081 user.UserService/Login
```

2. **Create User** (register):
```bash
grpcurl -plaintext \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "password123"
  }' \
  localhost:8081 user.UserService/CreateUser
```

3. **Get User** (with authentication):
```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"id": "USER_ID"}' \
  localhost:8081 user.UserService/GetUser
```

4. **List Users** (with pagination):
```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"page": 1, "page_size": 10}' \
  localhost:8081 user.UserService/ListUsers
```

### Authentication

gRPC uses **metadata** for JWT authentication:
- Public methods: `CreateUser`, `Login`
- Protected methods: `GetUser`, `ListUsers`, `UpdateUser`, `DeleteUser`

For protected methods, include the token in metadata:
```bash
-H "authorization: Bearer YOUR_TOKEN"
```

The gRPC auth interceptor validates tokens and adds user context.

### Regenerating Proto Files

If you modify `proto/user.proto`:

```bash
# Generate Go code from proto
make proto-gen

# Or manually
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

For complete gRPC setup instructions, see [GRPC_SETUP.md](GRPC_SETUP.md).

## 🔐 Security Features

- **Password Hashing**: bcrypt with cost factor 10
- **JWT Authentication**: HMAC-SHA256
- **Token Expiration**: Access tokens expire in 24 hours
- **Input Validation**: All inputs are validated
- **CORS**: Configured for cross-origin requests

## 📈 Production Considerations

### Before deploying to production:

1. **Change JWT Secret**: Use a strong, random secret key
2. **Use HTTPS**: Enable TLS/SSL
3. **Environment Variables**: Use secure secret management
4. **Database**: Use MongoDB replica set
5. **Logging**: Ensure JSON logging is enabled
6. **Monitoring**: Add application monitoring (Prometheus, etc.)
7. **Rate Limiting**: Implement rate limiting middleware
8. **CORS**: Restrict allowed origins

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👤 Author

**Supakorn Wongsangiam**

- GitHub: [@spksupakorn](https://github.com/spksupakorn)

## 🙏 Acknowledgments

- Clean Architecture by Robert C. Martin
- Go community for excellent libraries
- MongoDB team for the official Go driver

---

**Happy Coding! 🚀**
