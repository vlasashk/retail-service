package impl_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"route256/loms/internal/loms/models"
	"route256/loms/internal/loms/ports/ggrpc/impl"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/gojuno/minimock/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mocksToUse struct {
	lomsService *LomsServerMock
}

func initMocks(t *testing.T) *mocksToUse {
	mc := minimock.NewController(t)
	return &mocksToUse{
		lomsService: NewLomsServerMock(mc),
	}
}

func initUseCase(mocks *mocksToUse) *impl.Impl {
	return impl.New(
		zerolog.New(os.Stderr),
		mocks.lomsService,
	)
}

func TestOrderCreate(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultUserID := int64(1111)
	defaultOrderID := int64(2222)

	defaultReq := &lomsservicev1.OrderCreateRequest{
		User: defaultUserID,
		Items: []*lomsservicev1.Item{
			{
				Sku:   12314,
				Count: 5,
			},
			{
				Sku:   6713,
				Count: 7,
			},
		},
	}

	defaultOrder := models.Order{
		UserID: defaultUserID,
		Items: []models.Item{
			{
				SKUid: 12314,
				Count: 5,
			},
			{
				SKUid: 6713,
				Count: 7,
			},
		},
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		req         *lomsservicev1.OrderCreateRequest
		expectedErr error
		expectedRes *lomsservicev1.OrderCreateResponse
	}{
		{
			name: "OrderCreateSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderCreateMock.When(minimock.AnyContext, defaultOrder).Then(defaultOrderID, nil)
			},
			expectedErr: nil,
			req:         defaultReq,
			expectedRes: &lomsservicev1.OrderCreateResponse{
				OrderID: defaultOrderID,
			},
		},
		{
			name: "OrderCreateFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderCreateMock.When(minimock.AnyContext, defaultOrder).Then(0, anyErr)
			},
			expectedErr: status.Error(codes.Internal, "internal server error"),
			req:         defaultReq,
			expectedRes: nil,
		},
		{
			name: "OrderCreateNotFoundFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderCreateMock.When(minimock.AnyContext, defaultOrder).Then(0, models.ErrItemNotFound)
			},
			expectedErr: status.Error(codes.NotFound, models.ErrItemNotFound.Error()),
			req:         defaultReq,
			expectedRes: nil,
		},
		{
			name: "OrderCreateInsufficientStockFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderCreateMock.When(minimock.AnyContext, defaultOrder).Then(0, models.ErrInsufficientStock)
			},
			expectedErr: status.Error(codes.FailedPrecondition, models.ErrInsufficientStock.Error()),
			req:         defaultReq,
			expectedRes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			orderResp, err := uc.OrderCreate(context.Background(), tt.req)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedRes, orderResp)
		})
	}
}

func TestOrderInfo(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultUserID := int64(1111)
	defaultOrderID := int64(2222)

	defaultReq := &lomsservicev1.OrderInfoRequest{
		OrderID: defaultOrderID,
	}

	defaultOrder := models.Order{
		UserID: defaultUserID,
		Items: []models.Item{
			{
				SKUid: 12314,
				Count: 5,
			},
			{
				SKUid: 6713,
				Count: 7,
			},
		},
		Status: models.PayedStatus,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		req         *lomsservicev1.OrderInfoRequest
		expectedErr error
		expectedRes *lomsservicev1.OrderInfoResponse
	}{
		{
			name: "OrderInfoSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderInfoMock.When(minimock.AnyContext, defaultOrderID).Then(defaultOrder, nil)
			},
			expectedErr: nil,
			req:         defaultReq,
			expectedRes: &lomsservicev1.OrderInfoResponse{
				User: defaultUserID,
				Items: []*lomsservicev1.Item{
					{
						Sku:   12314,
						Count: 5,
					},
					{
						Sku:   6713,
						Count: 7,
					},
				},
				Status: lomsservicev1.Status_STATUS_PAYED,
			},
		},
		{
			name: "OrderInfoFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderInfoMock.When(minimock.AnyContext, defaultOrderID).Then(models.Order{}, anyErr)
			},
			expectedErr: status.Error(codes.Internal, "internal server error"),
			req:         defaultReq,
			expectedRes: nil,
		},
		{
			name: "OrderInfoNotFoundFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderInfoMock.When(minimock.AnyContext, defaultOrderID).Then(models.Order{}, models.ErrOrderNotFound)
			},
			expectedErr: status.Error(codes.NotFound, models.ErrOrderNotFound.Error()),
			req:         defaultReq,
			expectedRes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			orderResp, err := uc.OrderInfo(context.Background(), tt.req)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedRes, orderResp)
		})
	}
}

