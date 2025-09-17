package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"go-template/internal/user/model"
	"go-template/internal/user/repository"
	"go-template/internal/user/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GetUserHandler godoc
// @Summary Get user info
// @Description Get user data by ID
// @Tags user
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Router /user/{id} [get]
func GetUserHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	jwtRole, _ := c.Get("role")
	if jwtRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	user, err := userService.GetUserByID(id)
	if err == sql.ErrNoRows || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Printf("Query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUserHandler godoc
// @Summary Create new user
// @Description Add a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.User true "User Info"
// @Success 201 {object} map[string]interface{}
// @Router /user [post]
func CreateUserHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := c.MustGet("db").(*sql.DB)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	err := userService.RegisterUser(&user)
	if err != nil {
		log.Printf("Create failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// UpdateUserHandler godoc
// @Summary Edit user
// @Description Edit user data by ID
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.User true "User Info"
// @Success 200 {object} model.User
// @Router /user/{id} [put]
func UpdateUserHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = int(id)
	err = userService.UpdateUser(&user)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Printf("Update failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUserHandler godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags user
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204 {string} string ""
// @Router /user/{id} [delete]
func DeleteUserHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	err = userService.DeleteUser(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Printf("Delete failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

// RegisterUserHandler godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.User true "User Info"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func RegisterUserHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and password are required"})
		return
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	db := c.MustGet("db").(*sql.DB)
	err = db.QueryRow(`
		INSERT INTO "user" ("name", "email", "password", "role")
		VALUES ($1, $2, $3, $4)
		RETURNING "id", "createdAt"
	`, user.Name, user.Email, string(hashedPassword), "user").Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		log.Printf("Register failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Register failed"})
		return
	}
	user.Password = "" // Do not return password
	c.JSON(http.StatusCreated, user)
}

// LoginRequest represents the login request body
// swagger:model
type LoginRequest struct {
	Email    string `json:"email" example:"abc@gmail.com"`
	Password string `json:"password" example:"1234"`
}

// LoginHandler godoc
// @Summary Login
// @Description Login with email and password, returns JWT
// @Tags user
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
// LoginHandler handles user login requests
// 1. Parse login credentials from JSON body
// 2. Initialize repository and service layers
// 3. Authenticate user and generate JWT token
// 4. Return token or error response
func LoginHandler(c *gin.Context) {
	// Parse login credentials from request body
	var creds LoginRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		// If parsing fails, return 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// Get database connection from Gin context
	db := c.MustGet("db").(*sql.DB)
	// Initialize user repository
	repo := repository.NewUserRepository(db)
	// Initialize user service
	userService := service.NewUserService(repo)

	// Authenticate user and generate JWT token
	token, err := userService.LoginUser(creds.Email, creds.Password)
	if err != nil {
		// If authentication fails, return 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// Return JWT token in response
	c.JSON(http.StatusOK, gin.H{"token": token})
}
