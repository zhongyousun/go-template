package main

import (
	docs "go-template/docs"
	"go-template/internal/db"
	"go-template/internal/middleware"
	"go-template/internal/user/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Example API
// @version 1.0
// @description Gin + Swagger example
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables
	godotenv.Load()
	// Initialize DB
	db.InitDB()
	defer db.DB.Close()

	r := gin.Default()
	// Middleware: Inject DB into Gin Context
	r.Use(func(c *gin.Context) {
		c.Set("db", db.DB)
		c.Next()
	})
	// r.Use(middleware.AuthMiddleware())

	// Public routes
	r.POST("/login", handler.LoginHandler)
	r.POST("/register", handler.RegisterUserHandler)

	// Protected routes
	authorized := r.Group("/user", middleware.AuthMiddleware())
	{
		authorized.GET("/:id", handler.GetUserHandler)
		authorized.PUT("/:id", handler.UpdateUserHandler)
		authorized.DELETE("/:id", handler.DeleteUserHandler)
	}

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User API
	// handler.RegisterUserRoutes(r)

	r.Run(":8080")
}
