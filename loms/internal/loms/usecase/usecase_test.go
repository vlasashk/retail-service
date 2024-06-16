package usecase_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"route256/loms/internal/loms/models"
	"route256/loms/internal/loms/usecase"

	"github.com/gojuno/minimock/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type mocksToUse struct {
	stockProvider *StockProviderMock
	orderManager  *OrderManagerMock
}

func initMocks(t *testing.T) *mocksToUse {
	mc := minimock.NewController(t)
	return &mocksToUse{
		stockProvider: NewStockProviderMock(mc),
		orderManager:  NewOrderManagerMock(mc),
	}
}

func initUseCase(mocks *mocksToUse) *usecase.UseCase {
	return usecase.New(
		zerolog.New(os.Stderr),
		mocks.orderManager,
		mocks.stockProvider,
	)
}

func TestOrderCreate(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultUserID := int64(1111)
	defaultOrderID := int64(2222)

	defaultItems := []models.Item{
		{
			SKUid: 12314,
			Count: 5,
		},
		{
			SKUid: 6713,
			Count: 7,
		},
	}
	defaultOrder := models.Order{
		UserID: defaultUserID,
		Items:  defaultItems,
		Status: models.NewStatus,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		order       models.Order
		expectedErr error
		expectedRes int64
	}{
		{
			name: "OrderCreateSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.CreateMock.When(minimock.AnyContext, defaultOrder).Then(defaultOrderID, nil)
				m.stockProvider.ReserveMock.When(minimock.AnyContext, defaultOrder).Then(nil)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.AwaitingPaymentStatus).Then(nil)
			},
			expectedErr: nil,
			order:       defaultOrder,
			expectedRes: defaultOrderID,
		},
		{
			name: "OrderCreateSetAwaitingFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.CreateMock.When(minimock.AnyContext, defaultOrder).Then(defaultOrderID, nil)
				m.stockProvider.ReserveMock.When(minimock.AnyContext, defaultOrder).Then(nil)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.AwaitingPaymentStatus).Then(anyErr)
			},
			expectedErr: nil,
			order:       defaultOrder,
			expectedRes: defaultOrderID,
		},
		{
			name: "OrderCreateCreateFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.CreateMock.When(minimock.AnyContext, defaultOrder).Then(0, anyErr)
			},
			expectedErr: anyErr,
			order:       defaultOrder,
			expectedRes: 0,
		},
		{
			name: "OrderCreateCreateFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.CreateMock.When(minimock.AnyContext, defaultOrder).Then(0, anyErr)
			},
			expectedErr: anyErr,
			order:       defaultOrder,
			expectedRes: 0,
		},
		{
			name: "OrderCreateSetFailedStatusFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.CreateMock.When(minimock.AnyContext, defaultOrder).Then(defaultOrderID, nil)
				m.stockProvider.ReserveMock.When(minimock.AnyContext, defaultOrder).Then(anyErr)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.FailedStatus).Then(anyErr)
			},
			expectedErr: anyErr,
			order:       defaultOrder,
			expectedRes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			orderID, err := uc.OrderCreate(context.Background(), tt.order)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedRes, orderID)
		})
	}
}

