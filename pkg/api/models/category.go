package models

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

func (cm CategoryModel) Insert(category *Category) error {
	query := `
		INSERT INTO categories (name) 
		VALUES ($1) 
		RETURNING id
	`

	args := []interface{}{category.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := cm.DB.QueryRowContext(ctx, query, args...).Scan(&category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (cm CategoryModel) Get(id int) (*Category, error) {
	query := `
		SELECT id, name
		FROM categories
		WHERE id = $1
	`

	var category Category

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := cm.DB.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (сm CategoryModel) Update(category *Category) error {
	query := `
		UPDATE categories
		SET name = $2
		WHERE id = $1
		RETURNING id
	`

	args := []interface{}{category.ID, category.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := сm.DB.QueryRowContext(ctx, query, args...).Scan(&category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (сm CategoryModel) Delete(id int) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := сm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
