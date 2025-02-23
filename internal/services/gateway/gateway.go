//go:generate mockgen -source gateway.go -destination mocks/gateway.go -package mocks

package gateway

import (
	"errors"
	"fmt"
	"log"
	"payment-gateway/internal/models"
	"payment-gateway/internal/repository"
	"payment-gateway/internal/util"
)

const (
	GatewayError   = "No available gateway"
	gatewayErrPing = "Gateways are unhealthy/unavailable"
)

type ServiceGateway interface {
	GetGateway(countryID int) (*models.Gateway, error)
	Deposit(req models.Transaction) error
	Withdrawal(req models.Transaction) error
}

type serviceGateway struct {
	gatewayRepo repository.GatewayRepository
}

func NewServiceGateway(gatewayRepo repository.GatewayRepository) ServiceGateway {
	return &serviceGateway{
		gatewayRepo: gatewayRepo,
	}
}

func (s *serviceGateway) GetGateway(countryID int) (*models.Gateway, error) {
	gateways, err := s.gatewayRepo.GetAvailableGateways(countryID)
	if err != nil || len(gateways) == 0 {
		log.Printf("Error repo.GetAvailableGateways: %v", err)
		return nil, errors.New(GatewayError)
	}

	for _, gateway := range gateways {
		if s.ping() {
			return &gateway, nil
		}
	}

	return nil, errors.New(gatewayErrPing)
}

func (s *serviceGateway) Deposit(req models.Transaction) error {
	// external request to Gateway here
	amount := util.MaskData([]byte(fmt.Sprintf("%.2f", req.Amount)))

	log.Printf("Gateway Deposit succes with amount: %v", amount)

	return nil
}

func (s *serviceGateway) Withdrawal(req models.Transaction) error {
	// external request to Gateway here
	amount := util.MaskData([]byte(fmt.Sprintf("%.2f", req.Amount)))

	log.Printf("Gateway Withdrawal succes with amount: %v", amount)
	return nil
}

// ping check gateway
func (s *serviceGateway) ping() bool {
	// http request there on other call for check external gateway
	// I think that best choose will be chet it in background
	return true
}
