package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	Role      string    `json:"role"`
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GetUserHandler godoc
// @Summary Get user info
// @Description Get user data by ID
// @Tags user
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /user/{id} [get]
func GetUserHandler(c *gin.Context) {
	// Get db from context
	db := c.MustGet("db").(*sql.DB)

	id := c.Param("id")

	// Get JWT payload from Gin Context
	jwtRole, _ := c.Get("role")
	if jwtRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	ctx, cancel := context.WithTimeout(c, 3*time.Second) // 3 seconds timeout
	defer cancel()
	var user User
	err := db.QueryRowContext(ctx, `
		SELECT "id", "name", "email", "password", "createdAt", "role"
		FROM "user" 
		WHERE "id" = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.Role)

	if err == sql.ErrNoRows {
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
// @Param user body User true "User Info"
// @Success 201 {object} map[string]interface{}
// @Router /user [post]
func CreateUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := c.MustGet("db").(*sql.DB)
	err := db.QueryRow(`
		INSERT INTO "user" ("name", "email", "password", "role")
		VALUES ($1, $2, $3, $4)
		RETURNING "id", "createdAt"
	`, user.Name, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
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
// @Param user body User true "User Info"
// @Success 200 {object} User
// @Router /user/{id} [put]
func UpdateUserHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := db.Exec(`
		UPDATE "user" SET "name"=$1, "email"=$2, "password"=$3, "role"=$4 WHERE "id"=$5
	`, user.Name, user.Email, user.Password, user.Role, id)
	if err != nil {
		log.Printf("Update failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.ID, _ = strconv.Atoi(id)
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
	id := c.Param("id")
	res, err := db.Exec(`DELETE FROM "user" WHERE "id"=$1`, id)
	if err != nil {
		log.Printf("Delete failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
// @Param user body User true "User Info"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func RegisterUserHandler(c *gin.Context) {
	var user User
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
func LoginHandler(c *gin.Context) {
	var creds LoginRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	db := c.MustGet("db").(*sql.DB)
	var user User
	err := db.QueryRow(`SELECT "id", "name", "email", "password", "role" FROM "user" WHERE "email"=$1`, creds.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	} else if err != nil {
		log.Printf("Login query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
