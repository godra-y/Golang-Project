package models

import (
	"context"
	"database/sql"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  int     `json:"category_id"`
}

type ProductModel struct {
	DB *sql.DB
}

func (pm ProductModel) Insert(product *Product) error {
	query := `
        INSERT INTO products (title, description, price, category_id) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id
    `

	args := []interface{}{product.Title, product.Description, product.Price, product.CategoryID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pm.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (pm ProductModel) Get(id int) (*Product, error) {
	query := `
        SELECT id, title, description, price, category_id
        FROM products
        WHERE id = $1
    `

	var product Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pm.DB.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.CategoryID)
	if err != nil {

		return nil, err
	}

	return &product, nil
}

func (pm ProductModel) Update(product *Product) error {
	query := `
        UPDATE products
        SET title = $2, description = $3, price = $4, category_id = $5
        WHERE id = $1
    `

	args := []interface{}{product.ID, product.Title, product.Description, product.Price, product.CategoryID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := pm.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pm ProductModel) Delete(id int) error {
	query := `
        DELETE FROM products
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := pm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
