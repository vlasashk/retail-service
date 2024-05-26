package models

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrCartIsEmpty = errors.New("cart is empty")
)
