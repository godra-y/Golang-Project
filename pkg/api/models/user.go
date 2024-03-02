package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

//crud

func (m UserModel) CreateUser(user *User) (int, error) {
	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`
	args := []interface{}{user.Username, user.Email, user.PasswordHash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
}

func (m UserModel) GetUserByID(id int) (*User, error) {
	query := `
		SELECT id, username, email, password_hash
		FROM users
		WHERE id = $1
	`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m UserModel) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3
		WHERE id = $4
		RETURNING id
	`
	args := []interface{}{user.Username, user.Email, user.PasswordHash, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
}

func (m UserModel) DeleteUser(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