func TestOrderPay(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultUserID := int64(1111)
	defaultOrderID := int64(2222)

	defaultItems := []models.Item{
		{
			SKUid: 12314,
			Count: 5,
		},
		{
			SKUid: 6713,
			Count: 7,
		},
	}
	defaultOrder := models.Order{
		UserID: defaultUserID,
		Items:  defaultItems,
		Status: models.AwaitingPaymentStatus,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		order       models.Order
		expectedErr error
		expectedRes int64
	}{
		{
			name: "OrderPaySuccess",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
				m.stockProvider.ReserveRemoveMock.When(minimock.AnyContext, defaultOrder).Then(nil)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.PayedStatus).Then(nil)
			},
			expectedErr: nil,
			order:       defaultOrder,
			expectedRes: defaultOrderID,
		},
		{
			name: "OrderPaySetAwaitingFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
				m.stockProvider.ReserveRemoveMock.When(minimock.AnyContext, defaultOrder).Then(nil)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.PayedStatus).Then(anyErr)
			},
			expectedErr: nil,
			order:       defaultOrder,
			expectedRes: defaultOrderID,
		},
		{
			name: "OrderPayGetByOrderIDMockFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(models.Order{}, anyErr)
			},
			expectedErr: anyErr,
			order:       defaultOrder,
			expectedRes: 0,
		},
		{
			name: "OrderPayOrderWithWrongStatusFail",
			mockSetUp: func(m *mocksToUse) {
				order := defaultOrder
				order.Status = models.PayedStatus
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(order, nil)
			},
			expectedErr: models.ErrOrderStatusConflict,
			order:       defaultOrder,
			expectedRes: 0,
		},
		{
			name: "OrderPaySetFailedStatusFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
				m.stockProvider.ReserveRemoveMock.When(minimock.AnyContext, defaultOrder).Then(anyErr)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.FailedStatus).Then(anyErr)
			},
			expectedErr: anyErr,
			order:       defaultOrder,
			expectedRes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			err := uc.OrderPay(context.Background(), defaultOrderID)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestOrderCancel(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultUserID := int64(1111)
	defaultOrderID := int64(2222)

	defaultItems := []models.Item{
		{
			SKUid: 12314,
			Count: 5,
		},
		{
			SKUid: 6713,
			Count: 7,
		},
	}
	defaultOrder := models.Order{
		UserID: defaultUserID,
		Items:  defaultItems,
		Status: models.AwaitingPaymentStatus,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		order       models.Order
		expectedErr error
		expectedRes int64
	}{
		{
			name: "OrderCancelSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
				m.stockProvider.ReserveCancelMock.When(minimock.AnyContext, defaultOrder).Then(nil)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.CancelledStatus).Then(nil)
			},
			expectedErr: nil,
			order:       defaultOrder,
			expectedRes: defaultOrderID,
		},
		{
			name: "OrderCancelSetAwaitingFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
				m.stockProvider.ReserveCancelMock.When(minimock.AnyContext, defaultOrder).Then(nil)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.CancelledStatus).Then(anyErr)
			},
			expectedErr: nil,
			order:       defaultOrder,
			expectedRes: defaultOrderID,
		},
		{
			name: "OrderCancelGetByOrderIDMockFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(models.Order{}, anyErr)
			},
			expectedErr: anyErr,
			order:       defaultOrder,
			expectedRes: 0,
		},
		{
			name: "OrderCancelOrderWithWrongStatusFail",
			mockSetUp: func(m *mocksToUse) {
				order := defaultOrder
				order.Status = models.CancelledStatus
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(order, nil)
			},
			expectedErr: models.ErrOrderStatusConflict,
			order:       defaultOrder,
			expectedRes: 0,
		},
		{
			name: "OrderCancelSetFailedStatusFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
				m.stockProvider.ReserveCancelMock.When(minimock.AnyContext, defaultOrder).Then(anyErr)
				m.orderManager.SetStatusMock.When(minimock.AnyContext, defaultOrderID, models.FailedStatus).Then(anyErr)
			},
			expectedErr: anyErr,
			order:       defaultOrder,
			expectedRes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			err := uc.OrderCancel(context.Background(), defaultOrderID)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestOrderInfo(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultUserID := int64(1111)
	defaultOrderID := int64(2222)

	defaultItems := []models.Item{
		{
			SKUid: 12314,
			Count: 5,
		},
		{
			SKUid: 6713,
			Count: 7,
		},
	}
	defaultOrder := models.Order{
		UserID: defaultUserID,
		Items:  defaultItems,
		Status: models.AwaitingPaymentStatus,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		order       models.Order
		expectedErr error
	}{
		{
			name: "OrderInfoSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
			},
			expectedErr: nil,
			order:       defaultOrder,
		},
		{
			name: "OrderInfoFail",
			mockSetUp: func(m *mocksToUse) {
				m.orderManager.GetByOrderIDMock.When(minimock.AnyContext, defaultOrderID).Then(models.Order{}, anyErr)
			},
			expectedErr: anyErr,
			order:       models.Order{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			order, err := uc.OrderInfo(context.Background(), defaultOrderID)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.order, order)
		})
	}
}

func TestStockInfo(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultSKUID := uint32(1111)

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		expectedErr error
		expectedRes int64
	}{
		{
			name: "OrderInfoSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.stockProvider.GetBySKUMock.When(minimock.AnyContext, defaultSKUID).Then(123, nil)
			},
			expectedErr: nil,
			expectedRes: 123,
		},
		{
			name: "OrderInfoFail",
			mockSetUp: func(m *mocksToUse) {
				m.stockProvider.GetBySKUMock.When(minimock.AnyContext, defaultSKUID).Then(0, anyErr)
			},
			expectedErr: anyErr,
			expectedRes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			order, err := uc.StocksInfo(context.Background(), defaultSKUID)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedRes, order)
		})
	}
}
