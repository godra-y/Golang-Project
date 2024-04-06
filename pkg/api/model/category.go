package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/godra-y/go-project/pkg/api/validator"
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

func (cm CategoryModel) GetAll(name string, filters Filters) ([]*Category, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, name
        FROM categories
        WHERE (LOWER(name) = LOWER($1) OR $1 = '')
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3
    `, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, filters.limit(), filters.offset()}

	rows, err := cm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			cm.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var categories []*Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&totalRecords, &category.ID, &category.Name)
		if err != nil {
			return nil, Metadata{}, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {

		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return categories, metadata, nil
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

	return cm.DB.QueryRowContext(ctx, query, args...).Scan(&category.ID)
}

func (cm CategoryModel) Get(id int) (*Category, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name
		FROM categories
		WHERE id = $1
	`

	var category Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := cm.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve category with id: %v, %w", id, err)
	}
	return &category, nil
}

func (cm CategoryModel) Update(category *Category) error {
	query := `
		UPDATE categories
		SET name = $1
		WHERE id = $2
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return cm.DB.QueryRowContext(ctx, query, category.Name, category.ID).Scan(&category.ID)
}

func (сm CategoryModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := сm.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 100, "name", "must not be more than 100 bytes long")
}
