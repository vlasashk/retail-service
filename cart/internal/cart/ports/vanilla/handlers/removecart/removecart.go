package removeitem

import "context"

type CartDeleter interface {
	DeleteItemsByUserID(ctx context.Context, userID int64) error
}
