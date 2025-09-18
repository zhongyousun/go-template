package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"go-template/internal/common/commonmodel"
	orderrepo "go-template/internal/order/repository"
	usermodel "go-template/internal/user/model"
	"go-template/internal/user/repository"
	"go-template/internal/user/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

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
	orderRepo := orderrepo.NewOrderRepository(db)
	userService := service.NewUserService(repo, orderRepo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid user id",
			Code:    http.StatusBadRequest,
			Details: "User ID must be a valid integer",
		})
		return
	}
	jwtRole, _ := c.Get("role")
	if jwtRole != "admin" {
		c.JSON(http.StatusForbidden, commonmodel.ErrorResponse{
			Error:   "Forbidden",
			Code:    http.StatusForbidden,
			Details: "You do not have permission to access this resource",
		})
		return
	}
	userObj, err := userService.GetUserByID(id)
	if err == sql.ErrNoRows || userObj == nil {
		c.JSON(http.StatusNotFound, commonmodel.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Details: "No user found with the given ID",
		})
		return
	} else if err != nil {
		log.Printf("Query failed: %v", err)
		c.JSON(http.StatusInternalServerError, commonmodel.ErrorResponse{
			Error:   "Internal server error",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, userObj)
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
	var user usermodel.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}
	db := c.MustGet("db").(*sql.DB)
	repo := repository.NewUserRepository(db)
	orderRepo := orderrepo.NewOrderRepository(db)
	userService := service.NewUserService(repo, orderRepo)
	err := userService.RegisterUser(&user)
	if err != nil {
		log.Printf("Create failed: %v", err)
		c.JSON(http.StatusInternalServerError, commonmodel.ErrorResponse{
			Error:   "Create failed",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
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
	orderRepo := orderrepo.NewOrderRepository(db)
	userService := service.NewUserService(repo, orderRepo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid user id",
			Code:    http.StatusBadRequest,
			Details: "User ID must be a valid integer",
		})
		return
	}
	var user usermodel.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}
	user.ID = int(id)
	err = userService.UpdateUser(&user)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, commonmodel.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Details: "No user found with the given ID",
		})
		return
	} else if err != nil {
		log.Printf("Update failed: %v", err)
		c.JSON(http.StatusInternalServerError, commonmodel.ErrorResponse{
			Error:   "Update failed",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
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
	orderRepo := orderrepo.NewOrderRepository(db)
	userService := service.NewUserService(repo, orderRepo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid user id",
			Code:    http.StatusBadRequest,
			Details: "User ID must be a valid integer",
		})
		return
	}
	err = userService.DeleteUser(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, commonmodel.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Details: "No user found with the given ID",
		})
		return
	} else if err != nil {
		log.Printf("Delete failed: %v", err)
		c.JSON(http.StatusInternalServerError, commonmodel.ErrorResponse{
			Error:   "Delete failed",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
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
// @Failure 400 {object} commonmodel.ErrorResponse
// @Failure 500 {object} commonmodel.ErrorResponse
// @Router /register [post]
func RegisterUserHandler(c *gin.Context) {
	var user usermodel.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Details: "Request body is not valid JSON or missing required fields",
		})
		return
	}
	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Name, email, and password are required",
			Code:    http.StatusBadRequest,
			Details: "All fields are required for registration",
		})
		return
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonmodel.ErrorResponse{
			Error:   "Failed to hash password",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, commonmodel.ErrorResponse{
			Error:   "Register failed",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
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
// @Failure 400 {object} commonmodel.ErrorResponse
// @Failure 401 {object} commonmodel.ErrorResponse
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
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid input",
			Code:    http.StatusBadRequest,
			Details: "Request body is not valid JSON or missing required fields",
		})
		return
	}
	// Get database connection from Gin context
	db := c.MustGet("db").(*sql.DB)
	// Initialize user repository
	repo := repository.NewUserRepository(db)
	// Initialize order repository
	orderRepo := orderrepo.NewOrderRepository(db)
	// Initialize user service
	userService := service.NewUserService(repo, orderRepo)

	// Authenticate user and generate JWT token
	token, err := userService.LoginUser(creds.Email, creds.Password)
	if err != nil {
		// If authentication fails, return 401 Unauthorized
		c.JSON(http.StatusUnauthorized, commonmodel.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Details: err.Error(),
		})
		return
	}
	// Return JWT token in response
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ...existing code...
// GetUserWithOrdersHandler godoc
// @Summary Get user info with orders
// @Description Get user data and all orders by user ID
// @Tags user
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} model.UserWithOrders
// @Failure 400 {object} commonmodel.ErrorResponse
// @Failure 404 {object} commonmodel.ErrorResponse
// @Router /user/{id}/orders [get]
func GetUserWithOrdersHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	repo := repository.NewUserRepository(db)
	orderRepo := orderrepo.NewOrderRepository(db)
	userService := service.NewUserService(repo, orderRepo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.ErrorResponse{
			Error:   "Invalid user id",
			Code:    http.StatusBadRequest,
			Details: "User ID must be a valid integer",
		})
		return
	}
	jwtRole, _ := c.Get("role")
	if jwtRole != "admin" {
		c.JSON(http.StatusForbidden, commonmodel.ErrorResponse{
			Error:   "Forbidden",
			Code:    http.StatusForbidden,
			Details: "You do not have permission to access this resource",
		})
		return
	}
	result, err := userService.GetUserWithOrders(id)
	if err != nil {
		log.Printf("Query failed: %v", err)
		c.JSON(http.StatusInternalServerError, commonmodel.ErrorResponse{
			Error:   "Internal server error",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, commonmodel.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Details: "No user found with the given ID",
		})
		return
	}
	c.JSON(http.StatusOK, result)
}
