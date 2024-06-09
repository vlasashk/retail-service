package impl

import (
	"context"

	"route256/loms/internal/loms/models"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

type lomsServer interface {
	OrderCreate(ctx context.Context, order models.Order) (int64, error)
	OrderInfo(ctx context.Context, orderID int64) (models.Order, error)
	OrderPay(ctx context.Context, orderID int64) error
	OrderCancel(ctx context.Context, orderID int64) error
	StocksInfo(ctx context.Context, skuID uint32) (int64, error)
}

var _ lomsservicev1.LOMSServer = (*Impl)(nil)

type Impl struct {
	log    zerolog.Logger
	server lomsServer
	lomsservicev1.UnimplementedLOMSServer
}

func New(
	log zerolog.Logger,
	server lomsServer,
) *Impl {
	return &Impl{
		log:    log,
		server: server,
	}
}

func (i *Impl) OrderCreate(ctx context.Context, req *lomsservicev1.OrderCreateRequest) (*lomsservicev1.OrderCreateResponse, error) {
	orderID, err := i.server.OrderCreate(ctx, OrderToDomain(req))
	if err != nil {
		i.log.Error().Err(err).Send()
		return nil, mapError(err)
	}

	return &lomsservicev1.OrderCreateResponse{OrderID: orderID}, nil
}

func (i *Impl) OrderInfo(ctx context.Context, req *lomsservicev1.OrderInfoRequest) (*lomsservicev1.OrderInfoResponse, error) {
	order, err := i.server.OrderInfo(ctx, req.OrderID)
	if err != nil {
		i.log.Error().Err(err).Send()
		return nil, mapError(err)
	}

	return OrderToProto(order), nil
}
func (i *Impl) OrderPay(ctx context.Context, req *lomsservicev1.OrderPayRequest) (*emptypb.Empty, error) {
	err := i.server.OrderPay(ctx, req.OrderID)
	if err != nil {
		i.log.Error().Err(err).Send()
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
func (i *Impl) OrderCancel(ctx context.Context, req *lomsservicev1.OrderCancelRequest) (*emptypb.Empty, error) {
	err := i.server.OrderCancel(ctx, req.OrderID)
	if err != nil {
		i.log.Error().Err(err).Send()
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
func (i *Impl) StocksInfo(ctx context.Context, req *lomsservicev1.StocksInfoRequest) (*lomsservicev1.StocksInfoResponse, error) {
	count, err := i.server.StocksInfo(ctx, req.Sku)
	if err != nil {
		i.log.Error().Err(err).Send()
		return nil, mapError(err)
	}

	return &lomsservicev1.StocksInfoResponse{Count: uint64(count)}, nil
}
