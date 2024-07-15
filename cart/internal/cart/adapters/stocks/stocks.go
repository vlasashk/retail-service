package stocks

import (
	"context"
	"log"

	"route256/cart/config"
	"route256/cart/internal/cart/adapters/stocks/interceptors"
	"route256/cart/internal/cart/models"
	lomsservicev1 "route256/cart/pkg/api/loms/v1"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	baseURL string
	client  lomsservicev1.LOMSClient
	conn    *grpc.ClientConn
	log     zerolog.Logger
}

func New(cfg config.StocksProviderCfg, log zerolog.Logger) (*Client, error) {
	log.Debug().Str("host", cfg.Address).Msg("creating new stocks service client")
	backoffConfig := backoff.Config{
		BaseDelay:  cfg.BaseDelay,
		Multiplier: cfg.Multiplier,
		Jitter:     cfg.Jitter,
		MaxDelay:   cfg.MaxDelay,
	}

	conn, err := grpc.NewClient(cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoffConfig,
			MinConnectTimeout: cfg.MaxConnTimeout,
		}),
		grpc.WithChainUnaryInterceptor(
			interceptors.Metrics,
		),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	client := lomsservicev1.NewLOMSClient(conn)

	return &Client{
		baseURL: cfg.Address,
		log:     log,
		client:  client,
		conn:    conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) OrderCreate(ctx context.Context, order models.Order) (int64, error) {
	resp, err := c.client.OrderCreate(ctx, orderToDTO(order))
	if err != nil {
		return 0, handleError(err)
	}

	return resp.OrderID, nil
}

func (c *Client) StocksInfo(ctx context.Context, skuID int64) (uint64, error) {
	resp, err := c.client.StocksInfo(ctx, skuIDToDTO(skuID))
	if err != nil {
		return 0, handleError(err)
	}

	return resp.Count, nil
}

func handleError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		log.Fatalf("Non-gRPC error: %v", err)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.Wrap(models.ErrNotFound, err.Error())
	case codes.FailedPrecondition:
		return errors.Wrap(models.ErrInsufficientStock, err.Error())
	case codes.Internal:
		return errors.Wrap(models.ErrStockProvider, err.Error())
	default:
		return err
	}
}
