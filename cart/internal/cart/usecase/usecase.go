package usecase

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"route256/cart/internal/cart/models"

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

type UseCase struct {
	adder       cartAdder
	itemRemover itemRemover
	cartRemover cartRemover
	retriever   cartRetriever
	provider    productProvider
}

func New(
	adder cartAdder,
	itemRemover itemRemover,
	cartRemover cartRemover,
	retriever cartRetriever,
	provider productProvider,
) *UseCase {
	return &UseCase{
		adder:       adder,
		itemRemover: itemRemover,
		cartRemover: cartRemover,
		retriever:   retriever,
		provider:    provider,
	}
}

func (uc *UseCase) AddItem(ctx context.Context, userID int64, skuID int64, count uint16) error {
	_, err := uc.provider.GetProduct(ctx, skuID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return models.ErrItemProvider
		}
		return err
	}

	return uc.adder.AddItem(ctx, userID, skuID, count)
}

func (uc *UseCase) DeleteItem(ctx context.Context, userID int64, skuID int64) error {
	return uc.itemRemover.DeleteItem(ctx, userID, skuID)
}

func (uc *UseCase) DeleteItemsByUserID(ctx context.Context, userID int64) error {
	return uc.cartRemover.DeleteItemsByUserID(ctx, userID)
}

func (uc *UseCase) GetItemsByUserID(ctx context.Context, userID int64) (models.ItemsInCart, error) {
	itemSKUs, err := uc.retriever.GetItemsByUserID(ctx, userID)
	if err != nil {
		return models.ItemsInCart{}, err
	}

	cart, err := calcCart(ctx, uc.provider, itemSKUs)
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
