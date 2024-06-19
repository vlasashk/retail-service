//go:build integration

package suits

import (
	"context"
	"fmt"
	"os"

	"route256/loms/config"
	"route256/loms/internal/loms"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/migration"
	"route256/loms/pkg/pgconnect"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const migrationsPath = "../../migrations"

type PostgresSuit struct {
	suite.Suite
	client   lomsservicev1.LOMSClient
	conn     *grpc.ClientConn
	cancel   context.CancelFunc
	pool     *pgxpool.Pool
	orderIDs []int64
}

func (s *PostgresSuit) SetupSuite() {
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

	if err = healthCheck(s.conn, 10); err != nil {
		s.T().Fatal(err)
	}

	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.OrdersRepo.User,
		cfg.OrdersRepo.Password,
		cfg.OrdersRepo.Host,
		cfg.OrdersRepo.Port,
		cfg.OrdersRepo.Name)

	s.pool, err = pgconnect.Connect(ctx, url, zerolog.New(os.Stderr))
	if err != nil {
		s.T().Fatal(err)
	}

	if err = migration.Up(s.pool, migrationsPath); err != nil {
		s.T().Fatal(err)
	}
	s.TearDownTest()
}

func (s *PostgresSuit) TearDownSuite() {
	s.cancel()
	if s.conn != nil {
		_ = s.conn.Close()
	}

	if err := migration.Down(s.pool, migrationsPath); err != nil {
		s.T().Fatal(err)
	}
	s.pool.Close()
}

// setUpOrders вставляются 2 заказа на 3 SKU и база принимает следующий вид
// 1076963 - В stock-data.json 100 всего и 50 зарезервировано - 50 доступно
// 1148162 - В stock-data.json 150 всего и 50 зарезервировано - 100 доступно
// 1625903 - В stock-data.json 200 всего и 100 зарезервировано - 100 доступно
func (s *PostgresSuit) setUpOrders() {
	createOrderQry := `INSERT INTO orders.orders (user_id, status)
							VALUES ($1, $2)
							RETURNING id;`
	createItemQry := `INSERT INTO orders.order_items (sku_id, order_id, count)
							VALUES ($1, $2, $3);`

	var firstOrderID, secondOrderID int64
	ctx := context.Background()
	err := s.pool.QueryRow(ctx, createOrderQry, additionalReq.User, "AwaitingPayment").Scan(&firstOrderID)
	s.Require().NoError(err)
	for _, item := range additionalReq.Items {
		_, err = s.pool.Exec(ctx, createItemQry, item.Sku, firstOrderID, item.Count)
		s.Require().NoError(err)
	}

	err = s.pool.QueryRow(ctx, createOrderQry, defaultReq.User, "AwaitingPayment").Scan(&secondOrderID)
	s.Require().NoError(err)
	for _, item := range defaultReq.Items {
		_, err = s.pool.Exec(ctx, createItemQry, item.Sku, secondOrderID, item.Count)
		s.Require().NoError(err)
	}

	// updating stocks according to previous orders
	initStocks := `INSERT INTO stocks.stocks (id, available, reserved)
						VALUES (1076963, 100, 50),
							   (1148162, 150, 50),
							   (1625903, 200, 100),
							   (2618151, 50, 5)
						ON CONFLICT (id) DO UPDATE
						  SET available = EXCLUDED.available,
							  reserved  = EXCLUDED.reserved;`
	_, err = s.pool.Exec(context.Background(), initStocks)
	s.Require().NoError(err)

	s.orderIDs = append(s.orderIDs, firstOrderID, secondOrderID)
}

// SetupTest initializes default stock amount
func (s *PostgresSuit) SetupTest() {
	initStocks := `INSERT INTO stocks.stocks (id, available, reserved)
					VALUES (1076963, 100, 10),
					       (1148162, 150, 20),
					       (1625903, 200, 30),
					       (2618151, 50, 5)
						ON CONFLICT (id) DO UPDATE
						  SET available = EXCLUDED.available,
							  reserved  = EXCLUDED.reserved;`
	_, err := s.pool.Exec(context.Background(), initStocks)
	s.Require().NoError(err)
}

func (s *PostgresSuit) TearDownTest() {
	_, err := s.pool.Exec(context.Background(), "TRUNCATE TABLE orders.order_items;")
	s.Require().NoError(err)
	_, err = s.pool.Exec(context.Background(), "DELETE FROM orders.orders;")
	s.Require().NoError(err)

	// return default stocks amount to db
	s.SetupTest()
}

func (s *PostgresSuit) TestOrderInfo() {
	tests := []struct {
		name       string
		req        func() *lomsservicev1.OrderInfoRequest
		setup      func()
		cleanup    func()
		expectErr  error
		expectResp *lomsservicev1.OrderInfoResponse
	}{
		{
			name:    "OrderInfoSuccess",
			setup:   s.setUpOrders,
			cleanup: s.TearDownTest,
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
			// Не используем сетапер что бы спровоцировать ошибку
			setup:   func() {},
			cleanup: func() {},
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
				tt.cleanup()
				s.orderIDs = []int64{}
			}
		})
	}
}

func (s *PostgresSuit) TestOrderCreate() {
	tests := []struct {
		name       string
		req        *lomsservicev1.OrderCreateRequest
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

				s.TearDownTest()
			}
		})
	}
}

func (s *PostgresSuit) TestOrderCancel() {
	tests := []struct {
		name       string
		req        func() *lomsservicev1.OrderCancelRequest
		setup      func()
		cleanup    func()
		expectErr  error
		expectResp *lomsservicev1.OrderInfoResponse
	}{
		{
			name:    "OrderCancelSuccess",
			setup:   s.setUpOrders,
			cleanup: s.TearDownTest,
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
			cleanup: func() {},
			req: func() *lomsservicev1.OrderCancelRequest {
				return &lomsservicev1.OrderCancelRequest{
					OrderID: 1111,
				}
			},
			expectErr: status.Error(codes.NotFound, "order not found"),
		},
		{
			name:    "OrderCancelAlreadyCancelledFail",
			setup:   s.setUpOrders,
			cleanup: s.TearDownTest,
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

			tt.cleanup()
			// Обнуляем список созданных заказов в tt.setup()
			s.orderIDs = []int64{}
		})
	}
}

func (s *PostgresSuit) TestStocksInfo() {
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
func (s *PostgresSuit) TestOrderPay() {
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
