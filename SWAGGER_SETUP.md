# Swagger/OpenAPI Setup Guide

This document explains how Swagger is integrated into this Go REST API project.

## 📦 What's Installed

1. **Swagger Dependencies**:
   - `github.com/swaggo/swag` - CLI tool for generating Swagger docs
   - `github.com/swaggo/gin-swagger` - Gin middleware for Swagger
   - `github.com/swaggo/files` - Embedded Swagger UI files

2. **Generated Documentation**:
   - `docs/swagger.json` - OpenAPI spec in JSON format
   - `docs/swagger.yaml` - OpenAPI spec in YAML format
   - `docs/docs.go` - Go file with embedded docs

## 🚀 Quick Start

### 1. Install Swag CLI (First Time Only)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Generate Documentation

Every time you update API comments:

```bash
# Using Makefile
make swagger

# Or manually
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### 3. Build and Run

```bash
# Build (automatically generates Swagger docs)
make build

# Run
make run
```

### 4. Access Swagger UI

Open your browser:
```
http://localhost:8080/swagger/index.html
```

## 📝 How to Document Your APIs

### General API Information (in `cmd/api/main.go`)

```go
// @title User Management API
// @version 1.0
// @description RESTful API for user management with JWT authentication
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Endpoint Documentation (in controllers)

```go
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
    // ... implementation
}
```

### Protected Endpoints (requiring JWT)

Add `@Security BearerAuth`:

```go
// @Summary Get all users
// @Security BearerAuth
// @Tags users
// @Param page query int false "Page number"
// @Success 200 {object} dto.UsersListResponse
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
    // ... implementation
}
```

## 🔑 Testing with Swagger UI

### For Public Endpoints (No Auth Required)

1. Navigate to `/swagger/index.html`
2. Find the endpoint (e.g., "POST /auth/register")
3. Click "Try it out"
4. Fill in the request body
5. Click "Execute"
6. View the response

### For Protected Endpoints (Auth Required)

1. First, get a token by calling `/auth/login` or `/auth/register`
2. Copy the `access_token` from the response
3. Click the **"Authorize"** button at the top
4. Enter: `Bearer YOUR_ACCESS_TOKEN_HERE`
5. Click "Authorize"
6. Now you can call protected endpoints

## 📋 Swagger Annotation Reference

### Common Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Summary` | Short description | `@Summary Get user by ID` |
| `@Description` | Detailed description | `@Description Retrieves user details by ID` |
| `@Tags` | Group endpoints | `@Tags users` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Param` | Parameter definition | `@Param id path string true "User ID"` |
| `@Success` | Success response | `@Success 200 {object} dto.UserResponse` |
| `@Failure` | Error response | `@Failure 404 {object} dto.ErrorResponse` |
| `@Router` | Route path and method | `@Router /users/{id} [get]` |
| `@Security` | Security requirement | `@Security BearerAuth` |

### Parameter Types

```go
// Path parameter
@Param id path string true "User ID"

// Query parameter
@Param page query int false "Page number" default(1)

// Header parameter
@Param Authorization header string true "Bearer token"

// Body parameter
@Param request body dto.RegisterRequest true "User data"
```

## 🔄 Workflow

### When Adding a New Endpoint

1. **Write the handler** in your controller
2. **Add Swagger comments** above the handler function
3. **Define DTOs** if needed (request/response structs)
4. **Regenerate docs**: `make swagger`
5. **Build and test**: `make build && make run`
6. **Verify in Swagger UI**: http://localhost:8080/swagger/index.html

### When Modifying an Existing Endpoint

1. **Update the handler** logic
2. **Update Swagger comments** if signature changed
3. **Update DTOs** if request/response changed
4. **Regenerate docs**: `make swagger`
5. **Test the changes** in Swagger UI

## 🛠️ Troubleshooting

### "docs package not found"

**Solution**: Run `make swagger` to generate the docs package.

### Swagger UI shows old documentation

**Solution**: Regenerate docs and restart the server:
```bash
make swagger
make run
```

### Can't access protected endpoints in Swagger

**Solution**: 
1. Call `/auth/login` or `/auth/register` first
2. Copy the access token
3. Click "Authorize" button
4. Enter: `Bearer <your-token>`
5. Click "Authorize"

### Changes not reflected in Swagger UI

**Solution**:
1. Make sure you regenerated docs: `make swagger`
2. Restart the application
3. Hard refresh browser (Ctrl+Shift+R or Cmd+Shift+R)

## 📚 Additional Resources

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Gin Swagger Examples](https://github.com/swaggo/gin-swagger)

## ✅ Checklist

- [x] Swagger dependencies installed
- [x] Documentation generated in `docs/` folder
- [x] Swagger route added (`/swagger/*any`)
- [x] API metadata defined in `main.go`
- [x] All endpoints documented with comments
- [x] Request/Response DTOs properly tagged
- [x] JWT security definition configured
- [x] Makefile includes `swagger` target

## 🎯 Best Practices

1. **Keep docs up-to-date**: Regenerate after every API change
2. **Use clear descriptions**: Help API consumers understand endpoints
3. **Document all parameters**: Include type, required/optional, defaults
4. **Include examples**: Add example values in DTOs
5. **Group endpoints**: Use tags to organize related endpoints
6. **Document errors**: List all possible error responses
7. **Version your API**: Update version in main.go when making breaking changes

---

**Need help?** Check the existing controller comments for examples or refer to the Swag documentation.
