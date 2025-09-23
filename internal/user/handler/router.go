package handler

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.Engine) {
	r.GET("/user/:id", GetUserHandler)
	r.POST("/user", CreateUserHandler)
	r.PUT("/user/:id", UpdateUserHandler)
	r.DELETE("/user/:id", DeleteUserHandler)
	r.POST("/register", RegisterUserHandler)
	r.POST("/login", LoginHandler)
	r.GET("/user/:id/orders", GetUserWithOrdersHandler)

}
