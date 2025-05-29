package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"OrderMatchingEngine/internal/dto"
	"OrderMatchingEngine/internal/service"
)

type OrderController struct {
	orderService service.OrderService
}

func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	var placeOrderRequest dto.PlaceOrderRequest
	if err := ctx.ShouldBindJSON(&placeOrderRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.orderService.PlaceOrder(ctx, &placeOrderRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *OrderController) GetOrderBook(ctx *gin.Context) {
	req := &dto.GetOrderBookRequest{
		Symbol: ctx.Query("symbol"),
	}
	fmt.Println("blastoise")
	orderBook, err := c.orderService.GetOrderBook(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orderBook)
}

func (c *OrderController) CancelOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.orderService.CancelOrder(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *OrderController) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	order, err := c.orderService.GetOrder(ctx, id)
	if err != nil {
		if err.Error() == "order with ID not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func NewOrderController(service service.OrderService) *OrderController {
	return &OrderController{orderService: service}
}
