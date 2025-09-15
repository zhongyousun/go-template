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

For additional notes or examples, please update this document.
