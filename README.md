# Go Template

A modern, scalable Go web project template featuring:


## 🚀 Project Overview


This template provides a modern, scalable foundation for Go web projects, featuring:

- **Authentication & Authorization**: Secure JWT authentication and fine-grained, role-based API access control (admin/user).
  - Code: `internal/middleware/auth.go` (`AuthMiddleware`), `internal/user/handler/user.go` (`LoginHandler`, `RegisterUserHandler`)
- **User & Order Management**: Full CRUD operations for users and orders, with business logic separated by domain.
  - Code: `internal/user/handler/user.go`, `internal/user/service/user_service.go`, `internal/user/repository/user_repository.go`, `internal/order/handler/order.go`, `internal/order/service/order_service.go`, `internal/order/repository/order_repository.go`
- **Layered Architecture**: Clean separation of concerns—handlers (HTTP), services (business logic), repositories (DB/cache), and middleware (auth, etc.).
  - Code: See `internal/user/`, `internal/order/`, `internal/middleware/`, `internal/common/`, `internal/db/`
- **Transactional Operations**: Unit of Work pattern for atomic multi-table operations (e.g., register user and create order in one transaction), ensuring data consistency with automatic rollback on failure.
  - Code: `internal/user/handler/user.go` (`RegisterUserWithOrderHandler`), `internal/user/service/user_service.go` (`RegisterUserWithOrder`), `internal/db/transaction_manager.go`, `internal/user/repository/user_repository.go`, `internal/order/repository/order_repository.go`
- **Redis Caching**: Fast user lookup with Redis, seamlessly falling back to the database if needed.
  - Code: `pkg/redisclient/redis.go`, `internal/user/repository/user_repository.go` (`GetUserByIDWithCache`)
- **API Documentation**: Auto-generated Swagger docs for easy API exploration and testing.
  - Usage: Start the server and open [Swagger UI](http://localhost:8080/swagger/index.html) to explore and test the API
  - To update API docs after code changes:
    1. Regenerate docs: `swag init -g cmd/server/main.go -o docs`
    2. Restart the server: `go run cmd/server/main.go`
- **Testing Support**: Comprehensive unit and integration tests under `test/`, easily run with `go test ./...`.
  - Code: `test/`

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

## Quick Start


To get started with this template:

1. **Install Go dependencies**
  ```
  go mod tidy
  ```
2. **Copy and edit environment variables**
  ```
  cp .env.example .env
  # Edit .env to match your database and environment settings
  ```
3. **Generate Swagger API documentation** (requires [swag](https://github.com/swaggo/swag) installed)
  ```
  swag init -g cmd/server/main.go -o docs
  ```
4. **Start the server**
  ```
  go run cmd/server/main.go
  ```
5. **Explore and test the API**
  - Open [Swagger UI](http://localhost:8080/swagger/index.html) in your browser
  - Use Swagger UI or Postman for API testing (JWT required for protected endpoints)




