package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/godra-y/go-project/pkg/api/validator"
	"log"
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
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (pm ProductModel) GetAll(title string, price int, id int, filters Filters) ([]*Product, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, title, description, price, category_id
		FROM products
		WHERE (LOWER(title) = LOWER($1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, filters.limit(), filters.offset()}

	rows, err := pm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			pm.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var products []*Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&totalRecords, &product.ID, &product.Title, &product.Description, &product.Price, &product.CategoryID)
		if err != nil {
			return nil, Metadata{}, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return products, metadata, nil
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

	return pm.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
}

func (pm ProductModel) Get(id int) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, description, price, category_id
		FROM products
		WHERE id = $1
		`

	var product Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := pm.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve product with id: %v, %w", id, err)
	}
	return &product, nil
}

func (pm ProductModel) GetProductsByCategory(categoryID int, title string, filters Filters) ([]*Product, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, title, description, price, category_id
		FROM products
		WHERE category_id = $1 AND (LOWER(title) = LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{categoryID, title, filters.limit(), filters.offset()}

	rows, err := pm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			pm.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var products []*Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&totalRecords, &product.ID, &product.Title, &product.Description, &product.Price, &product.CategoryID)
		if err != nil {
			return nil, Metadata{}, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return products, metadata, nil
}

func (pm ProductModel) Update(product *Product) error {
	query := `
		UPDATE products
		SET title = $1, description = $2, price = $3, category_id = $4
		WHERE id = $5
		RETURNING id
		`

	args := []interface{}{product.Title, product.Description, product.Price, product.CategoryID, product.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return pm.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
}

func (pm ProductModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM products
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := pm.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Title != "", "title", "must be provided")
	v.Check(len(product.Title) <= 100, "title", "must not be more than 100 bytes long")
	v.Check(len(product.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	v.Check(product.Price >= 0, "price", "must be a non-negative value")
}
