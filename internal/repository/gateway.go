//go:generate mockgen -source gateway.go -destination mocks/gateway.go -package mocks
package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"payment-gateway/internal/models"
)

type GatewayRepository interface {
	GetAvailableGateways(countryID int) ([]models.Gateway, error)
	CreateGateway(gateway models.Gateway) error
	GetGateways() ([]models.Gateway, error)
}

type gatewayRepository struct {
	db *sql.DB
}

func NewGatewayRepository(db *sql.DB) GatewayRepository {
	return &gatewayRepository{
		db: db,
	}
}

func (r *gatewayRepository) GetAvailableGateways(countryID int) ([]models.Gateway, error) {
	query := `
		SELECT g.id, 
		       g.name, 
		       g.data_format_supported, 
		       g.priority
		FROM gateways g
		JOIN gateway_countries gc ON g.id = gc.gateway_id
		WHERE gc.country_id = $1 AND g.status = 'active'
		ORDER BY g.priority ASC
	`
	rows, err := r.db.Query(query, countryID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gateway: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var gateways []models.Gateway
	for rows.Next() {
		var gateway models.Gateway
		if err := rows.Scan(&gateway.ID, &gateway.Name, &gateway.DataFormatSupported, &gateway.Priority); err != nil {
			return nil, fmt.Errorf("failed to scan gateway: %v", err)
		}
		gateways = append(gateways, gateway)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning rows: %v", err)
	}

	return gateways, nil
}

func (r *gatewayRepository) CreateGateway(gateway models.Gateway) error {
	query := `INSERT INTO gateway (name, data_format_supported, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRow(query, gateway.Name, gateway.DataFormatSupported, time.Now(), time.Now()).Scan(&gateway.ID)
	if err != nil {
		return fmt.Errorf("failed to insert gateway: %v", err)
	}
	return nil
}

func (r *gatewayRepository) GetGateways() ([]models.Gateway, error) {
	rows, err := r.db.Query(`SELECT id, name, data_format_supported, created_at, updated_at FROM gateway`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gateway: %v", err)
	}
	defer rows.Close()

	var gateways []models.Gateway
	for rows.Next() {
		var gateway models.Gateway
		if err := rows.Scan(&gateway.ID, &gateway.Name, &gateway.DataFormatSupported, &gateway.CreatedAt, &gateway.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan gateway: %v", err)
		}
		gateways = append(gateways, gateway)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return gateways, nil
}
