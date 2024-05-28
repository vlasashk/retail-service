package getcart

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/common"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/pkg/constants"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

//go:generate mockery --name=CartRetriever
type CartRetriever interface {
	GetItemsByUserID(ctx context.Context, userID int64) ([]models.Item, error)
}

type Handler struct {
	retriever       CartRetriever
	productProvider common.ProductProvider
	log             zerolog.Logger
}

func New(log zerolog.Logger, retriever CartRetriever, provider common.ProductProvider) *Handler {
	return &Handler{
		retriever:       retriever,
		productProvider: provider,
		log:             log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localLog := h.log.With().Str("handler", "get_cart").Logger()

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		localLog.Error().Err(err).Int64("user_id", userID).Send()
		errhandle.NewErr(constants.ErrInvalidUserID).Send(w, localLog, http.StatusBadRequest)
		return
	}

	itemSKUs, err := h.retriever.GetItemsByUserID(r.Context(), userID)
	if err != nil {
		localLog.Error().Err(err).Send()
		if errors.Is(err, models.ErrCartIsEmpty) {
			errhandle.NewErr(constants.ErrEmptyCart).Send(w, localLog, http.StatusNotFound)
			return
		}
		errhandle.NewErr(constants.ErrGetItems).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	cart, err := calcCart(r.Context(), h.productProvider, itemSKUs)
	if err != nil {
		localLog.Error().Err(err).Send()
		if errors.Is(err, models.ErrCartIsEmpty) {
			errhandle.NewErr(constants.ErrEmptyCart).Send(w, localLog, http.StatusNotFound)
			return
		}
		errhandle.NewErr(constants.ErrCartCheckout).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	cartResp := cartToDTO(cart)

	data, err := json.Marshal(cartResp)
	if err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(constants.ErrMarshal).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		localLog.Error().Err(err).Send()
	}
}

func calcCart(ctx context.Context, provider common.ProductProvider, itemSKUs []models.Item) (models.Cart, error) {
	cart := models.Cart{
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
			itemInfo, err := provider.GetProduct(ctx, item.SkuId)
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
		return models.Cart{}, err
	}

	close(itemChan)
	wg.Wait()

	sort.Slice(cart.Items, func(i, j int) bool {
		return cart.Items[i].SkuId < cart.Items[j].SkuId
	})

	return cart, nil
}
