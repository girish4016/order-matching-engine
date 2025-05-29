package service

import (
	"OrderMatchingEngine/internal/dto"
	"OrderMatchingEngine/internal/models"
	"OrderMatchingEngine/internal/transformer"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	PlaceOrder(ctx context.Context, req *dto.PlaceOrderRequest) (*dto.PlaceOrderResponse, error)
	GetOrder(ctx context.Context, id string) (*dto.Order, error)
	CancelOrder(ctx context.Context, id string) (*dto.Order, error)
	GetOrderBook(ctx context.Context, req *dto.GetOrderBookRequest) (*dto.GetOrderBookResponse, error)
}

type orderService struct {
	OrderBookMap map[string]*dto.OrderBook
	DB           *sql.DB
	TradeService TradeService
}

func (s *orderService) PlaceOrder(ctx context.Context, req *dto.PlaceOrderRequest) (*dto.PlaceOrderResponse, error) {
	orderBook, ok := s.OrderBookMap[req.Symbol]
	if !ok {
		orderBook = dto.NewOrderBook(req.Symbol)
		s.OrderBookMap[req.Symbol] = orderBook
	}
	order := &dto.Order{
		ID:                uuid.New().String(),
		Symbol:            req.Symbol,
		Side:              req.Side,
		Type:              req.Type,
		Price:             req.Price,
		Quantity:          req.Quantity,
		RemainingQuantity: req.Quantity,
		Status:            dto.StatusPending,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	matches := orderBook.AddOrder(order)

	fmt.Print(matches)
	// repo changes to be added here
	orderObj := *transformer.TransformOrderToModel(order)
	err := orderObj.CreateOrder(ctx, s.DB)
	if err != nil {
		return nil, err
	}

	matchOrdersObj := transformer.TransformOrderBookMatchesToOrderModels(matches)
	err = orderObj.BatchUpdateOrders(ctx, s.DB, matchOrdersObj)
	if err != nil {
		return nil, err
	}

	// trades to be recorded here using trade repository
	err = s.TradeService.BulkCreateTrades(ctx, order.ID, order.Side, order.Symbol, matches)
	if err != nil {
		return nil, err
	}

	return &dto.PlaceOrderResponse{
		Order: order,
	}, nil
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*dto.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	order := models.Order{
		ID: id,
	}
	err := order.GetOrder(ctx, s.DB)
	if err != nil {
		return nil, err
	}

	orderStruct := transformer.TransformOrderModelToOrder(&order)

	return orderStruct, nil
}

func (s *orderService) CancelOrder(ctx context.Context, id string) (*dto.Order, error) {
	orderModel := &models.Order{
		ID: id,
	}
	err := orderModel.GetOrder(ctx, s.DB)
	if err != nil {
		return nil, err
	}

	orderModel.Status = string(dto.StatusCancelled)
	orderModel.UpdatedAt = time.Now().UTC()
	err = orderModel.UpdateOrder(ctx, s.DB)

	orderStruct := transformer.TransformOrderModelToOrder(orderModel)

	s.OrderBookMap[orderModel.Symbol].RemoveOrder(orderModel.ID)

	return orderStruct, nil
}

func (s *orderService) GetOrderBook(ctx context.Context, req *dto.GetOrderBookRequest) (*dto.GetOrderBookResponse, error) {
	fmt.Println("blastoise2")
	if orderBook, ok := s.OrderBookMap[req.Symbol]; ok {
		fmt.Println("blastoise3")
		return &dto.GetOrderBookResponse{
			OrderBook: *orderBook,
		}, nil
	} else {
		fmt.Println("blastoise4")
		s.OrderBookMap[req.Symbol] = dto.NewOrderBook(req.Symbol)
		return &dto.GetOrderBookResponse{
			OrderBook: *s.OrderBookMap[req.Symbol],
		}, nil
	}
}

func NewOrderService(db *sql.DB, tradeService *TradeService) OrderService {
	return &orderService{
		OrderBookMap: make(map[string]*dto.OrderBook),
		DB:           db,
		TradeService: *tradeService,
	}
}
