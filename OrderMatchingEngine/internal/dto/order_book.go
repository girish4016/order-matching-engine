package dto

import (
	"sort"
	"time"
)

type OrderBook struct {
	Symbol string           `json:"symbol"`
	Bids   []OrderBookEntry `json:"bids"`
	Asks   []OrderBookEntry `json:"asks"`
}

func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{
		Symbol: symbol,
		Bids:   make([]OrderBookEntry, 0),
		Asks:   make([]OrderBookEntry, 0),
	}
}

type OrderBookEntry struct {
	OrderId   string
	Timestamp time.Time
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
}

type OrderBookMatches struct {
	MatchOrderID      string    `json:"match_order_id"`
	Price             float64   `json:"price"`
	Quantity          float64   `json:"quantity"`
	RemainingQuantity float64   `json:"remaining_quantity"`
	Timestamp         time.Time `json:"timestamp"`
}

func (o *OrderBook) AddOrder(order *Order) []*OrderBookMatches {
	matchingOrders := o.MatchOrder(order)

	if order.RemainingQuantity == 0 {
		order.Status = StatusFilled
	} else if order.Type == TypeLimit {
		if order.RemainingQuantity == order.Quantity {
			order.Status = StatusPending
		} else {
			order.Status = StatusPartial
		}
		if order.Side == SideBuy {
			o.Bids = append(o.Bids, OrderBookEntry{
				OrderId:   order.ID,
				Timestamp: order.CreatedAt,
				Price:     order.Price,
				Quantity:  order.Quantity,
			})
		} else {
			o.Asks = append(o.Asks, OrderBookEntry{
				OrderId:   order.ID,
				Timestamp: order.CreatedAt,
				Price:     order.Price,
				Quantity:  order.Quantity,
			})
		}
	} else {
		order.Status = StatusCancelled
	}

	return matchingOrders
}

func (o *OrderBook) MatchOrder(order *Order) []*OrderBookMatches {
	if order.Side == SideBuy {
		if len(o.Asks) == 0 {
			return nil
		}
		// Sort asks by price in descending order of price and timestamp
		sort.Slice(o.Asks, func(i, j int) bool {
			if o.Asks[i].Price == o.Asks[j].Price {
				return o.Asks[i].Timestamp.After(o.Asks[j].Timestamp)
			}
			return o.Asks[i].Price > o.Asks[j].Price
		})
		matches := make([]*OrderBookMatches, 0)
		index := len(o.Asks) - 1
		for index >= 0 && order.RemainingQuantity > 0 {
			ask := &o.Asks[index]
			if ask.Price > order.Price && order.Type == TypeLimit {
				break
			}

			qty := order.RemainingQuantity
			if order.RemainingQuantity > ask.Quantity {
				qty = ask.Quantity
				order.RemainingQuantity -= qty
				ask.Quantity = 0
			} else {
				ask.Quantity -= qty
				order.RemainingQuantity = 0
			}

			match := &OrderBookMatches{
				MatchOrderID:      ask.OrderId,
				Price:             ask.Price,
				Quantity:          qty,
				RemainingQuantity: ask.Quantity,
				Timestamp:         time.Now().UTC(),
			}

			if ask.Quantity == 0 {
				o.Asks = o.Asks[:index]
			}

			matches = append(matches, match)
			index--
		}
		return matches
	} else {
		if o.Bids == nil {
			return nil
		}
		// Sort bids by price in ascending order of price and descending order of timestamp
		sort.Slice(o.Bids, func(i, j int) bool {
			if o.Bids[i].Price == o.Bids[j].Price {
				return o.Bids[i].Timestamp.After(o.Bids[j].Timestamp)
			}
			return o.Bids[i].Price < o.Bids[j].Price
		})
		matches := make([]*OrderBookMatches, 0)
		index := len(o.Bids) - 1
		for index >= 0 && order.RemainingQuantity > 0 {
			bid := &o.Bids[index]
			if bid.Price < order.Price && order.Type == TypeLimit {
				break
			}

			qty := order.RemainingQuantity
			if order.RemainingQuantity > bid.Quantity {
				qty = bid.Quantity
				order.RemainingQuantity -= qty
				bid.Quantity = 0
			} else {
				order.RemainingQuantity = 0
				bid.Quantity -= qty
			}

			match := &OrderBookMatches{
				MatchOrderID:      order.ID,
				Price:             bid.Price,
				Quantity:          qty,
				RemainingQuantity: bid.Quantity,
				Timestamp:         time.Now().UTC(),
			}

			if bid.Quantity == 0 {
				o.Bids = o.Bids[:index]
			}

			matches = append(matches, match)
			index--
		}
		return matches
	}
}

func (o *OrderBook) RemoveOrder(orderId string) {
	for i := 0; i < len(o.Bids); i++ {
		if o.Bids[i].OrderId == orderId {
			o.Bids = append(o.Bids[:i], o.Bids[i+1:]...)
			return
		}
	}
	
	for i := 0; i < len(o.Asks); i++ {
		if o.Asks[i].OrderId == orderId {
			o.Asks = append(o.Asks[:i], o.Asks[i+1:]...)
			return
		}
	}
}