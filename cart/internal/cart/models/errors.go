package models

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrItemNotFound   = errors.New("item not found")
	ErrItemProvider   = errors.New("failed to request item info")
	ErrAddItem        = errors.New("failed to add item")
	ErrCartIsEmpty    = errors.New("cart is empty or doesn't exist")
	ErrCartCheckout   = errors.New("failed to checkout cart")
	ErrBadProductData = errors.New("product data is invalid")
	ErrReadBody       = errors.New("failed to read body")
	ErrJSONProcessing = errors.New("failed to process request body")
	ErrBadCount       = errors.New("invalid amount of products")
	ErrInvalidUserID  = errors.New("invalid user_id value")
	ErrInvalidSKUID   = errors.New("invalid sku_id value")
	ErrRemoveItem     = errors.New("failed to remove item")
	ErrRemoveCart     = errors.New("failed to remove cart")
	ErrInternalError  = errors.New("internal server error")
	ErrBadProductInfo = errors.New("bad product info")
)
