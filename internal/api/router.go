package api

import (
	"database/sql"
	"net/http"
	"payment-gateway/internal/kafka"
	repo "payment-gateway/internal/repository"
	"payment-gateway/internal/services/gateway"
	"payment-gateway/internal/services/transaction"

	"github.com/gorilla/mux"
)

type DiContainer struct {
	handler *Handler
}

func GetContainer(db *sql.DB, kf kafka.KafkaPublisher) *DiContainer {
	gatewayRepo := repo.NewGatewayRepository(db)
	userRepo := repo.NewUserRepository(db)
	transRepo := repo.NewTransactionRepository(db)

	gatewayService := gateway.NewServiceGateway(gatewayRepo)

	transactionService := transaction.NewTransactionService(gatewayService, userRepo, transRepo, kf)

	handler := NewHandler(transactionService)

	return &DiContainer{
		handler: handler,
	}

}

func SetupRouter(di *DiContainer) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/deposit", http.HandlerFunc(di.handler.DepositHandler)).Methods("POST")
	router.Handle("/withdrawal", http.HandlerFunc(di.handler.WithdrawalHandler)).Methods("POST")
	router.Handle("/callback", http.HandlerFunc(di.handler.CallbackHandler)).Methods("GET")

	return router
}
