package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Order struct {
	ID                string    `json:"id"`
	Symbol            string    `json:"symbol"`
	Side              string    `json:"side"`
	Type              string    `json:"type"`
	Price             float64   `json:"price"`
	Quantity          float64   `json:"quantity"`
	RemainingQuantity float64   `json:"remaining_quantity"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Order model
func (Order) TableName() string {
	return "orders"
}

func (o *Order) CreateOrder(ctx context.Context, db *sql.DB) error {
	query := "INSERT INTO orders (id, symbol, side, type, price, quantity, remaining_quantity, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, query, o.ID, o.Symbol, o.Side, o.Type, o.Price, o.Quantity, o.RemainingQuantity, o.Status, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		fmt.Printf("Error creating order: %v", err)
		return err
	}
	return nil
}

func (o *Order) UpdateOrder(ctx context.Context, db *sql.DB) error {
	query := "UPDATE orders SET remaining_quantity = ?, updated_at = ?, status = ? WHERE id = ?"
	_, err := db.ExecContext(ctx, query, o.RemainingQuantity, o.UpdatedAt, o.Status, o.ID)
	if err != nil {
		fmt.Printf("Error updating order: %v", err)
		return err
	}
	return nil
}

func (o *Order) BatchUpdateOrders(ctx context.Context, db *sql.DB, orders []*Order) error {
	for _, order := range orders {
		query := "UPDATE orders SET remaining_quantity = ?, updated_at = ?, status = ? WHERE id = ?"
		_, err := db.ExecContext(ctx, query, order.RemainingQuantity, order.UpdatedAt, order.Status, order.ID)
		if err != nil {
			fmt.Printf("Error updating order %s: %v", order.ID, err)
			return err
		}
	}
	return nil
}

func (o *Order) GetOrder(ctx context.Context, db *sql.DB) error {
	query := "SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at, updated_at FROM orders WHERE id = ?"
	row := db.QueryRowContext(ctx, query, o.ID)

	err := row.Scan(&o.ID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Quantity, &o.RemainingQuantity, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order with ID %s not found", o.ID)
		}
		fmt.Printf("Error retrieving order: %v", err)
		return err
	}
	return nil
}