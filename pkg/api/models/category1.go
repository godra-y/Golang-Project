package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *CategoryModel) Insert(category *Category) (int, error) {
	query := `
        INSERT INTO categories (name) 
        VALUES ($1) 
        RETURNING id
    `
	args := []interface{}{category.Name}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m CategoryModel) Get(id int) (*Category, error) {
	query := `
		SELECT id, name
		FROM categories
		WHERE id = $1
	`
	var category Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (m CategoryModel) Update(category *Category) error {
	query := `
		UPDATE categories
		SET name = $1
		WHERE id = $2
		RETURNING id
	`
	args := []interface{}{category.Name, category.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&category.ID)
}

func (m CategoryModel) Delete(id int) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
