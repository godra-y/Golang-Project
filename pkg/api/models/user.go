package models

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

func (um UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (username, email, password_hash) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `

	args := []interface{}{user.Username, user.Email, user.PasswordHash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (um UserModel) Get(id int) (*User, error) {
	query := `
        SELECT id, username, email, password_hash
        FROM users
        WHERE id = $1
    `

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (um UserModel) Update(user *User) error {
	query := `
        UPDATE users
        SET username = $2, email = $3, password_hash = $4
        WHERE id = $1
    `

	args := []interface{}{user.ID, user.Username, user.Email, user.PasswordHash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := um.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (um UserModel) Delete(id int) error {
	query := `
        DELETE FROM users
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := um.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
