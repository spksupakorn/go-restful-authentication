# Quick Reference - gRPC API

## 🚀 Server Ports
```
HTTP REST API: http://localhost:8080
gRPC API:      localhost:8081
Swagger UI:    http://localhost:8080/swagger/index.html
```

## 📝 Quick Commands

### Start Server
```bash
make run
# or
./bin/api
```

### Generate Proto
```bash
make proto-gen
```

### Build
```bash
make build
```

## 🔑 Authentication

### Get Token (Login)
```bash
grpcurl -plaintext \
  -d '{"email":"user@example.com","password":"pass123"}' \
  localhost:8081 user.UserService/Login
```

### Use Token
```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_TOKEN" \
  -d '{"page":1,"page_size":10}' \
  localhost:8081 user.UserService/ListUsers
```

## 📋 Available Methods

| Method | Auth Required | Description |
|--------|--------------|-------------|
| `CreateUser` | ❌ No | Register new user |
| `Login` | ❌ No | Get JWT tokens |
| `GetUser` | ✅ Yes | Get user by ID |
| `ListUsers` | ✅ Yes | Get all users (paginated) |
| `UpdateUser` | ✅ Yes | Update user details |
| `DeleteUser` | ✅ Yes | Delete user |

## 🧪 Test Commands

### List Services
```bash
grpcurl -plaintext localhost:8081 list
```

### Describe Service
```bash
grpcurl -plaintext localhost:8081 describe user.UserService
```

### Create User
```bash
grpcurl -plaintext \
  -d '{
    "name":"John Doe",
    "email":"john@example.com",
    "password":"password123"
  }' \
  localhost:8081 user.UserService/CreateUser
```

### Get User
```bash
grpcurl -plaintext \
  -H "authorization: Bearer TOKEN" \
  -d '{"id":"USER_ID"}' \
  localhost:8081 user.UserService/GetUser
```

### Update User
```bash
grpcurl -plaintext \
  -H "authorization: Bearer TOKEN" \
  -d '{
    "id":"USER_ID",
    "name":"Jane Doe",
    "email":"jane@example.com"
  }' \
  localhost:8081 user.UserService/UpdateUser
```

### Delete User
```bash
grpcurl -plaintext \
  -H "authorization: Bearer TOKEN" \
  -d '{"id":"USER_ID"}' \
  localhost:8081 user.UserService/DeleteUser
```

## 🐛 Troubleshooting

### Connection Refused
```bash
# Check if server is running
curl http://localhost:8080/health

# Check logs
tail -f logs/app.log
```

### Invalid Token
```bash
# Token format must be:
-H "authorization: Bearer YOUR_TOKEN"

# NOT:
-H "Authorization: YOUR_TOKEN"
```

### Proto Changes Not Applied
```bash
# Regenerate proto files
make proto-gen

# Rebuild
make build
```

## 📚 Documentation

- Full Setup: [GRPC_SETUP.md](GRPC_SETUP.md)
- Implementation Details: [GRPC_IMPLEMENTATION_SUMMARY.md](GRPC_IMPLEMENTATION_SUMMARY.md)
- Main README: [README.md](README.md)

---

**Quick Tip**: Use the `-v` flag with grpcurl for verbose output to see request/response headers!
