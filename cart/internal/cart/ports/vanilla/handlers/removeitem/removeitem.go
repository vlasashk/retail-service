package removeitem

import "context"

type CartDeleter interface {
	DeleteItem(ctx context.Context, userID, skuID int64) error
}
