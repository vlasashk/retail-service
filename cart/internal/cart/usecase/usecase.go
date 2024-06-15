package usecase

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"route256/cart/internal/cart/models"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

//go:generate minimock -i cartAdder -p usecase_test
type cartAdder interface {
	AddItem(ctx context.Context, userID, skuID int64, count uint16) error
}

//go:generate minimock -i itemRemover -p usecase_test
type itemRemover interface {
	DeleteItem(ctx context.Context, userID, skuID int64) error
}

//go:generate minimock -i cartRemover -p usecase_test
type cartRemover interface {
	DeleteItemsByUserID(ctx context.Context, userID int64) error
}

//go:generate minimock -i cartRetriever -p usecase_test
type cartRetriever interface {
	GetItemsByUserID(ctx context.Context, userID int64) ([]models.Item, error)
}

//go:generate minimock -i productProvider -p usecase_test
type productProvider interface {
	GetProduct(ctx context.Context, sku int64) (models.ItemDescription, error)
}

//go:generate minimock -i stocksProvider -p usecase_test
type stocksProvider interface {
	OrderCreate(ctx context.Context, order models.Order) (int64, error)
	StocksInfo(ctx context.Context, skuID int64) (uint64, error)
}

type UseCase struct {
	adder         cartAdder
	itemRemover   itemRemover
	cartRemover   cartRemover
	retriever     cartRetriever
	prodProvider  productProvider
	stockProvider stocksProvider
}

func New(
	adder cartAdder,
	itemRemover itemRemover,
	cartRemover cartRemover,
	retriever cartRetriever,
	prodProvider productProvider,
	stockProvider stocksProvider,
) *UseCase {
	return &UseCase{
		adder:         adder,
		itemRemover:   itemRemover,
		cartRemover:   cartRemover,
		retriever:     retriever,
		prodProvider:  prodProvider,
		stockProvider: stockProvider,
	}
}

func (uc *UseCase) AddItem(ctx context.Context, userID int64, skuID int64, count uint16) error {
	_, err := uc.prodProvider.GetProduct(ctx, skuID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return models.ErrItemProvider
		}
		return err
	}

	available, err := uc.stockProvider.StocksInfo(ctx, skuID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return models.ErrStockProvider
		}
		return err
	}

	if available < uint64(count) {
		return models.ErrInsufficientStock
	}

	return uc.adder.AddItem(ctx, userID, skuID, count)
}

func (uc *UseCase) DeleteItem(ctx context.Context, userID int64, skuID int64) error {
	return uc.itemRemover.DeleteItem(ctx, userID, skuID)
}

func (uc *UseCase) DeleteItemsByUserID(ctx context.Context, userID int64) error {
	return uc.cartRemover.DeleteItemsByUserID(ctx, userID)
}

func (uc *UseCase) CartCheckout(ctx context.Context, userID int64) (int64, error) {
	itemSKUs, err := uc.retriever.GetItemsByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	order := models.Order{
		UserID: userID,
		Items:  itemSKUs,
	}

	orderID, err := uc.stockProvider.OrderCreate(ctx, order)
	if err != nil {
		return 0, err
	}

	// Не будем прокидывать ошибку на весь хендлер, если не удалось удалить корзину (Только залогируем).
	// Так как заказ был успешно зарегистрирован и стоки забронированы для пользователя
	err = uc.cartRemover.DeleteItemsByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Caller().Send()
	}

	return orderID, nil
}

func (uc *UseCase) GetItemsByUserID(ctx context.Context, userID int64) (models.ItemsInCart, error) {
	itemSKUs, err := uc.retriever.GetItemsByUserID(ctx, userID)
	if err != nil {
		return models.ItemsInCart{}, err
	}

	cart, err := calcCart(ctx, uc.prodProvider, itemSKUs)
	if err != nil {
		return models.ItemsInCart{}, err
	}

	return cart, nil
}

func calcCart(ctx context.Context, provider productProvider, itemSKUs []models.Item) (models.ItemsInCart, error) {
	cart := models.ItemsInCart{
		Items:      make([]models.Item, 0, len(itemSKUs)),
		TotalPrice: 0,
	}
	itemChan := make(chan models.Item, len(itemSKUs))

	// горутина для сбора асинхронных ответов от Product Service для формирования корзины
	// (должна завершиться только после получения всех ответов от сервиса)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for item := range itemChan {
			cart.TotalPrice += uint32(item.Count) * item.Info.Price
			cart.Items = append(cart.Items, item)
		}
	}()

	// errgroup для асинхронного опроса Product Service
	eg, gCtx := errgroup.WithContext(ctx)
	for _, item := range itemSKUs {
		item := item
		eg.Go(func() error {
			ctx, cancel := context.WithTimeout(gCtx, time.Second*2)
			defer cancel()

			itemInfo, err := provider.GetProduct(ctx, item.SkuID)
			if err != nil {
				return err
			}

			item.Info = itemInfo
			itemChan <- item
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		close(itemChan)
		return models.ItemsInCart{}, err
	}

	close(itemChan)
	wg.Wait()

	sort.Slice(cart.Items, func(i, j int) bool {
		return cart.Items[i].SkuID < cart.Items[j].SkuID
	})

	return cart, nil
}
