package resources

import (
	"context"

	"route256/cart/config"
	"route256/cart/internal/cart/adapters/inmem"
	"route256/cart/internal/cart/adapters/prodservice"
	"route256/cart/internal/cart/adapters/stocks"
	"route256/cart/internal/cart/tracing"
	"route256/cart/internal/cart/usecase"
	"route256/cart/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log           zerolog.Logger
	UseCase       *usecase.UseCase
	stopResources []func() error
}

func NewResources(ctx context.Context, cfg config.Config) (Resources, error) {
	log, err := logger.New(cfg.LoggerLVL, cfg.ServiceName)
	if err != nil {
		return Resources{}, err
	}

	stocksProvider, err := stocks.New(cfg.StocksProvider, log)
	if err != nil {
		return Resources{}, err
	}

	inMemStorage := inmem.NewStorage()
	productProvider := prodservice.New(cfg.ProductProvider, log)

	trace, err := tracing.New(ctx, cfg.Telemetry)
	if err != nil {
		return Resources{}, err
	}

	return Resources{
		Log: log,
		UseCase: usecase.New(
			inMemStorage,
			inMemStorage,
			inMemStorage,
			inMemStorage,
			productProvider,
			stocksProvider,
		),
		stopResources: []func() error{
			stocksProvider.Close,
			trace.Close,
		},
	}, nil
}

func (r Resources) Stop() {
	for _, stop := range r.stopResources {
		if err := stop(); err != nil {
			r.Log.Error().Err(err).Send()
		}
	}
}
