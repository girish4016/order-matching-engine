package api

import (
	"OrderMatchingEngine/internal/controllers"
	"OrderMatchingEngine/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Router struct {
	orderController *controllers.OrderController
	tradeController *controllers.TradeController
}

func NewRouter(db *sql.DB) *Router {
	tradeService := service.NewTradeService(db)
	orderService := service.NewOrderService(db, &tradeService)
	return &Router{
		orderController: controllers.NewOrderController(
			orderService,
		),
		tradeController: controllers.NewTradeController(
			tradeService,
		),
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	// Place new order
	engine.POST("/orders", r.orderController.PlaceOrder)

	// Get specific order
	engine.GET("/order/:id", r.orderController.GetOrder)

	// // Cancel/Delete order
	engine.DELETE("/order/:id", r.orderController.CancelOrder)

	// Get orderbook for a symbol
	engine.GET("/orderbook", r.orderController.GetOrderBook)

	// // Trades endpoint - this wasn't implemented in the controller yet
	engine.GET("/trades/:symbol", r.tradeController.ListTrades)
}
