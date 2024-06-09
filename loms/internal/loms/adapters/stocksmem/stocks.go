package stocksmem

import (
	"context"
	"sync"

	"route256/loms/internal/loms/models"
)

type Stocks struct {
	mu       *sync.RWMutex
	stocks   map[uint32]int64 // Общее количество доступных для заказа товаров (skuID -> amount). Количество уменьшается только после оплаты
	reserved map[uint32]int64 // Зарезервированные товары (не оплаченные) (skuID -> amount)
}

func New() *Stocks {
	return &Stocks{
		mu:       &sync.RWMutex{},
		stocks:   make(map[uint32]int64),
		reserved: make(map[uint32]int64),
	}
}

func (s *Stocks) Reserve(_ context.Context, order models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	changesSnapshot := make(map[uint32]int64)

	// для восстановления данных в исходное состояние в случае возникновения ошибки
	defer func() {
		if err != nil {
			for sku, amount := range changesSnapshot {
				s.reserved[sku] -= amount
			}
		}
	}()

	for _, item := range order.Items {
		available, ok := s.stocks[item.SKUid]
		if !ok {
			err = models.ErrItemNotFound
			return err
		}

		toReserve := int64(item.Count)
		if available < s.reserved[item.SKUid]+toReserve {
			err = models.ErrInsufficientStock
			return err
		}

		changesSnapshot[item.SKUid] += toReserve
		s.reserved[item.SKUid] += toReserve
	}

	return nil
}

func (s *Stocks) ReserveRemove(_ context.Context, order models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	changesSnapshot := make(map[uint32]int64)

	// для восстановления данных в исходное состояние в случае возникновения ошибки
	defer func() {
		if err != nil {
			for sku, amount := range changesSnapshot {
				s.reserved[sku] += amount
				s.stocks[sku] += amount
			}
		}
	}()

	for _, item := range order.Items {
		// не проверяем s.stocks так как согласованность s.reserved c s.stocks проверяется в Reserve()
		reserved, ok := s.reserved[item.SKUid]
		if !ok {
			err = models.ErrItemNotFound
			return err
		}

		toRemove := int64(item.Count)
		if reserved < toRemove {
			err = models.ErrReservationConflict
			return err
		}

		changesSnapshot[item.SKUid] += toRemove
		s.reserved[item.SKUid] -= toRemove
		s.stocks[item.SKUid] -= toRemove
	}

	return nil
}

func (s *Stocks) ReserveCancel(_ context.Context, order models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	changesSnapshot := make(map[uint32]int64)

	// для восстановления данных в исходное состояние в случае возникновения ошибки
	defer func() {
		if err != nil {
			for sku, amount := range changesSnapshot {
				s.reserved[sku] += amount
			}
		}
	}()

	for _, item := range order.Items {
		reserved, ok := s.reserved[item.SKUid]
		if !ok {
			err = models.ErrItemNotFound
			return err
		}

		toRemove := int64(item.Count)
		if reserved < toRemove {
			err = models.ErrReservationConflict
			return err
		}

		changesSnapshot[item.SKUid] += toRemove
		s.reserved[item.SKUid] -= toRemove
	}

	return nil
}

func (s *Stocks) GetBySKU(_ context.Context, skuID uint32) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.stocks[skuID]; !ok {
		return 0, models.ErrItemNotFound
	}
	// нет проверки на наличие в s.reserved так как нерелевантна (в любом случае получим zero value)
	total := s.stocks[skuID] - s.reserved[skuID]

	return total, nil
}

// UploadStockData - добавляет данные из stock-data.json. Перезаписывает данные в случае их наличия
func (s *Stocks) UploadStockData(_ context.Context, stocks []models.Stock) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, stock := range stocks {
		s.stocks[stock.SKUid] += int64(stock.Count)
		s.reserved[stock.SKUid] = int64(stock.Reserved)
	}

	return nil
}
