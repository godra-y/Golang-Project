package models

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type OrderModel struct {
	DB *sql.DB
}

func (om OrderModel) Insert(order *Order) error {
	query := `
        INSERT INTO orders (user_id, total_price, status) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `

	args := []interface{}{order.UserID, order.TotalPrice, order.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := om.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (om OrderModel) Get(id int) (*Order, error) {
	query := `
        SELECT id, user_id, total_price, status
        FROM orders
        WHERE id = $1
    `

	var order Order

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := om.DB.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (om OrderModel) Update(order *Order) error {
	query := `
        UPDATE orders
        SET user_id = $2, total_price = $3, status = $4
        WHERE id = $1
    `

	args := []interface{}{order.ID, order.UserID, order.TotalPrice, order.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := om.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (om OrderModel) Delete(id int) error {
	query := `
        DELETE FROM orders
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := om.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
