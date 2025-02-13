package models

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrItemNotFound        = errors.New("SKU not found")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrReservationConflict = errors.New("reservation conflict")
	ErrOrderStatusConflict = errors.New("operation is not allowed for order with current status")
)
