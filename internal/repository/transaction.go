//go:generate mockgen -source transaction.go -destination mocks/transaction.go -package mocks
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"payment-gateway/internal/models"
	"time"
)

type TransactionRepository interface {
	CreateTransaction(transaction models.Transaction) (int, error)
	GetTransactions() ([]models.Transaction, error)
	UpdateStatus(transactionID int, status string) error
	GetTransaction(transactionID int) (*models.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreateTransaction(transaction models.Transaction) (int, error) {
	query := `INSERT INTO transactions (amount, type, status, gateway_id, country_id, user_id, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRow(query, transaction.Amount, transaction.Type, transaction.Status, transaction.GatewayID, transaction.CountryID, transaction.UserID, time.Now()).Scan(&transaction.ID)
	if err != nil {
		return transaction.ID, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return transaction.ID, nil
}

func (r *transactionRepository) GetTransactions() ([]models.Transaction, error) {
	rows, err := r.db.Query(`SELECT id, amount, type, status, user_id, gateway_id, country_id, created_at FROM transactions`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Amount, &transaction.Type, &transaction.Status, &transaction.UserID, &transaction.GatewayID, &transaction.CountryID, &transaction.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) UpdateStatus(transactionID int, status string) error {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, transactionID)
	return err
}

func (r *transactionRepository) GetTransaction(transactionID int) (*models.Transaction, error) {
	query := `
        SELECT id, amount, type, status, user_id, gateway_id, country_id, created_at 
        FROM transactions 
        WHERE id = $1
    `

	var transaction models.Transaction

	err := r.db.QueryRow(query, transactionID).Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.Type,
		&transaction.Status,
		&transaction.UserID,
		&transaction.GatewayID,
		&transaction.CountryID,
		&transaction.CreatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("transaction not found with ID: %d", transactionID)
	case err != nil:
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	default:
		return &transaction, nil
	}
}
