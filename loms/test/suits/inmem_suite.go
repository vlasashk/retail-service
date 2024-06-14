//go:build integration

package suits

import (
	"context"
	"time"

	"route256/loms/config"
	"route256/loms/internal/loms"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	defaultUserID     = 111111
	defaultAlphaSKUID = 1076963 // В stock-data.json 100 всего и 10 зарезервировано
	defaultBetaSKUID  = 1148162 // В stock-data.json 150 всего и 20 зарезервировано
	additionalSKUID   = 1625903 // В stock-data.json 200 всего и 30 зарезервировано
	paySKUID          = 2618151 // В stock-data.json 50 всего и 5 зарезервировано
)

var defaultReq = &lomsservicev1.OrderCreateRequest{
	User: defaultUserID,
	Items: []*lomsservicev1.Item{
		{
			Sku:   defaultAlphaSKUID,
			Count: 10,
		},
		{
			Sku:   defaultBetaSKUID,
			Count: 30,
		},
	},
}

var additionalReq = &lomsservicev1.OrderCreateRequest{
	User: defaultUserID,
	Items: []*lomsservicev1.Item{
		{
			Sku:   defaultAlphaSKUID,
			Count: 30,
		},
		{
			Sku:   additionalSKUID,
			Count: 70,
		},
	},
}

type InmemSuite struct {
	suite.Suite
	client   lomsservicev1.LOMSClient
	conn     *grpc.ClientConn
	cancel   context.CancelFunc
	orderIDs []int64
}

func (s *InmemSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	cfg, err := config.New()
	if err != nil {
		s.T().Fatal(err)
	}

	go func() {
		if err = loms.Run(ctx, cfg); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	s.conn, err = grpc.NewClient("localhost:"+cfg.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.T().Fatal(err)
	}

	s.client = lomsservicev1.NewLOMSClient(s.conn)

	if err = s.healthCheck(10); err != nil {
		s.T().Fatal(err)
	}

}

func (s *InmemSuite) TearDownSuite() {
	s.cancel()
	if s.conn != nil {
		_ = s.conn.Close()
	}
}

func (s *InmemSuite) healthCheck(attempts int) error {
	var err error
	healthClient := grpc_health_v1.NewHealthClient(s.conn)

	for attempts > 0 {
		req := &grpc_health_v1.HealthCheckRequest{Service: ""}

		resp, err := healthClient.Check(context.Background(), req)
		if err == nil && resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			log.Info().Msg("Service is ready to accept requests")
			return nil
		}

		log.Debug().Err(err).Int("attempts left", attempts).Msg("Service is not available for integration rests")
		time.Sleep(time.Second)
		attempts--
	}

	return err
}

// setupOrders вставляются 2 заказа на 3 SKU и база принимает следующий вид
// 1076963 - В stock-data.json 100 всего и 50 зарезервировано - 50 доступно
// 1148162 - В stock-data.json 150 всего и 50 зарезервировано - 100 доступно
// 1625903 - В stock-data.json 200 всего и 100 зарезервировано - 100 доступно
func (s *InmemSuite) setupOrders() {
	firstOrderID, err := s.client.OrderCreate(context.Background(), additionalReq)
	s.Require().NoError(err)
	secondOrderID, err := s.client.OrderCreate(context.Background(), defaultReq)
	s.Require().NoError(err)

	s.orderIDs = append(s.orderIDs, firstOrderID.OrderID, secondOrderID.OrderID)
}

func (s *InmemSuite) removeOrders(skuIDs ...int64) {
	for _, skuID := range skuIDs {
		_, err := s.client.OrderCancel(context.Background(), &lomsservicev1.OrderCancelRequest{
			OrderID: skuID,
		})
		s.Require().NoError(err)
	}
}

func (s *InmemSuite) TestOrderInfo() {
	tests := []struct {
		name       string
		req        func() *lomsservicev1.OrderInfoRequest
		setup      func()
		cleanup    func(...int64)
		expectErr  error
		expectResp *lomsservicev1.OrderInfoResponse
	}{
		{
			name:    "OrderInfoSuccess",
			setup:   s.setupOrders,
			cleanup: s.removeOrders,
			req: func() *lomsservicev1.OrderInfoRequest {
				// Так как используем сетапер то ожидаем что создалось 2 заказа
				s.Require().Len(s.orderIDs, 2)

				return &lomsservicev1.OrderInfoRequest{
					OrderID: s.orderIDs[len(s.orderIDs)-1],
				}
			},
			expectResp: &lomsservicev1.OrderInfoResponse{
				User:   defaultReq.User,
				Items:  defaultReq.Items,
				Status: lomsservicev1.Status_STATUS_AWAITING_PAYMENT,
			},
		},
		{
			name: "OrderInfoNotFoundItemFail",
			// Не используем сетапер что бы спровоцирвать ошибку
			setup:   func() {},
			cleanup: func(_ ...int64) {},
			req: func() *lomsservicev1.OrderInfoRequest {
				return &lomsservicev1.OrderInfoRequest{
					OrderID: 1111,
				}
			},
			expectErr: status.Error(codes.NotFound, "order not found"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()

			request := tt.req()
			orderResp, err := s.client.OrderInfo(context.Background(), request)

			s.ErrorIs(err, tt.expectErr)

			s.True(proto.Equal(tt.expectResp, orderResp))
			if tt.expectErr == nil {
				tt.cleanup(s.orderIDs...)
				s.orderIDs = []int64{}
			}
		})
	}
}