func TestOrderPay(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultOrderID := int64(2222)

	defaultReq := &lomsservicev1.OrderPayRequest{
		OrderID: defaultOrderID,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		req         *lomsservicev1.OrderPayRequest
		expectedErr error
	}{
		{
			name: "OrderPaySuccess",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderPayMock.When(minimock.AnyContext, defaultOrderID).Then(nil)
			},
			expectedErr: nil,
			req:         defaultReq,
		},
		{
			name: "OrderPayFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderPayMock.When(minimock.AnyContext, defaultOrderID).Then(anyErr)
			},
			expectedErr: status.Error(codes.Internal, "internal server error"),
			req:         defaultReq,
		},
		{
			name: "OrderPayNotFoundFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderPayMock.When(minimock.AnyContext, defaultOrderID).Then(models.ErrOrderNotFound)
			},
			expectedErr: status.Error(codes.NotFound, models.ErrOrderNotFound.Error()),
			req:         defaultReq,
		},
		{
			name: "OrderPayReservationConflictFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderPayMock.When(minimock.AnyContext, defaultOrderID).Then(models.ErrReservationConflict)
			},
			expectedErr: status.Error(codes.FailedPrecondition, models.ErrReservationConflict.Error()),
			req:         defaultReq,
		},
		{
			name: "OrderPayPaymentStatusConflictFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderPayMock.When(minimock.AnyContext, defaultOrderID).Then(models.ErrOrderStatusConflict)
			},
			expectedErr: status.Error(codes.FailedPrecondition, models.ErrOrderStatusConflict.Error()),
			req:         defaultReq,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			_, err := uc.OrderPay(context.Background(), tt.req)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestOrderCancel(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultOrderID := int64(2222)

	defaultReq := &lomsservicev1.OrderCancelRequest{
		OrderID: defaultOrderID,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		req         *lomsservicev1.OrderCancelRequest
		expectedErr error
	}{
		{
			name: "OrderCancelSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderCancelMock.When(minimock.AnyContext, defaultOrderID).Then(nil)
			},
			expectedErr: nil,
			req:         defaultReq,
		},
		{
			name: "OrderCancelFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderCancelMock.When(minimock.AnyContext, defaultOrderID).Then(anyErr)
			},
			expectedErr: status.Error(codes.Internal, "internal server error"),
			req:         defaultReq,
		},
		{
			name: "OrderCancelReservationConflictFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.OrderCancelMock.When(minimock.AnyContext, defaultOrderID).Then(models.ErrReservationConflict)
			},
			expectedErr: status.Error(codes.FailedPrecondition, models.ErrReservationConflict.Error()),
			req:         defaultReq,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			_, err := uc.OrderCancel(context.Background(), tt.req)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestStocksInfo(t *testing.T) {
	t.Parallel()

	anyErr := errors.New("any error")

	defaultSkuID := uint32(2222)
	defaultCount := int64(5)

	defaultReq := &lomsservicev1.StocksInfoRequest{
		Sku: defaultSkuID,
	}

	tests := []struct {
		name        string
		mockSetUp   func(*mocksToUse)
		req         *lomsservicev1.StocksInfoRequest
		expectedErr error
		expectedRes *lomsservicev1.StocksInfoResponse
	}{
		{
			name: "StocksInfoSuccess",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.StocksInfoMock.When(minimock.AnyContext, defaultSkuID).Then(defaultCount, nil)
			},
			expectedErr: nil,
			req:         defaultReq,
			expectedRes: &lomsservicev1.StocksInfoResponse{
				Count: uint64(defaultCount),
			},
		},
		{
			name: "StocksInfoFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.StocksInfoMock.When(minimock.AnyContext, defaultSkuID).Then(0, anyErr)
			},
			expectedErr: status.Error(codes.Internal, "internal server error"),
			req:         defaultReq,
			expectedRes: nil,
		},
		{
			name: "StocksInfoReservationConflictFail",
			mockSetUp: func(m *mocksToUse) {
				m.lomsService.StocksInfoMock.When(minimock.AnyContext, defaultSkuID).Then(0, models.ErrItemNotFound)
			},
			expectedErr: status.Error(codes.NotFound, models.ErrItemNotFound.Error()),
			req:         defaultReq,
			expectedRes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := initMocks(t)
			tt.mockSetUp(mocks)
			uc := initUseCase(mocks)

			count, err := uc.StocksInfo(context.Background(), tt.req)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedRes, count)
		})
	}
}
