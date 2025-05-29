package controllers

import (
	"OrderMatchingEngine/internal/dto"
	"OrderMatchingEngine/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TradeController struct {
	tradeService service.TradeService
}

func NewTradeController(tradeService service.TradeService) *TradeController {
	return &TradeController{tradeService: tradeService}
}

func (t *TradeController) ListTrades(ctx *gin.Context) {
	req := &dto.ListTradesRequest{
		Symbol: ctx.Param("symbol"),
	}
	trades, err := t.tradeService.ListTrades(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, trades)
}
