package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Order struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userId"`
	TotalPrice int    `json:"totalPrice"`
	Status     string `json:"status"`
}

type OrderModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m OrderModel) CreateOrder(order *Order) (int, error) {
	query := `
		INSERT INTO orders (user_id, total_price, status) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`
	args := []interface{}{order.UserID, order.TotalPrice, order.Status}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID)
}

func (m OrderModel) GetOrder(id int) (*Order, error) {
	query := `
		SELECT id, user_id, total_price, status
		FROM orders
		WHERE id = $1
	`
	var order Order
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (m OrderModel) UpdateOrder(order *Order) error {
	query := `
		UPDATE orders
		SET user_id = $1, total_price = $2, status = $3
		WHERE id = $4
		RETURNING id
	`
	args := []interface{}{order.UserID, order.TotalPrice, order.Status, order.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID)
}

func (m OrderModel) DeleteOrder(id int) error {
	query := `
		DELETE FROM orders
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
