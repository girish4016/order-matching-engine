package models

import (
	"context"
	"database/sql"
	"time"
)

type Trade struct {
	ID          string    `json:"id"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
	Symbol      string    `json:"symbol"`
}

func (t *Trade) ListTrades(ctx context.Context, db *sql.DB) ([]*Trade, error) {
	query := "SELECT * FROM trades WHERE symbol = ?"
	rows, err := db.QueryContext(ctx, query, t.Symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trades := []*Trade{}
	for rows.Next() {
		var trade Trade
		err := rows.Scan(&trade.ID, &trade.BuyOrderID, &trade.SellOrderID, &trade.Price, &trade.Quantity, &trade.Timestamp, &trade.Symbol)
		if err != nil {
			return nil, err
		}
		trades = append(trades, &trade)
	}
	return trades, nil
}

func (t *Trade) BatchCreateTrades(ctx context.Context, db *sql.DB, trades []*Trade) error {
	query := "INSERT INTO trades (id, buy_order_id, sell_order_id, price, quantity, timestamp, symbol) VALUES (?, ?, ?, ?, ?, ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStmt := tx.StmtContext(ctx, stmt)
	defer txStmt.Close()

	for _, trade := range trades {
		_, err := txStmt.ExecContext(ctx,
			trade.ID,
			trade.BuyOrderID,
			trade.SellOrderID,
			trade.Price,
			trade.Quantity,
			trade.Timestamp,
			trade.Symbol,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *Trade) TableName() string {
	return "trades"
}
