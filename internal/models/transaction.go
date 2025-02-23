package models

import "time"

const (
	TransactionStatusPending = "pending"
	TransactionStatusFailed  = "failed"
	TransactionStatusDone    = "done"
)

const (
	TransactionTypeDeposit    = "deposit"
	TransactionTypeWithdrawal = "withdrawal"
)

type Transaction struct {
	ID        int
	Amount    float64
	Type      string
	Status    string
	UserID    int
	GatewayID int
	CountryID int
	CreatedAt time.Time
}
