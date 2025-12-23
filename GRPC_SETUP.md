# gRPC Setup and Usage Guide

This document explains how to use the gRPC API alongside the REST API for user management.

## 📦 What's Included

### Protocol Buffer Definition
- **File**: `proto/user.proto`
- **Package**: `user`
- **Service**: `UserService`

### gRPC Server
- **Port**: `8081` (HTTP port + 1)
- **Authentication**: JWT token via metadata
- **Security**: Token-based interceptor

### Available Methods

#### Public Methods (No Authentication Required)
1. **CreateUser** - Register a new user
2. **Login** - Authenticate and get tokens

#### Protected Methods (JWT Required)
3. **GetUser** - Get user by ID
4. **ListUsers** - Get all users with pagination
5. **UpdateUser** - Update user details
6. **DeleteUser** - Delete a user

## 🚀 Quick Start

### 1. Install Required Tools

**Protocol Buffer Compiler** (if not already installed):
```bash
brew install protobuf
```

**Go Plugins**:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**grpcurl** (for testing):
```bash
brew install grpcurl
```

### 2. Generate gRPC Code

```bash
# Using Makefile
make proto

# Or manually
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/user.proto
```

### 3. Build and Run

```bash
# Build application (includes proto generation)
make build

# Run both HTTP and gRPC servers
make run
```

You should see:
```
HTTP Server: http://localhost:8080
gRPC Server: localhost:8081
Swagger UI: http://localhost:8080/swagger/index.html
```

## 🔑 Authentication

### For Public Methods (CreateUser, Login)

No authentication required. Just call the method directly.

### For Protected Methods (GetUser, ListUsers, UpdateUser, DeleteUser)

**Include JWT token in metadata:**

```bash
# Login first to get a token
grpcurl -plaintext \
  -d '{"email":"user@example.com","password":"password123"}' \
  localhost:8081 user.UserService/Login

# Use the access_token in subsequent requests
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"id":"USER_ID"}' \
  localhost:8081 user.UserService/GetUser
```

## 📝 Usage Examples

### 1. Create User (Register)

```bash
grpcurl -plaintext \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }' \
  localhost:8081 user.UserService/CreateUser
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

### 2. Login

```bash
grpcurl -plaintext \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }' \
  localhost:8081 user.UserService/Login
```

**Save the access_token from the response for authenticated calls.**

### 3. Get User by ID (Protected)

```bash
# Replace YOUR_ACCESS_TOKEN with your actual token
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"id":"6756c8a9f1234567890abcde"}' \
  localhost:8081 user.UserService/GetUser
```

### 4. List Users with Pagination (Protected)

```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10
  }' \
  localhost:8081 user.UserService/ListUsers
```

### 5. Update User (Protected)

```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "id": "6756c8a9f1234567890abcde",
    "name": "John Updated",
    "email": "john.updated@example.com"
  }' \
  localhost:8081 user.UserService/UpdateUser
```

### 6. Delete User (Protected)

```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"id":"6756c8a9f1234567890abcde"}' \
  localhost:8081 user.UserService/DeleteUser
```

## 🔍 Introspection

### List All Services

```bash
grpcurl -plaintext localhost:8081 list
```

**Output:**
```
grpc.reflection.v1alpha.ServerReflection
user.UserService
```

### List Methods in UserService

```bash
grpcurl -plaintext localhost:8081 list user.UserService
```

**Output:**
```
user.UserService.CreateUser
user.UserService.DeleteUser
user.UserService.GetUser
user.UserService.ListUsers
user.UserService.Login
user.UserService.UpdateUser
```

### Describe a Method

```bash
grpcurl -plaintext localhost:8081 describe user.UserService.CreateUser
```

## 🐍 Using gRPC with Python Client

### Install Python gRPC

```bash
pip install grpcio grpcio-tools
```

### Generate Python Client

```bash
python -m grpc_tools.protoc -I./proto \
  --python_out=./clients/python \
  --grpc_python_out=./clients/python \
  proto/user.proto
```

### Python Client Example

```python
import grpc
from proto import user_pb2
from proto import user_pb2_grpc

