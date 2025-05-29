package dto

import "time"

type Side string

type Type string

type Status string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"

	TypeLimit  Type = "limit"
	TypeMarket Type = "market"

	StatusPending   Status = "pending"
	StatusPartial   Status = "partial"
	StatusFilled    Status = "filled"
	StatusCancelled Status = "cancelled"
)

type Order struct {
	ID                string    `json:"id"`
	Symbol            string    `json:"symbol"`
	Side              Side      `json:"side"`
	Type              Type      `json:"type"`
	Price             float64   `json:"price"`
	Quantity          float64   `json:"quantity"`
	RemainingQuantity float64   `json:"remaining_quantity"`
	Status            Status    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PlaceOrderRequest struct {
	Symbol   string  `json:"symbol"`
	Side     Side    `json:"side"`
	Type     Type    `json:"type"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type PlaceOrderResponse struct {
	Order *Order `json:"order"`
}

type GetOrderBookRequest struct {
	Symbol string `json:"symbol"`
}

type GetOrderBookResponse struct {
	OrderBook OrderBook `json:"order_book"`
}
