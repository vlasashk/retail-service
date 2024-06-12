package models

type Order struct {
	UserID int64
	Items  []Item
}
