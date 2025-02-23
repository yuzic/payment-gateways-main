//go:generate mockgen -source transaction.go -destination mocks/transaction.go -package mocks

package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/models"
	"payment-gateway/internal/repository"
	"payment-gateway/internal/services/gateway"
	"payment-gateway/internal/util"
)

type transactionService struct {
	gateway   gateway.ServiceGateway
	userRepo  repository.UserRepository
	transRepo repository.TransactionRepository
	publisher kafka.KafkaPublisher
}

type TransactionService interface {
	Deposit(req models.TransactionRequest) (*models.Transaction, error)
	Withdrawal(req models.TransactionRequest) (*models.Transaction, error)
	UpdateStatus(txID int, gatewayID int64, status string) error
}

const (
	amountErr  = "invalid amount, must be greater than zero"
	userErr    = "invalid user"
	maxRetries = 3
)

func NewTransactionService(
	gw gateway.ServiceGateway,
	userRepo repository.UserRepository,
	transRepo repository.TransactionRepository,
	kafkaPublisher kafka.KafkaPublisher,

) TransactionService {
	return &transactionService{
		gateway:   gw,
		userRepo:  userRepo,
		transRepo: transRepo,
		publisher: kafkaPublisher,
	}
}

func (s *transactionService) Deposit(req models.TransactionRequest) (*models.Transaction, error) {
	if err := s.validateTransaction(req); err != nil {
		return nil, err
	}

	// some different business logic

	return s.transaction(req, models.TransactionTypeDeposit)
}

func (s *transactionService) Withdrawal(req models.TransactionRequest) (*models.Transaction, error) {
	if err := s.validateTransaction(req); err != nil {
		return nil, err
	}

	// some different business logic
	return s.transaction(req, models.TransactionTypeWithdrawal)
}

func (s *transactionService) transaction(req models.TransactionRequest, transactionType string) (*models.Transaction, error) {
	user, err := s.userRepo.GetUserByID(req.UserID)
	if err != nil {
		log.Printf("Error db.GetUserByID: %v", err)
		return nil, err
	}

	gateway, err := s.gateway.GetGateway(req.UserID)
	if err != nil {
		return nil, err
	}

	tx := models.Transaction{
		UserID:    user.ID,
		Amount:    req.Amount,
		GatewayID: gateway.ID,
		CountryID: user.CountryID,
		Status:    models.TransactionStatusPending,
		Type:      transactionType,
	}

	tx.ID, err = s.transRepo.CreateTransaction(tx)
	if err != nil {
		log.Printf("Error db.CreateTransaction: %v", err)
		return nil, err
	}

	err = s.gateway.Deposit(tx)
	if err != nil {
		log.Printf("Error gateway.Deposit: %v", err)
		return nil, err
	}

	if err = util.RetryOperation(func() error {
		err = s.gateway.Deposit(tx)
		return err
	}, maxRetries); err != nil {

		if err = s.transRepo.UpdateStatus(tx.ID, models.TransactionStatusFailed); err != nil {
			return nil, errors.New("error s.transRepo.UpdateStatus")
		}

		return nil, errors.New(" util.RetryOperation")
	}

	txByte, err := json.Marshal(tx)
	if err != nil {
		log.Printf("Error json.Marshal: %v", err)
		return nil, err
	}

	err = s.publisher.PublishTransaction(
		context.Background(),
		"txn-12345",
		txByte,
		"application/json",
	)

	if err != nil {
		log.Printf("Failed to publish transaction: %v", err)
		return nil, err
	}

	return &tx, nil
}

func (s *transactionService) UpdateStatus(txID int, gatewayID int64, status string) error {
	statusTx := models.TransactionStatusPending
	// check status from external gateways
	switch status {
	case models.TransactionStatusPending:
		statusTx = status
	case models.TransactionStatusFailed:
		statusTx = status
	case models.TransactionStatusDone:
		statusTx = status
	default:
		statusTx = models.TransactionStatusPending
	}

	_, err := s.transRepo.GetTransaction(txID)
	if err != nil {
		log.Printf("Error db.GetTransaction: %v", err)
		return err
	}

	return s.transRepo.UpdateStatus(txID, statusTx)
}

func (s *transactionService) validateTransaction(req models.TransactionRequest) error {
	if req.Amount <= 0 {
		return errors.New(amountErr)
	}

	if req.UserID <= 0 {
		return errors.New(userErr)
	}

	return nil
}
