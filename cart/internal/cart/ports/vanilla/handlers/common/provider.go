package common

import (
	"context"

	"route256/cart/internal/cart/models"
)

//go:generate mockery --name=ProductProvider
type ProductProvider interface {
	GetProduct(ctx context.Context, sku int64) (models.ItemDescription, error)
}
