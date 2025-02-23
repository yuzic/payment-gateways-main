package api

import (
	"log"
	"net/http"
	"strconv"

	"payment-gateway/internal/models"
	"payment-gateway/internal/services/transaction"
	"payment-gateway/internal/util"
)

type Handler struct {
	transactionService transaction.TransactionService
}

func NewHandler(transactionService transaction.TransactionService) *Handler {
	return &Handler{
		transactionService: transactionService,
	}
}

// DepositHandler handles deposit requests (feel free to update how user is passed to the request)
// Sample Request (POST /deposit):
//
//	{
//	    "amount": 100.00,
//	    "user_id": 1,
//	    "currency": "EUR"
//	}
func (h *Handler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	var request models.TransactionRequest
	err := util.DecodeRequest(r, &request)
	if err != nil {
		log.Printf("Error util.DecodeRequest: %v", err)
		http.Error(w, "Error deposit", http.StatusInternalServerError)

		return
	}

	tx, err := h.transactionService.Deposit(request)
	if err != nil {
		log.Printf("Error h.TransactionService.Deposit: %v", err)
		http.Error(w, "Error deposit", http.StatusInternalServerError)
		return
	}

	err = util.EncodeResponse(w, r, models.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Transaction deposit successfully",
		Data:       newDataResp(tx),
	})

	if err != nil {
		log.Printf("Error EncodeResponse: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

}

// WithdrawalHandler handles withdrawal requests (feel free to update how user is passed to the request)
// Sample Request (POST /deposit):
//
//	{
//	    "amount": 100.00,
//	    "user_id": 1,
//	    "currency": "EUR"
//	}
func (h *Handler) WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	var request models.TransactionRequest
	err := util.DecodeRequest(r, &request)
	if err != nil {
		log.Printf("Error util.DecodeRequest: %v", err)
		http.Error(w, "Error withdrawal", http.StatusInternalServerError)

		return
	}

	tx, err := h.transactionService.Withdrawal(request)
	if err != nil {
		log.Printf("Error h.TransactionService.Withdrawal: %v", err)
		http.Error(w, "Error deposit", http.StatusInternalServerError)
		return
	}

	err = util.EncodeResponse(w, r, models.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Transaction withdrawal successfully",
		Data:       newDataResp(tx),
	})

	if err != nil {
		log.Printf("Error EncodeResponse: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// CallbackHandler handle postback query from payment system
// (GET /callback?id=101&status=done&gateway=1)
func (h *Handler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")
	gateway := r.URL.Query().Get("gateway")

	if gateway == "" || status == "" {
		log.Printf("Error Query Params: gateway=%s status=%s", gateway, status)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	txID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error Parse Transaction ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	gatewayID, err := strconv.ParseInt(gateway, 10, 64)
	if err != nil {
		log.Printf("Error Parse Transaction ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.transactionService.UpdateStatus(int(txID), gatewayID, status)
	if err != nil {
		log.Printf("Error h.TransactionService.UpdateStatus: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = util.EncodeResponse(w, r, models.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Transaction Callback successfully",
		Data:       DataResp{},
	})

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

type DataResp map[string]interface{}

func newDataResp(tx *models.Transaction) DataResp {
	return DataResp{
		"transactionID": tx.ID,
		"status":        tx.Status,
	}
}
