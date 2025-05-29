package service

import (
	"OrderMatchingEngine/internal/dto"
	"OrderMatchingEngine/internal/models"
	"OrderMatchingEngine/internal/transformer"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type TradeService interface {
	ListTrades(ctx context.Context, req *dto.ListTradesRequest) (*dto.ListTradesResponse, error)
	BulkCreateTrades(ctx context.Context, orderId string, orderSide dto.Side, symbol string, matches []*dto.OrderBookMatches) error 
}

type tradeService struct {
	db *sql.DB
}

func NewTradeService(db *sql.DB) TradeService {
	return &tradeService{db: db}
}

func (s *tradeService) ListTrades(ctx context.Context, req *dto.ListTradesRequest) (*dto.ListTradesResponse, error) {
	tradeModel := &models.Trade{
		Symbol: req.Symbol,
	}
	tradeModels, err := tradeModel.ListTrades(ctx, s.db)
	if err != nil {
		return nil, err
	}

	tradeStructs := transformer.TransformTradeModelsToTrades(tradeModels)

	return &dto.ListTradesResponse{
		Trades: tradeStructs,
	}, nil
}

func (s *tradeService) BulkCreateTrades(ctx context.Context, orderId string, orderSide dto.Side, symbol string, matches []*dto.OrderBookMatches) error {

	trades := make([]*models.Trade, len(matches))
	for i, match := range matches {
		buyID := ""
		sellID := ""
		if orderSide == dto.SideBuy {
			buyID = orderId
			sellID = match.MatchOrderID
		} else {
			buyID = match.MatchOrderID
			sellID = orderId
		}
		trades[i] = &models.Trade{
			ID: uuid.New().String(),
			BuyOrderID: buyID,
			SellOrderID: sellID,
			Price: match.Price,
			Quantity: match.Quantity,
			Timestamp: match.Timestamp,
			Symbol: symbol,
		}
	}

	tradeModel := &models.Trade{}
	err := tradeModel.BatchCreateTrades(ctx, s.db, trades)
	if err != nil {
		return err
	}
	return nil
}