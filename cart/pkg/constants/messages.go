package constants

const (
	ErrInvalidUserID = "invalid user_id value"
	ErrInvalidSKUID  = "invalid sku_id value"
	ErrUnmarshal     = "failed to unmarshal body"
	ErrMarshal       = "failed to marshal body"
	ErrReadBody      = "failed to read body"
	ErrBadCount      = "invalid amount of products"
	ErrItemNotFound  = "item not found"
	ErrGetItem       = "failed to get item"
	ErrGetItems      = "failed to get items"
	ErrAddItem       = "failed to add item"
	ErrEmptyCart     = "cart is empty or doesn't exist"
	ErrCartCheckout  = "failed to checkout cart"
	ErrRemoveItem    = "failed to remove item"
	ErrRemoveCart    = "failed to remove cart"
	ErrInternalError = "internal server error"
)
