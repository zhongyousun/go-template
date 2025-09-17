package handler

import "github.com/gin-gonic/gin"

func RegisterOrderRoutes(r *gin.Engine) {
	r.GET("/order/:id", GetOrderHandler)
}
