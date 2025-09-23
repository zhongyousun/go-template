package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "go-template/docs"
	"go-template/internal/db"
	"go-template/internal/middleware"
	orderhandler "go-template/internal/order/handler"
	userhandler "go-template/internal/user/handler"
	"go-template/pkg/redisclient"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Example: Initialize sql.DB (legacy/original version)
	// You can use this if you want to use the standard library database/sql API instead of GORM.
	// By default, this project uses GORM for database operations, but you can switch to sql.DB if needed.
	db.InitDB()
	defer db.DB.Close()

	// Example: Initialize GORM DB (recommended/primary usage)
	// This is the main database connection for most use cases in this project.
	// If you want to use GORM's ORM features, use this connection.
	gormDSN := os.Getenv("POSTGRES_CONN")
	gormDB, err := gorm.Open(postgres.Open(gormDSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database (gorm): " + err.Error())
	}

	r := gin.Default()
	// Middleware: Inject both DB into Gin Context
	r.Use(func(c *gin.Context) {
		c.Set("db", db.DB)
		c.Set("gorm", gormDB)
		c.Next()
	})

	// Init Redis
	redisclient.Init()

	// Public routes
	r.POST("/login", userhandler.LoginHandler)
	r.POST("/register", userhandler.RegisterUserHandler)
	r.GET("/userwithcache/:id", userhandler.GetUserWithCacheHandler)
	r.POST("/register_with_order", userhandler.RegisterUserWithOrderHandler)

	// Protected routes
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

	// Order API
	orderhandler.RegisterOrderRoutes(r)

	r.Run(":8080")
}
