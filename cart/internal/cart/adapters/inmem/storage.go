package inmem

import (
	"context"
	"sync"

	"route256/cart/internal/cart/models"
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

	s.carts[userID].deleteItem(skuID)

	return nil
}

func (s *Storage) DeleteItemsByUserID(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.carts[userID]; !ok {
		return models.ErrCartIsEmpty
	}

	delete(s.carts, userID)

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
