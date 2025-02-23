package models

import "time"

type Country struct {
	ID        int
	Name      string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
