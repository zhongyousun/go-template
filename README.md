# Go Template

A modern, scalable Go web project template featuring:

- **JWT authentication** and role-based access control
- **User and Order management**
- **Layered architecture** (handler, service, repository)
- **Redis caching** for fast user lookup
- **Swagger API documentation**
- **Extensible, clean project structure**

---

## Project Structure

```
go.mod
go.sum
cmd/
    server/
        main.go
internal/
    user/
        handler/
        repository/
        service/
    order/
        handler/
        repository/
        service/
    middleware/
    common/
        commonmodel/
    db/
pkg/
    redisclient/
configs/
docs/
test/
scripts/
build/
assets/
```

---

## Main Features & Entry Points

### 1. JWT Authentication & Role-Based Access
- Secure login and registration with JWT token issuance
- Role-based API protection (admin/user)
- **Key code:**
  - Middleware: `internal/middleware/auth.go` (`AuthMiddleware`)
  - Login/Registration: `internal/user/handler/user.go` (`LoginHandler`, `RegisterUserHandler`)

### 2. User Management
- Register, login, CRUD for users
- Get user info, get user with orders, get user with Redis cache
- **Key code:**
  - Handlers: `internal/user/handler/user.go` (`GetUserHandler`, `CreateUserHandler`, `UpdateUserHandler`, `DeleteUserHandler`, `GetUserWithOrdersHandler`, `GetUserWithCacheHandler`)
  - Service: `internal/user/service/user_service.go` (`RegisterUser`, `LoginUser`, `GetUserByID`, `GetUserWithOrders`, `GetUserByIDWithCache`, `UpdateUser`, `DeleteUser`)
  - Repository: `internal/user/repository/user_repository.go` (CRUD, cache logic)

### 3. Order Management
- Get order by ID, get all orders for a user
- **Key code:**
  - Handlers: `internal/order/handler/order.go` (`GetOrderHandler`)
  - Service: `internal/order/service/order_service.go` (`GetOrderByID`)
  - Repository: `internal/order/repository/order_repository.go` (`GetOrderByID`, `GetOrdersByUserID`)

### 4. Redis Caching
- Fast user lookup with Redis, fallback to DB
- **Key code:**
  - Redis client: `pkg/redisclient/redis.go`
  - User repository: `GetUserByIDWithCache`

### 5. Layered Architecture
- Handlers: HTTP request/response logic
- Services: Business logic, orchestration
- Repositories: Data access (DB, cache)
- Middleware: Cross-cutting concerns (auth)

---

## Quick Start

1. Install dependencies: `go mod tidy`
2. Copy environment file: `cp .env.example .env`
3. Configure your database and environment variables in `.env`
4. Start the server: `go run cmd/server/main.go`
5. Open [Swagger UI](http://localhost:8080/swagger/index.html) to explore and test the API

---

## Troubleshooting & Maintenance

- **API Testing:** Use Swagger UI or Postman. For protected endpoints, include a valid JWT token.
- **Common Issues:**
  - Database errors: Check `.env` and DB status
  - JWT errors: Ensure `Authorization` header is set
- **Run all tests:** `go test ./...`
- **Update dependencies:** `go get -u` and `go mod tidy`
- **Keep this README and API docs up to date as you add features.**

---

**For more details, see code comments and Swagger documentation.**