func (s *InmemSuite) TestOrderCreate() {
	tests := []struct {
		name       string
		req        *lomsservicev1.OrderCreateRequest
		cleanup    func(...int64)
		expectErr  error
		expectResp *lomsservicev1.OrderInfoResponse
	}{
		{
			name: "OrderCreateSuccess",
			req:  defaultReq,
			expectResp: &lomsservicev1.OrderInfoResponse{
				User:   defaultReq.User,
				Items:  defaultReq.Items,
				Status: lomsservicev1.Status_STATUS_AWAITING_PAYMENT,
			},
			cleanup: s.removeOrders,
		},
		{
			name: "OrderCreateNotFoundItemFail",
			req: &lomsservicev1.OrderCreateRequest{
				User: defaultUserID,
				Items: []*lomsservicev1.Item{
					{
						Sku:   111,
						Count: 10,
					},
				},
			},
			expectErr: status.Error(codes.NotFound, "SKU not found"),
		},
		{
			name: "OrderCreateInsufficientStockFail",
			req: &lomsservicev1.OrderCreateRequest{
				User: defaultUserID,
				Items: []*lomsservicev1.Item{
					{
						Sku:   defaultAlphaSKUID,
						Count: 50000,
					},
				},
			},
			expectErr: status.Error(codes.FailedPrecondition, "insufficient stock"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			createResp, err := s.client.OrderCreate(context.Background(), tt.req)
			s.ErrorIs(err, tt.expectErr)

			if tt.expectErr == nil {

				infoResp, err := s.client.OrderInfo(context.Background(), &lomsservicev1.OrderInfoRequest{
					OrderID: createResp.OrderID,
				})
				s.Require().NoError(err)

				s.True(proto.Equal(tt.expectResp, infoResp))

				tt.cleanup(createResp.OrderID)
			}
		})
	}
}

func (s *InmemSuite) TestOrderCancel() {
	tests := []struct {
		name       string
		req        func() *lomsservicev1.OrderCancelRequest
		setup      func()
		cleanup    func(...int64)
		expectErr  error
		expectResp *lomsservicev1.OrderInfoResponse
	}{
		{
			name:    "OrderCancelSuccess",
			setup:   s.setupOrders,
			cleanup: s.removeOrders,
			req: func() *lomsservicev1.OrderCancelRequest {
				s.Require().Len(s.orderIDs, 2)
				req := &lomsservicev1.OrderCancelRequest{
					OrderID: s.orderIDs[len(s.orderIDs)-1],
				}

				// Удаляем из списка, так как заказ уже отменен
				s.orderIDs = s.orderIDs[:len(s.orderIDs)-1]

				return req
			},
			expectResp: &lomsservicev1.OrderInfoResponse{
				User:   defaultReq.User,
				Items:  defaultReq.Items,
				Status: lomsservicev1.Status_STATUS_CANCELLED,
			},
		},
		{
			name:    "OrderCancelOrderNotFoundFail",
			setup:   func() {},
			cleanup: func(_ ...int64) {},
			req: func() *lomsservicev1.OrderCancelRequest {
				return &lomsservicev1.OrderCancelRequest{
					OrderID: 1111,
				}
			},
			expectErr: status.Error(codes.NotFound, "order not found"),
		},
		{
			name:    "OrderCancelAlreadyCancelledFail",
			setup:   s.setupOrders,
			cleanup: s.removeOrders,
			req: func() *lomsservicev1.OrderCancelRequest {
				s.Require().Len(s.orderIDs, 2)

				req := &lomsservicev1.OrderCancelRequest{
					OrderID: s.orderIDs[len(s.orderIDs)-1],
				}

				// Отменяем заказ, что бы спровоцировать ошибку
				_, err := s.client.OrderCancel(context.Background(), req)
				s.NoError(err)

				// Удаляем из списка, так как заказ уже отменен
				s.orderIDs = s.orderIDs[:len(s.orderIDs)-1]

				return req
			},
			expectErr: status.Error(codes.FailedPrecondition, "operation is not allowed for order with current status"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()

			request := tt.req()
			_, err := s.client.OrderCancel(context.Background(), request)
			s.ErrorIs(err, tt.expectErr)

			if tt.expectErr == nil {
				infoResp, err := s.client.OrderInfo(context.Background(), &lomsservicev1.OrderInfoRequest{
					OrderID: request.OrderID,
				})
				s.Require().NoError(err)

				s.True(proto.Equal(tt.expectResp, infoResp))
			}

			tt.cleanup(s.orderIDs...)
			// Обнуляем список созданных заказов в tt.setup()
			s.orderIDs = []int64{}
		})
	}
}

