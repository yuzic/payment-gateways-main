package models

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	CountryID int
	CreatedAt time.Time
	UpdatedAt time.Time
}
