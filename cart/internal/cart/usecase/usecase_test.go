package usecase_test

import (
	"context"
	"errors"
	"testing"

	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/usecase"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

type mocksToUse struct {
	Adder           *CartAdderMock
	ItemRemoverMock *ItemRemoverMock
	CartRemover     *CartRemoverMock
	ProductProvider *ProductProviderMock
	CartRetriever   *CartRetrieverMock
	StockProvider   *StocksProviderMock
}

func initMocks(t *testing.T) *mocksToUse {
	mc := minimock.NewController(t)
	return &mocksToUse{
		Adder:           NewCartAdderMock(mc),
		ItemRemoverMock: NewItemRemoverMock(mc),
		CartRemover:     NewCartRemoverMock(mc),
		ProductProvider: NewProductProviderMock(mc),
		CartRetriever:   NewCartRetrieverMock(mc),
		StockProvider:   NewStocksProviderMock(mc),
	}
}

func initUseCase(mocks *mocksToUse) *usecase.UseCase {
	return usecase.New(
		mocks.Adder,
		mocks.ItemRemoverMock,
		mocks.CartRemover,
		mocks.CartRetriever,
		mocks.ProductProvider,
		mocks.StockProvider,
	)
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestAddItem(t *testing.T) {
	t.Parallel()
	testItem := models.Item{
		SkuID: 1000,
		Count: 5,
		Info: models.ItemDescription{
			Name:  "TEST",
			Price: 1000,
		},
	}

	stockData := uint64(10)

	anyErr := errors.New("any error")

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse, int64)
		userID      int64
		expectedErr error
	}{
		{
			name: "AddItemSuccess",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItem.SkuID).Then(testItem.Info, nil)
				m.StockProvider.StocksInfoMock.When(minimock.AnyContext, testItem.SkuID).Then(stockData, nil)
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(nil)
			},
			userID:      999,
			expectedErr: nil,
		},
		{
			name: "AddItemProductDoesntExist",
			mockSetUp: func(m *mocksToUse, _ int64) {
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItem.SkuID).Then(models.ItemDescription{}, models.ErrNotFound)
			},
			userID:      42,
			expectedErr: models.ErrNotFound,
		},
		{
			name: "AddItemProductProviderErr",
			mockSetUp: func(m *mocksToUse, _ int64) {
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItem.SkuID).Then(models.ItemDescription{}, errors.New("any error"))
			},
			userID:      13,
			expectedErr: models.ErrItemProvider,
		},
		{
			name: "AddItemStockDoesntExist",
			mockSetUp: func(m *mocksToUse, _ int64) {
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItem.SkuID).Then(testItem.Info, nil)
				m.StockProvider.StocksInfoMock.When(minimock.AnyContext, testItem.SkuID).Then(0, models.ErrNotFound)
			},
			userID:      42,
			expectedErr: models.ErrNotFound,
		},
		{
			name: "AddItemStockProviderErr",
			mockSetUp: func(m *mocksToUse, _ int64) {
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItem.SkuID).Then(testItem.Info, nil)
				m.StockProvider.StocksInfoMock.When(minimock.AnyContext, testItem.SkuID).Then(0, errors.New("any error"))
			},
			userID:      13,
			expectedErr: models.ErrStockProvider,
		},
		{
			name: "AddItemInsufficientStockErr",
			mockSetUp: func(m *mocksToUse, _ int64) {
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItem.SkuID).Then(testItem.Info, nil)
				m.StockProvider.StocksInfoMock.When(minimock.AnyContext, testItem.SkuID).Then(0, nil)
			},
			userID:      13,
			expectedErr: models.ErrInsufficientStock,
		},
		{
			name: "AddItemAdderErr",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItem.SkuID).Then(testItem.Info, nil)
				m.StockProvider.StocksInfoMock.When(minimock.AnyContext, testItem.SkuID).Then(stockData, nil)
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(anyErr)
			},
			userID:      13,
			expectedErr: anyErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks, tt.userID)
			uc := initUseCase(mocks)

			err := uc.AddItem(context.Background(), tt.userID, testItem.SkuID, testItem.Count)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestGetItemsByUserID(t *testing.T) {
	t.Parallel()
	testItemAlpha := models.Item{
		SkuID: 1000,
		Count: 5,
	}
	testItemBeta := models.Item{
		SkuID: 5000,
		Count: 7,
	}

	itemAlphaDescr := models.ItemDescription{
		Name:  "ALPHA",
		Price: 1000,
	}
	itemBetaDescr := models.ItemDescription{
		Name:  "Beta",
		Price: 500,
	}

	anyErr := errors.New("any error")

	tests := []struct {
		name          string
		mockSetUp     func(*mocksToUse, int64)
		userID        int64
		expectedErr   error
		expectedItems models.ItemsInCart
	}{
		{
			name: "GetItemsSuccess_SingleItem",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then([]models.Item{testItemAlpha}, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemAlpha.SkuID).Then(itemAlphaDescr, nil)
			},
			userID:      999,
			expectedErr: nil,
			expectedItems: models.ItemsInCart{
				Items: []models.Item{
					{
						SkuID: testItemAlpha.SkuID,
						Count: testItemAlpha.Count,
						Info:  itemAlphaDescr,
					},
				},
				TotalPrice: uint32(testItemAlpha.Count) * itemAlphaDescr.Price,
			},
		},
		{
			name: "GetItemsSuccess_DoubleItems_Sorted",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then([]models.Item{testItemAlpha, testItemBeta}, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemAlpha.SkuID).Then(itemAlphaDescr, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemBeta.SkuID).Then(itemBetaDescr, nil)
			},
			userID:      999,
			expectedErr: nil,
			expectedItems: models.ItemsInCart{
				Items: []models.Item{
					{
						SkuID: testItemAlpha.SkuID,
						Count: testItemAlpha.Count,
						Info:  itemAlphaDescr,
					},
					{
						SkuID: testItemBeta.SkuID,
						Count: testItemBeta.Count,
						Info:  itemBetaDescr,
					},
				},
				TotalPrice: uint32(testItemAlpha.Count)*itemAlphaDescr.Price + uint32(testItemBeta.Count)*itemBetaDescr.Price,
			},
		},
		{
			name: "GetItemsSuccess_DoubleItems_UnSorted",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then([]models.Item{testItemBeta, testItemAlpha}, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemAlpha.SkuID).Then(itemAlphaDescr, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemBeta.SkuID).Then(itemBetaDescr, nil)
			},
			userID:      999,
			expectedErr: nil,
			expectedItems: models.ItemsInCart{
				Items: []models.Item{
					{
						SkuID: testItemAlpha.SkuID,
						Count: testItemAlpha.Count,
						Info:  itemAlphaDescr,
					},
					{
						SkuID: testItemBeta.SkuID,
						Count: testItemBeta.Count,
						Info:  itemBetaDescr,
					},
				},
				TotalPrice: uint32(testItemAlpha.Count)*itemAlphaDescr.Price + uint32(testItemBeta.Count)*itemBetaDescr.Price,
			},
		},
		{
			name: "GetItemsErr",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then(nil, anyErr)
			},
			userID:      42,
			expectedErr: anyErr,
		},
		{
			name: "GetItemsProviderErr_ForAlpha",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then([]models.Item{testItemBeta, testItemAlpha}, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemBeta.SkuID).Then(itemBetaDescr, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemAlpha.SkuID).Then(models.ItemDescription{}, anyErr)
			},
			userID:      13,
			expectedErr: anyErr,
		},
		{
			name: "GetItemsProviderErr_ForBeta",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then([]models.Item{testItemBeta, testItemAlpha}, nil)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemBeta.SkuID).Then(models.ItemDescription{}, anyErr)
				m.ProductProvider.GetProductMock.When(minimock.AnyContext, testItemAlpha.SkuID).Then(itemAlphaDescr, nil)
			},
			userID:      13,
			expectedErr: anyErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks, tt.userID)
			uc := initUseCase(mocks)

			items, err := uc.GetItemsByUserID(context.Background(), tt.userID)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, items, tt.expectedItems)
		})
	}
}

