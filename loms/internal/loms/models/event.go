package models

import "time"

type Event struct {
	ID        int64
	OrderID   int64
	Status    string
	CreatedAt time.Time
}
