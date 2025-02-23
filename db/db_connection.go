package db

import (
	"database/sql"
	"log"

	"payment-gateway/internal/util"

	_ "github.com/lib/pq"
)

const maxRetries = 5

// InitializeDB initializes the database connection
func InitializeDB(dataSourceName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	err = util.RetryOperation(func() error {
		db, err = sql.Open("postgres", dataSourceName)
		if err != nil {
			return err
		}

		return db.Ping()
	}, maxRetries)

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
		return nil, err
	}

	return db, nil
}
