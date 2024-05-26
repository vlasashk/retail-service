package common

import (
	"context"

	"route256/cart/internal/cart/models"
)

type ProductProvider interface {
	GetProduct(ctx context.Context, sku int64) (models.ItemDescription, error)
}
