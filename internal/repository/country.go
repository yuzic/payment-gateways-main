//go:generate mockgen -source country.go -destination mocks/country.go -package mocks
package repository

import (
	"database/sql"
	"fmt"
	"time"

	"payment-gateway/internal/models"
)

type CountryRepository interface{}

type countryRepository struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) CountryRepository {
	return &countryRepository{
		db: db,
	}
}

func (r *countryRepository) CreateCountry(db *sql.DB, country models.Country) error {
	query := `INSERT INTO countries (name, code, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := db.QueryRow(query, country.Name, country.Code, time.Now(), time.Now()).Scan(&country.ID)
	if err != nil {
		return fmt.Errorf("failed to insert country: %w", err)
	}
	return nil
}

func (r *countryRepository) GetCountries(db *sql.DB) ([]models.Country, error) {
	rows, err := db.Query(`SELECT id, name, code, created_at, updated_at FROM countries`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(&country.ID, &country.Name, &country.Code, &country.CreatedAt, &country.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan country: %v", err)
		}
		countries = append(countries, country)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *countryRepository) GetSupportedCountriesByGateway(db *sql.DB, gatewayID int) ([]models.Country, error) {
	query := `
		SELECT c.id AS country_id, c.name AS country_name
		FROM countries c
		JOIN gateway_countries gc ON c.id = gc.country_id
		WHERE gc.gateway_id = $1
		ORDER BY c.name
	`

	rows, err := db.Query(query, gatewayID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch countries for gateway %d: %v", gatewayID, err)
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(&country.ID, &country.Name); err != nil {
			return nil, fmt.Errorf("failed to scan country: %v", err)
		}
		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return countries, nil
}
