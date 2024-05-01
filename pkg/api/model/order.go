package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/godra-y/go-project/pkg/api/validator"
	"log"
	"time"
)

type Order struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	CreatedAt string `json:"created_at"`
}

type OrderModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (om OrderModel) GetAll(filters Filters) ([]*Order, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, product_id, quantity, created_at
		FROM orders
		ORDER BY %s %s, id ASC
		LIMIT $1 OFFSET $2
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{filters.limit(), filters.offset()}

	rows, err := om.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			om.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var orders []*Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&totalRecords, &order.ID, &order.ProductID, &order.Quantity, &order.CreatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return orders, metadata, nil
}

func (om OrderModel) Insert(order *Order) error {
	query := `
		INSERT INTO orders (product_id, quantity) 
		VALUES ($1, $2) 
		RETURNING id, created_at
		`

	args := []interface{}{order.ProductID, order.Quantity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return om.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID, &order.CreatedAt)
}

func (om OrderModel) Get(id int) (*Order, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, product_id, quantity, created_at
		FROM orders
		WHERE id = $1
		`

	var order Order
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := om.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&order.ID, &order.ProductID, &order.Quantity, &order.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve order with id: %v, %w", id, err)
	}
	return &order, nil
}

func (om OrderModel) GetOrdersByProduct(productID int, filters Filters) ([]*Order, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, product_id, quantity, created_at
		FROM orders
		WHERE product_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{productID, filters.limit(), filters.offset()}

	rows, err := om.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			om.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var orders []*Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&totalRecords, &order.ID, &order.ProductID, &order.Quantity, &order.CreatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return orders, metadata, nil
}

func (om OrderModel) Update(order *Order) error {
	query := `
		UPDATE orders
		SET product_id = $1, quantity = $2
		WHERE id = $3
		RETURNING id
		`

	args := []interface{}{order.ProductID, order.Quantity, order.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return om.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID)
}

func (om OrderModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM orders
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := om.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateOrder(v *validator.Validator, order *Order) {
	v.Check(order.ProductID > 0, "product_id", "must be a positive value")
	v.Check(order.Quantity > 0, "quantity", "must be a positive value")
}