# Create channel
channel = grpc.insecure_channel('localhost:8081')
stub = user_pb2_grpc.UserServiceStub(channel)

# Login
login_request = user_pb2.LoginRequest(
    email='john@example.com',
    password='password123'
)
login_response = stub.Login(login_request)
access_token = login_response.access_token

# Get user with authentication
metadata = [('authorization', f'Bearer {access_token}')]
get_user_request = user_pb2.GetUserRequest(id='USER_ID')
user_response = stub.GetUser(get_user_request, metadata=metadata)

print(f"User: {user_response.user.name}")
```

## 🔧 Go Client Example

```go
package main

import (
    "context"
    "log"
    
    pb "github.com/spksupakorn/go-restful-authentication/proto/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/metadata"
)

func main() {
    // Connect to server
    conn, err := grpc.Dial("localhost:8081", 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewUserServiceClient(conn)
    ctx := context.Background()
    
    // Login
    loginResp, err := client.Login(ctx, &pb.LoginRequest{
        Email:    "john@example.com",
        Password: "password123",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    accessToken := loginResp.AccessToken
    
    // Create context with token metadata
    md := metadata.Pairs("authorization", "Bearer "+accessToken)
    authCtx := metadata.NewOutgoingContext(ctx, md)
    
    // Get user
    userResp, err := client.GetUser(authCtx, &pb.GetUserRequest{
        Id: "USER_ID",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("User: %s", userResp.User.Name)
}
```

## 📊 gRPC vs REST API

| Feature | REST API | gRPC API |
|---------|----------|----------|
| **Port** | 8080 | 8081 |
| **Protocol** | HTTP/1.1 JSON | HTTP/2 Protobuf |
| **Documentation** | Swagger UI | Proto file + Reflection |
| **Authentication** | Header: `Authorization: Bearer TOKEN` | Metadata: `authorization: Bearer TOKEN` |
| **Performance** | Good | Excellent (binary protocol) |
| **Browser Support** | Native | Requires grpc-web |

## 🔐 Security Features

### JWT Token Authentication
- **Interceptor**: `internal/grpc/interceptors/auth_interceptor.go`
- **Public Methods**: Login, CreateUser (no auth required)
- **Protected Methods**: All other methods require valid JWT in metadata

### Error Handling
- Custom errors are converted to gRPC status codes:
  - `400` → `InvalidArgument`
  - `401` → `Unauthenticated`
  - `403` → `PermissionDenied`
  - `404` → `NotFound`
  - `500` → `Internal`

### Metadata Format
```
Key: authorization
Value: Bearer <access_token>
```

## 🛠️ Development Workflow

### 1. Modify Proto File
Edit `proto/user.proto` to add/modify messages or methods.

### 2. Regenerate Code
```bash
make proto
```

### 3. Implement Server Methods
Update `internal/grpc/user_server.go` with new implementations.

### 4. Test with grpcurl
```bash
grpcurl -plaintext localhost:8081 list user.UserService
```

### 5. Update Clients
Regenerate client code for Python, Go, or other languages.

## 🐛 Troubleshooting

### "connection refused" error
**Solution**: Make sure the server is running on port 8081.
```bash
make run
```

### "Unauthenticated" error
**Solution**: Include JWT token in metadata:
```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_TOKEN" \
  ...
```

### "protoc: command not found"
**Solution**: Install Protocol Buffer compiler:
```bash
brew install protobuf
```

### "plugin not found" error
**Solution**: Install Go protoc plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Check if gRPC server is running
```bash
grpcurl -plaintext localhost:8081 list
```

## 📚 Resources

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [grpcurl Tool](https://github.com/fullstorydev/grpcurl)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)

## ✅ Testing Checklist

- [ ] Server starts on port 8081
- [ ] Can list services with grpcurl
- [ ] CreateUser works without authentication
- [ ] Login returns valid tokens
- [ ] GetUser requires authentication
- [ ] Invalid token returns "Unauthenticated"
- [ ] All CRUD operations work correctly

---

**Need help?** Check the logs or refer to the main README.md for general setup instructions.