func (s *InmemSuite) TestStocksInfo() {
	tests := []struct {
		name       string
		req        *lomsservicev1.StocksInfoRequest
		expectErr  error
		expectResp *lomsservicev1.StocksInfoResponse
	}{
		{
			name: "OrderInfoSuccess",
			req: &lomsservicev1.StocksInfoRequest{
				Sku: defaultAlphaSKUID,
			},
			expectResp: &lomsservicev1.StocksInfoResponse{
				Count: 90,
			},
		},
		{
			name: "OrderInfoNotFoundItemFail",
			req: &lomsservicev1.StocksInfoRequest{
				Sku: 1,
			},
			expectErr: status.Error(codes.NotFound, "SKU not found"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orderResp, err := s.client.StocksInfo(context.Background(), tt.req)
			s.ErrorIs(err, tt.expectErr)

			s.True(proto.Equal(tt.expectResp, orderResp))
		})
	}
}

// TestOrderPay для этого теста берем SKU который не сможем откатить (paySKUID)
// поэтому нужно следить за порядком выполнения этого теста и отслеживать стоки
func (s *InmemSuite) TestOrderPay() {
	var commonOrderID int64

	items := []*lomsservicev1.Item{
		{
			Sku:   paySKUID,
			Count: 15,
		},
	}
	tests := []struct {
		name            string
		req             func() *lomsservicev1.OrderPayRequest
		expectErr       error
		expectResp      *lomsservicev1.OrderInfoResponse
		availableBefore uint64
		availableAfter  uint64
	}{
		{
			name: "OrderPaySuccess",
			req: func() *lomsservicev1.OrderPayRequest {
				infoResp, err := s.client.OrderCreate(context.Background(), &lomsservicev1.OrderCreateRequest{
					User:  defaultUserID,
					Items: items,
				})
				s.Require().NoError(err)

				commonOrderID = infoResp.OrderID

				return &lomsservicev1.OrderPayRequest{
					OrderID: infoResp.OrderID,
				}
			},
			expectResp: &lomsservicev1.OrderInfoResponse{
				User:   defaultUserID,
				Items:  items,
				Status: lomsservicev1.Status_STATUS_PAYED,
			},
			availableBefore: 45,
			availableAfter:  30,
		},
		{
			name: "OrderPayOrderNotFoundFail",
			req: func() *lomsservicev1.OrderPayRequest {
				return &lomsservicev1.OrderPayRequest{
					OrderID: 1111,
				}
			},
			availableBefore: 30,
			availableAfter:  30,
			expectErr:       status.Error(codes.NotFound, "order not found"),
		},
		{
			name: "OrderPayAlreadyPayedFail",
			req: func() *lomsservicev1.OrderPayRequest {
				return &lomsservicev1.OrderPayRequest{
					OrderID: commonOrderID,
				}
			},
			availableBefore: 30,
			availableAfter:  30,
			expectErr:       status.Error(codes.FailedPrecondition, "operation is not allowed for order with current status"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			available, err := s.client.StocksInfo(context.Background(), &lomsservicev1.StocksInfoRequest{Sku: paySKUID})
			s.NoError(err, tt.expectErr)
			s.Equal(tt.availableBefore, available.Count)

			request := tt.req()
			_, err = s.client.OrderPay(context.Background(), request)
			s.ErrorIs(err, tt.expectErr)

			if tt.expectErr == nil {
				infoResp, err := s.client.OrderInfo(context.Background(), &lomsservicev1.OrderInfoRequest{
					OrderID: request.OrderID,
				})
				s.Require().NoError(err)

				s.True(proto.Equal(tt.expectResp, infoResp))
			}

			available, err = s.client.StocksInfo(context.Background(), &lomsservicev1.StocksInfoRequest{Sku: paySKUID})
			s.NoError(err, tt.expectErr)
			s.Equal(tt.availableAfter, available.Count)

		})
	}
}
