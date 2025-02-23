package models

import "time"

type Gateway struct {
	ID                  int
	Name                string
	DataFormatSupported string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Priority            int
	Status              string
}
