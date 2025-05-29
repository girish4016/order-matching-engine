package dto

import "time"

type Trade struct {
	TradeId     string    `json:"trade_id"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
	Symbol      string    `json:"symbol"`
}

type ListTradesRequest struct {
	Symbol string `json:"symbol"`
}

type ListTradesResponse struct {
	Trades []*Trade `json:"trades"`
}