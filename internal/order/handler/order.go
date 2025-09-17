package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"go-template/internal/order/repository"
	"go-template/internal/order/service"

	"github.com/gin-gonic/gin"
)

// GetOrderHandler godoc
// @Summary Get order info
// @Description Get order data by ID
// @Tags order
// @Param id path int true "Order ID"
// @Success 200 {object} model.Order
// @Router /order/{id} [get]
func GetOrderHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	repo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(repo)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}
	order, err := orderService.GetOrderByID(id)
	if err == sql.ErrNoRows || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, order)
}
