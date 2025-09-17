
## Table of Contents
- [Folder Purpose](#folder-purpose)
- [JWT Authentication & Role-Based API Protection](#jwt-authentication--role-based-api-protection)
- [Service vs Repository Layer: Usage and Best Practices](#service-vs-repository-layer-usage-and-best-practices)


# Go Template Project Structure

```
go.mod
go.sum
cmd/
    server/
        main.go
internal/
    order/
    payment/
    user/
        handler/
        repository/
        service/
pkg/
configs/
docs/
test/
scripts/
build/
assets/
```

## Folder Purpose
- **cmd/**: Project entry points, each subfolder is an executable (e.g. server).
- **internal/**: Internal code, not accessible from outside the project.
- **pkg/**: Shared libraries or utilities for external use.
- **configs/**: Configuration files (e.g. config.yaml, config.json).
- **docs/**: Project documentation, API docs, design notes.
- **test/**: Integration or end-to-end test code.
- **scripts/**: Automation scripts (e.g. deployment, data migration).
- **build/**: CI/CD and build-related files.
- **assets/**: Static resources (e.g. images, frontend files).

## JWT Authentication & Role-Based API Protection

This project implements JWT authentication and role-based API protection using Gin.

### How it works
- **Login**: `/login` issues a JWT token with user info and role in the payload.
- **JWT Middleware**: Checks the token, parses claims, and stores user info (id, email, role) in Gin Context.
- **Role-Based Protection**: Handlers (e.g. `GetUserHandler`) can access caller info via `c.Get("role")`, `c.Get("email")`, etc. and control access based on role.

### Key Code Locations
- **JWT Middleware**: [`internal/middleware/auth.go`](internal/middleware/auth.go)
- **Login & Token Generation**: [`internal/user/handler/user.go`](internal/user/handler/user.go), function `LoginHandler`
- **API Protection Example**: [`internal/user/handler/user.go`](internal/user/handler/user.go), function `GetUserHandler`

### Usage Example
```go
// In handler
role, _ := c.Get("role")
if role == "admin" {
    // allow admin access
} else {
    // restrict access
}
```

### How to test
1. Register and login to get a JWT token.
2. Use the token in the `Authorization` header (format: `Bearer <token>`) for protected APIs.
3. Try accessing APIs with different roles to see permission control in action.

See the above files for implementation details and reference code.

---

## Service vs Repository Layer: Usage and Best Practices

This project follows a clean separation between the **service** and **repository** layers:

- **Repository Layer** (`internal/user/repository/`):
    - Responsible for direct data access (e.g., SQL queries, CRUD operations).
    - Exposes interfaces (e.g., `UserRepository`) and concrete implementations (e.g., `userRepository`).
    - Should not contain business logic—only data persistence and retrieval.

- **Service Layer** (`internal/user/service/`):
    - Responsible for business logic, validation, and orchestration.
    - Depends on repository interfaces, not concrete implementations (for testability and flexibility).
    - Handles tasks like password hashing, authentication, and combining multiple repository calls.
    - Exposes methods for handlers to use (e.g., `RegisterUser`, `LoginUser`, `GetUserByID`).

### Example Usage

**Handler Layer:**

```go
// In handler
db := c.MustGet("db").(*sql.DB)
repo := repository.NewUserRepository(db)
userService := service.NewUserService(repo)
user, err := userService.GetUserByID(id)
```

**Service Layer:**

```go
// In service
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
        return s.Repo.GetUserByID(id)
}
```

**Repository Layer:**

```go
// In repository
func (r *userRepository) GetUserByID(id int64) (*model.User, error) {
        // SQL query to fetch user by ID
}
```

### Why This Separation?

- Keeps business logic out of data access code.
- Makes unit testing easier (service can be tested with mock repositories).
- Allows swapping data sources (e.g., switch from SQL to NoSQL) with minimal changes.
- Improves code readability and maintainability.

**Tip:**
Always let handlers call service methods, and let services call repository methods. Avoid letting handlers call repository methods directly.

For additional notes or examples, please update this document.

