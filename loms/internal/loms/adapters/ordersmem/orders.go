package ordersmem

import (
	"context"
	"math/rand/v2"
	"sync"
	"time"

	"route256/loms/internal/loms/models"
)

type userOrder struct {
	models.Order
	createdAt time.Time
}

type Orders struct {
	mu     sync.RWMutex
	orders map[int64]userOrder
}

func New() *Orders {
	return &Orders{
		mu:     sync.RWMutex{},
		orders: make(map[int64]userOrder),
	}
}

func (o *Orders) Create(_ context.Context, order models.Order) (int64, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	orderID := rand.Int64()

	for _, exist := o.orders[orderID]; exist; {
		orderID = rand.Int64()
	}

	orderToCreat := userOrder{
		Order:     order,
		createdAt: time.Now(),
	}
	orderToCreat.Status = models.NewStatus

	o.orders[orderID] = orderToCreat

	return orderID, nil
}

func (o *Orders) SetStatus(_ context.Context, orderID int64, status models.OrderStatus) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	order, ok := o.orders[orderID]
	if !ok {
		return models.ErrOrderNotFound
	}

	order.Status = status
	o.orders[orderID] = order

	return nil
}

func (o *Orders) GetByOrderID(_ context.Context, orderID int64) (models.Order, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	order, ok := o.orders[orderID]
	if !ok {
		return models.Order{}, models.ErrOrderNotFound
	}

	return order.Order, nil
}
