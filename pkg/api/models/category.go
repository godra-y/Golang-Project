package models

import (
	"database/sql"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (cm *CategoryModel) Insert(category *Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	err := cm.DB.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (cm *CategoryModel) Get(id int) (*Category, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`
	var category Category
	err := cm.DB.QueryRow(query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (cm *CategoryModel) Update(category *Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err := cm.DB.Exec(query, category.Name, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (cm *CategoryModel) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := cm.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
