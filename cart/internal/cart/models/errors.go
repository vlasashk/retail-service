package models

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrItemProvider   = errors.New("failed to request item info")
	ErrCartIsEmpty    = errors.New("cart is empty")
	ErrBadProductData = errors.New("product data is invalid")
)
