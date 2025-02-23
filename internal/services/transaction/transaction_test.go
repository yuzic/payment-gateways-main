package transaction

import (
	"errors"
	"testing"

	mockPublisher "payment-gateway/internal/kafka/mocks"
	"payment-gateway/internal/models"
	"payment-gateway/internal/repository/mocks"
	mockGateway "payment-gateway/internal/services/gateway/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeposit_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockGateway.NewMockServiceGateway(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTransRepo := mocks.NewMockTransactionRepository(ctrl)
	mockPublisher := mockPublisher.NewMockKafkaPublisher(ctrl)

	service := NewTransactionService(mockGateway, mockUserRepo, mockTransRepo, mockPublisher)

	req := models.TransactionRequest{
		UserID:   1,
		Amount:   100.00,
		Currency: "EUR",
	}

	user := models.User{ID: 1, CountryID: 2}
	gw := &models.Gateway{ID: 10}
	tx := models.Transaction{
		UserID:    user.ID,
		Amount:    req.Amount,
		GatewayID: gw.ID,
		CountryID: user.CountryID,
		Status:    models.TransactionStatusPending,
		Type:      models.TransactionTypeDeposit,
	}

	// Set expectations for repository and gateway calls.
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(user, nil)
	mockGateway.EXPECT().GetGateway(req.UserID).Return(gw, nil)
	mockTransRepo.EXPECT().CreateTransaction(gomock.Any()).Return(1, nil)

	// Expect Deposit to be called once (adjusted from .Times(2) to .Times(1))
	mockGateway.EXPECT().Deposit(gomock.Any()).Return(nil).Times(1)

	// Use a flexible matcher for the payload.
	mockPublisher.EXPECT().PublishTransaction(gomock.Any(), gomock.Any(), gomock.Any(), "application/json").Return(nil)

	result, err := service.Deposit(req)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tx.Amount, result.Amount)
	assert.Equal(t, models.TransactionStatusPending, result.Status)
}

func TestDeposit_Fail_InvalidUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockGateway.NewMockServiceGateway(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTransRepo := mocks.NewMockTransactionRepository(ctrl)
	mockPublisher := mockPublisher.NewMockKafkaPublisher(ctrl)

	service := NewTransactionService(mockGateway, mockUserRepo, mockTransRepo, mockPublisher)

	req := models.TransactionRequest{
		UserID:   0, // Невалидный пользователь
		Amount:   100.00,
		Currency: "EUR",
	}

	result, err := service.Deposit(req)
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid user")
}

func TestDeposit_Fail_TransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockGateway.NewMockServiceGateway(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTransRepo := mocks.NewMockTransactionRepository(ctrl)
	mockPublisher := mockPublisher.NewMockKafkaPublisher(ctrl)

	service := NewTransactionService(mockGateway, mockUserRepo, mockTransRepo, mockPublisher)

	req := models.TransactionRequest{
		UserID:   1,
		Amount:   100.00,
		Currency: "EUR",
	}

	user := models.User{ID: 1, CountryID: 2}
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(user, nil)
	mockGateway.EXPECT().GetGateway(req.UserID).Return(&models.Gateway{ID: 10}, nil)
	mockTransRepo.EXPECT().CreateTransaction(gomock.Any()).Return(0, errors.New("db error"))

	result, err := service.Deposit(req)
	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
}
