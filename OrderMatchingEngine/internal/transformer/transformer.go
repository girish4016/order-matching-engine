package transformer

import (
	"OrderMatchingEngine/internal/dto"
	"OrderMatchingEngine/internal/models"
)

func TransformOrderToModel(order *dto.Order) *models.Order {
	return &models.Order{
		ID:                order.ID,
		Symbol:            order.Symbol,
		Side:              string(order.Side),
		Type:              string(order.Type),
		Price:             order.Price,
		Quantity:          order.Quantity,
		RemainingQuantity: order.RemainingQuantity,
		Status:            string(order.Status),
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}
}

func TransformOrderModelToOrder(orderModel *models.Order) *dto.Order {
	return &dto.Order{
		ID:                orderModel.ID,
		Symbol:            orderModel.Symbol,
		Side:              dto.Side(orderModel.Side),
		Type:              dto.Type(orderModel.Type),
		Price:             orderModel.Price,
		Quantity:          orderModel.Quantity,
		RemainingQuantity: orderModel.RemainingQuantity,
		Status:            dto.Status(orderModel.Status),
		CreatedAt:         orderModel.CreatedAt,
		UpdatedAt:         orderModel.UpdatedAt,
	}
}

func TransformOrderBookMatchToOrderModel(orderBookMatch *dto.OrderBookMatches) *models.Order {
	status := dto.StatusFilled
	if orderBookMatch.RemainingQuantity > 0 {
		status = dto.StatusPartial
	}

	return &models.Order{
		ID:                orderBookMatch.MatchOrderID,
		RemainingQuantity: orderBookMatch.RemainingQuantity,
		Status:            string(status),
		UpdatedAt:         orderBookMatch.Timestamp,
	}
}

func TransformOrderBookMatchesToOrderModels(orders []*dto.OrderBookMatches) []*models.Order {
	models := make([]*models.Order, len(orders))
	for i, order := range orders {
		models[i] = TransformOrderBookMatchToOrderModel(order)
	}
	return models
}

func TransformTradeModelToTrade(tradeModel *models.Trade) *dto.Trade {
	return &dto.Trade{
		TradeId:     tradeModel.ID,
		BuyOrderID:  tradeModel.BuyOrderID,
		SellOrderID: tradeModel.SellOrderID,
		Price:       tradeModel.Price,
		Quantity:    tradeModel.Quantity,
		Timestamp:   tradeModel.Timestamp,
		Symbol:      tradeModel.Symbol,
	}
}

func TransformTradeModelsToTrades(tradeModels []*models.Trade) []*dto.Trade {
	trades := make([]*dto.Trade, len(tradeModels))
	for i, trade := range tradeModels {
		trades[i] = TransformTradeModelToTrade(trade)
	}
	return trades
}