func TestCartCheckout(t *testing.T) {
	t.Parallel()
	defaultUserID := int64(999)

	testItemAlpha := models.Item{
		SkuID: 1000,
		Count: 5,
	}
	testItemBeta := models.Item{
		SkuID: 5000,
		Count: 7,
	}

	defaultOrder := models.Order{
		UserID: defaultUserID,
		Items:  []models.Item{testItemAlpha, testItemBeta},
	}

	anyErr := errors.New("any error")

	tests := []struct {
		name            string
		mockSetUp       func(*mocksToUse, int64)
		userID          int64
		expectedErr     error
		expectedOrderID int64
	}{
		{
			name: "CartCheckoutSuccess",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then([]models.Item{testItemAlpha, testItemBeta}, nil)
				m.StockProvider.OrderCreateMock.When(minimock.AnyContext, defaultOrder).Then(1000, nil)
				m.CartRemover.DeleteItemsByUserIDMock.When(minimock.AnyContext, userID).Then(nil)
			},
			userID:          defaultUserID,
			expectedErr:     nil,
			expectedOrderID: 1000,
		},
		{
			name: "CartCheckoutGetItemsErr",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then(nil, anyErr)
			},
			userID:      42,
			expectedErr: anyErr,
		},
		{
			name: "CartCheckoutDeleteItemsErr",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRetriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then([]models.Item{testItemAlpha, testItemBeta}, nil)
				m.StockProvider.OrderCreateMock.When(minimock.AnyContext, defaultOrder).Then(1000, nil)
				m.CartRemover.DeleteItemsByUserIDMock.When(minimock.AnyContext, userID).Then(errors.New("any error"))
			},
			userID:          defaultUserID,
			expectedErr:     nil,
			expectedOrderID: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks, tt.userID)
			uc := initUseCase(mocks)

			orderID, err := uc.CartCheckout(context.Background(), tt.userID)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, orderID, tt.expectedOrderID)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()
	testItem := models.Item{
		SkuID: 1000,
		Count: 5,
		Info: models.ItemDescription{
			Name:  "TEST",
			Price: 1000,
		},
	}

	anyErr := errors.New("any error")

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse, int64)
		userID      int64
		expectedErr error
	}{
		{
			name: "DeleteItemSuccess",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.ItemRemoverMock.DeleteItemMock.When(minimock.AnyContext, userID, testItem.SkuID).Then(nil)
			},
			userID:      999,
			expectedErr: nil,
		},
		{
			name: "DeleteItemErr",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.ItemRemoverMock.DeleteItemMock.When(minimock.AnyContext, userID, testItem.SkuID).Then(anyErr)
			},
			userID:      13,
			expectedErr: anyErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks, tt.userID)
			uc := initUseCase(mocks)

			err := uc.DeleteItem(context.Background(), tt.userID, testItem.SkuID)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestDeleteItemsByUserID(t *testing.T) {
	t.Parallel()
	anyErr := errors.New("any error")

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse, int64)
		userID      int64
		expectedErr error
	}{
		{
			name: "DeleteItemSuccess",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRemover.DeleteItemsByUserIDMock.When(minimock.AnyContext, userID).Then(nil)
			},
			userID:      999,
			expectedErr: nil,
		},
		{
			name: "DeleteItemErr",
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.CartRemover.DeleteItemsByUserIDMock.When(minimock.AnyContext, userID).Then(anyErr)
			},
			userID:      13,
			expectedErr: anyErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks, tt.userID)
			uc := initUseCase(mocks)

			err := uc.DeleteItemsByUserID(context.Background(), tt.userID)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
