package impl

import (
	"route256/loms/internal/loms/models"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"
)

func ItemToDomain(item *lomsservicev1.Item) models.Item {
	return models.Item{
		SKUid: item.Sku,
		Count: item.Count,
	}
}

func ItemToProto(item models.Item) *lomsservicev1.Item {
	return &lomsservicev1.Item{
		Sku:   item.SKUid,
		Count: item.Count,
	}
}

func OrderToDomain(req *lomsservicev1.OrderCreateRequest) models.Order {
	items := make([]models.Item, len(req.Items))
	for i, item := range req.Items {
		items[i] = ItemToDomain(item)
	}
	return models.Order{
		UserID: req.User,
		Items:  items,
		Status: models.UnknownStatus,
	}
}

func OrderToProto(order models.Order) *lomsservicev1.OrderInfoResponse {
	items := make([]*lomsservicev1.Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = ItemToProto(item)
	}
	return &lomsservicev1.OrderInfoResponse{
		User:   order.UserID,
		Items:  items,
		Status: StatusToProto(order.Status),
	}
}

func StatusToProto(status models.OrderStatus) lomsservicev1.Status {
	switch status {
	case models.NewStatus:
		return lomsservicev1.Status_STATUS_NEW
	case models.AwaitingPaymentStatus:
		return lomsservicev1.Status_STATUS_AWAITING_PAYMENT
	case models.PayedStatus:
		return lomsservicev1.Status_STATUS_PAYED
	case models.CancelledStatus:
		return lomsservicev1.Status_STATUS_CANCELLED
	case models.FailedStatus:
		return lomsservicev1.Status_STATUS_FAILED
	default:
		return lomsservicev1.Status_STATUS_UNKNOWN
	}
}
