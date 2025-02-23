//go:generate mockgen -source user.go -destination mocks/user.go -package mocks
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"payment-gateway/internal/models"
	"time"
)

type UserRepository interface {
	CreateUser(user models.User) error
	GetUserByID(userID int) (models.User, error)
	GetUsers() ([]models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user models.User) error {
	query := `INSERT INTO users (username, email, country_id, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRow(query, user.Username, user.Email, user.CountryID, time.Now(), time.Now()).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

// GetUserByID Get use by id
func (r *userRepository) GetUserByID(userID int) (models.User, error) {
	var user models.User

	query := `SELECT 
    			id, 
    			username, 
    			email, 
    			country_id, 
    			created_at, 
    			updated_at 
			  FROM users WHERE id = $1`

	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CountryID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("no user found with id %d", userID)
		}
		return models.User{}, fmt.Errorf("failed to fetch user: %v", err)
	}

	return user, nil
}

func (r *userRepository) GetUsers() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT id, username, email, country_id, created_at, updated_at FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CountryID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
