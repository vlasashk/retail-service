package pgorders

import (
	"time"

	"route256/loms/internal/loms/models"
)

const (
	unknownStatus         string = "Unknown"
	newStatus             string = "New"
	awaitingPaymentStatus string = "AwaitingPayment"
	payedStatus           string = "Payed"
	cancelledStatus       string = "Cancelled"
	failedStatus          string = "Failed"
)

type orderDTO struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type itemDTO struct {
	SKU       uint32    `db:"sku_id"`
	OrderID   int64     `db:"order_id"`
	Count     uint32    `db:"count"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func statusToDTO(status models.OrderStatus) string {
	switch status {
	case models.NewStatus:
		return newStatus
	case models.AwaitingPaymentStatus:
		return awaitingPaymentStatus
	case models.PayedStatus:
		return payedStatus
	case models.CancelledStatus:
		return cancelledStatus
	case models.FailedStatus:
		return failedStatus
	default:
		return unknownStatus
	}
}

func statusToDomain(status string) models.OrderStatus {
	switch status {
	case newStatus:
		return models.NewStatus
	case awaitingPaymentStatus:
		return models.AwaitingPaymentStatus
	case payedStatus:
		return models.PayedStatus
	case cancelledStatus:
		return models.CancelledStatus
	case failedStatus:
		return models.FailedStatus
	default:
		return models.UnknownStatus
	}
}

func orderToDTO(order models.Order) orderDTO {
	return orderDTO{
		UserID:    order.UserID,
		Status:    statusToDTO(order.Status),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func orderToDomain(order orderDTO, items []itemDTO) models.Order {
	return models.Order{
		UserID: order.UserID,
		Items:  itemsToDomain(items),
		Status: statusToDomain(order.Status),
	}
}

func itemsToDTO(orderID int64, items []models.Item) []itemDTO {
	itemsDTO := make([]itemDTO, 0, len(items))

	for _, item := range items {
		itemsDTO = append(itemsDTO, itemDTO{
			SKU:       item.SKUid,
			OrderID:   orderID,
			Count:     item.Count,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return itemsDTO
}

func itemsToDomain(itemsDTO []itemDTO) []models.Item {
	items := make([]models.Item, 0, len(itemsDTO))

	for _, item := range itemsDTO {
		items = append(items, models.Item{
			SKUid: item.SKU,
			Count: item.Count,
		})
	}

	return items
}
