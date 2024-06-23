package pgorders

import (
	"time"

	"route256/loms/internal/loms/adapters/pgorders/pgordersqry"
	"route256/loms/internal/loms/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func statusToDTO(status models.OrderStatus) pgordersqry.OrdersOrderStatus {
	switch status {
	case models.NewStatus:
		return pgordersqry.OrdersOrderStatusNew
	case models.AwaitingPaymentStatus:
		return pgordersqry.OrdersOrderStatusAwaitingPayment
	case models.PayedStatus:
		return pgordersqry.OrdersOrderStatusPayed
	case models.CancelledStatus:
		return pgordersqry.OrdersOrderStatusCancelled
	case models.FailedStatus:
		return pgordersqry.OrdersOrderStatusFailed
	default:
		return pgordersqry.OrdersOrderStatusUnknown
	}
}

func statusToDomain(status pgordersqry.OrdersOrderStatus) models.OrderStatus {
	switch status {
	case pgordersqry.OrdersOrderStatusNew:
		return models.NewStatus
	case pgordersqry.OrdersOrderStatusAwaitingPayment:
		return models.AwaitingPaymentStatus
	case pgordersqry.OrdersOrderStatusPayed:
		return models.PayedStatus
	case pgordersqry.OrdersOrderStatusCancelled:
		return models.CancelledStatus
	case pgordersqry.OrdersOrderStatusFailed:
		return models.FailedStatus
	default:
		return models.UnknownStatus
	}
}

func orderToDTO(order models.Order) pgordersqry.OrdersOrder {
	return pgordersqry.OrdersOrder{
		UserID: order.UserID,
		Status: statusToDTO(order.Status),
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

func orderToDomain(order pgordersqry.GetOrderByIdRow, items []pgordersqry.GetOrderItemsRow) models.Order {
	return models.Order{
		UserID: order.UserID,
		Items:  itemsToDomain(items),
		Status: statusToDomain(order.Status),
	}
}

func itemsToDTO(orderID int64, items []models.Item) []pgordersqry.InsertOrderItemsParams {
	itemsDTO := make([]pgordersqry.InsertOrderItemsParams, 0, len(items))

	for _, item := range items {
		itemsDTO = append(itemsDTO, pgordersqry.InsertOrderItemsParams{
			SkuID:   int64(item.SKUid),
			OrderID: orderID,
			Count:   int64(item.Count),
			CreatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		})
	}

	return itemsDTO
}

func itemsToDomain(itemsDTO []pgordersqry.GetOrderItemsRow) []models.Item {
	items := make([]models.Item, 0, len(itemsDTO))

	for _, item := range itemsDTO {
		items = append(items, models.Item{
			SKUid: uint32(item.SkuID),
			Count: uint32(item.Count),
		})
	}

	return items
}
