package models

import (
	"database/sql"
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

func (pm *ProductModel) Insert(product *Product) error {
	query := `INSERT INTO products (title, description, price, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := pm.DB.QueryRow(query, product.Title, product.Description, product.Price, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (pm *ProductModel) Get(id int) (*Product, error) {
	query := `SELECT id, title, description, price, category_id FROM products WHERE id = $1`
	var product Product
	err := pm.DB.QueryRow(query, id).Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.CategoryID)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (pm *ProductModel) Update(product *Product) error {
	query := `UPDATE products SET title = $1, description = $2, price = $3, category_id = $4 WHERE id = $5`
	_, err := pm.DB.Exec(query, product.Title, product.Description, product.Price, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (pm *ProductModel) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := pm.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
