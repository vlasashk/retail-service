package inmem

import (
	"sync"

	"route256/cart/internal/cart/models"
)

type cart struct {
	mu    sync.RWMutex
	items map[int64]uint16
}

func newCart() *cart {
	return &cart{
		mu:    sync.RWMutex{},
		items: make(map[int64]uint16),
	}
}

func (c *cart) addItem(skuID int64, count uint16) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[skuID]; ok {
		count += c.items[skuID]
	} else {
		skuCount.Inc()
	}

	c.items[skuID] = count
}

func (c *cart) deleteItem(skuID int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[skuID]; !ok {
		return false
	}

	delete(c.items, skuID)

	return true
}

func (c *cart) getItems() []models.Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res := make([]models.Item, 0, len(c.items))

	for skuID, count := range c.items {
		res = append(res, models.Item{SkuID: skuID, Count: count})
	}

	return res
}

func (c *cart) getItemsAmount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
