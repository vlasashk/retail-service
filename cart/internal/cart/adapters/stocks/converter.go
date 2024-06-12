package stocks

import (
	"route256/cart/internal/cart/models"
	lomsservicev1 "route256/cart/pkg/api/loms/v1"
)

func orderToDTO(order models.Order) *lomsservicev1.OrderCreateRequest {
	items := make([]*lomsservicev1.Item, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, &lomsservicev1.Item{
			Sku:   uint32(item.SkuID),
			Count: uint32(item.Count),
		})
	}

	return &lomsservicev1.OrderCreateRequest{
		User:  order.UserID,
		Items: items,
	}
}

func skuIDToDTO(skuID int64) *lomsservicev1.StocksInfoRequest {
	return &lomsservicev1.StocksInfoRequest{
		Sku: uint32(skuID),
	}
}
