package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"payment-gateway/internal/models"
	"strconv"
	"testing"
)

// MockTransactionService implements TransactionService for testing
type MockTransactionService struct {
	shouldFailDeposit    bool
	shouldFailWithdrawal bool
	shouldFailUpdate     bool
	lastTransactionID    int
	lastGatewayID        int64
	lastStatus           string
}

func (m *MockTransactionService) Deposit(req models.TransactionRequest) (*models.Transaction, error) {
	if m.shouldFailDeposit {
		return nil, errors.New("deposit failed")
	}
	return &models.Transaction{ID: 123, Status: "processed"}, nil
}

func (m *MockTransactionService) Withdrawal(req models.TransactionRequest) (*models.Transaction, error) {
	if m.shouldFailWithdrawal {
		return nil, errors.New("withdrawal failed")
	}
	return &models.Transaction{ID: 456, Status: "processed"}, nil
}

func (m *MockTransactionService) UpdateStatus(txID int, gatewayID int64, status string) error {
	if m.shouldFailUpdate {
		return errors.New("update failed")
	}
	m.lastTransactionID = txID
	m.lastGatewayID = gatewayID
	m.lastStatus = status
	return nil
}

func TestDepositHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		serviceFail    bool
		wantStatusCode int
	}{
		{
			name:           "successful deposit",
			requestBody:    `{"amount":100.00,"user_id":1,"currency":"EUR"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "invalid request body",
			requestBody:    `invalid json`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "service failure",
			requestBody:    `{"amount":100.00,"user_id":1,"currency":"EUR"}`,
			serviceFail:    true,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockTransactionService{shouldFailDeposit: tt.serviceFail}
			handler := NewHandler(mockService)

			req := httptest.NewRequest("POST", "/deposit", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.DepositHandler(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}
		})
	}
}

func TestWithdrawalHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		serviceFail    bool
		wantStatusCode int
	}{
		{
			name:           "successful withdrawal",
			requestBody:    `{"amount":50.00,"user_id":1,"currency":"USD"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "invalid request body",
			requestBody:    `invalid json`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "service failure",
			requestBody:    `{"amount":50.00,"user_id":1,"currency":"USD"}`,
			serviceFail:    true,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockTransactionService{shouldFailWithdrawal: tt.serviceFail}
			handler := NewHandler(mockService)

			req := httptest.NewRequest("POST", "/withdraw", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.WithdrawalHandler(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}
		})
	}
}

func TestCallbackHandler(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		serviceFail    bool
		wantStatusCode int
	}{
		{
			name: "successful callback",
			queryParams: map[string]string{
				"id":      "123",
				"status":  "completed",
				"gateway": "456",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "missing parameters",
			queryParams: map[string]string{
				"id": "123",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid transaction ID",
			queryParams: map[string]string{
				"id":      "invalid",
				"status":  "completed",
				"gateway": "456",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "service update failure",
			queryParams: map[string]string{
				"id":      "123",
				"status":  "completed",
				"gateway": "456",
			},
			serviceFail:    true,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockTransactionService{shouldFailUpdate: tt.serviceFail}
			handler := NewHandler(mockService)

			req := httptest.NewRequest("GET", "/callback", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handler.CallbackHandler(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}

			// Verify service call parameters for successful cases
			if tt.wantStatusCode == http.StatusOK {
				expectedID, _ := strconv.ParseInt(tt.queryParams["id"], 10, 64)
				expectedGateway, _ := strconv.ParseInt(tt.queryParams["gateway"], 10, 64)

				if mockService.lastTransactionID != int(expectedID) ||
					mockService.lastGatewayID != expectedGateway ||
					mockService.lastStatus != tt.queryParams["status"] {
					t.Errorf("service called with unexpected parameters: got %d/%d/%s want %d/%d/%s",
						mockService.lastTransactionID,
						mockService.lastGatewayID,
						mockService.lastStatus,
						expectedID,
						expectedGateway,
						tt.queryParams["status"],
					)
				}
			}
		})
	}
}
