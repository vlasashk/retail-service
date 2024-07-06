package inmem

import (
	"context"
	"sync"

	"route256/cart/internal/cart/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cartCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "inmem_cart_count",
		Help: "Current number of carts in the in-memory storage",
	})
	skuCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "inmem_skus_count",
		Help: "Current number of SKUs in the in-memory storage in all carts",
	})
)

type Storage struct {
	mu    sync.RWMutex
	carts map[int64]*cart
}

func NewStorage() *Storage {
	return &Storage{
		mu:    sync.RWMutex{},
		carts: make(map[int64]*cart),
	}
}

func (s *Storage) AddItem(_ context.Context, userID, skuID int64, count uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.carts[userID]; !ok {
		s.carts[userID] = newCart()
		cartCount.Inc()
	}

	s.carts[userID].addItem(skuID, count)

	return nil
}

func (s *Storage) DeleteItem(_ context.Context, userID, skuID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.carts[userID]; !ok {
		return models.ErrCartIsEmpty
	}

	if s.carts[userID].deleteItem(skuID) {
		skuCount.Dec()
	}

	return nil
}

func (s *Storage) DeleteItemsByUserID(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.carts[userID]; !ok {
		return models.ErrCartIsEmpty
	}

	skuCount.Sub(float64(s.carts[userID].getItemsAmount()))

	delete(s.carts, userID)
	cartCount.Dec()

	return nil
}

func (s *Storage) GetItemsByUserID(_ context.Context, userID int64) ([]models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.carts[userID]; !ok {
		return nil, models.ErrCartIsEmpty
	}

	items := s.carts[userID].getItems()
	if len(items) == 0 {
		return nil, models.ErrCartIsEmpty
	}

	return items, nil
}
