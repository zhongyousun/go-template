package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "go-template/docs"
	"go-template/internal/db"
	"go-template/internal/middleware"
	orderhandler "go-template/internal/order/handler"
	userhandler "go-template/internal/user/handler"
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
	// These endpoints are open to everyone (no authentication required)
	r.POST("/login", userhandler.LoginHandler)
	r.POST("/register", userhandler.RegisterUserHandler)

	// Protected routes
	// These endpoints require authentication (JWT or session)
	authorized := r.Group("/user", middleware.AuthMiddleware())
	{
		authorized.GET("/:id", userhandler.GetUserHandler)
		authorized.PUT("/:id", userhandler.UpdateUserHandler)
		authorized.DELETE("/:id", userhandler.DeleteUserHandler)
		authorized.GET("/:id/orders", userhandler.GetUserWithOrdersHandler)
	}

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User API
	// handler.RegisterUserRoutes(r)

	// Order API
	orderhandler.RegisterOrderRoutes(r)

	r.Run(":8080")
}
