package resources

import (
	"context"

	"route256/loms/config"
	"route256/loms/internal/loms/adapters/pgorders"
	"route256/loms/internal/loms/adapters/pgstocks"
	"route256/loms/internal/loms/tracing"
	"route256/loms/internal/loms/usecase"

	"route256/loms/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log           zerolog.Logger
	UseCase       *usecase.UseCase
	stopResources []func() error
}

func New(ctx context.Context, cfg config.Config) (Resources, error) {
	log, err := logger.New(cfg.LoggerLVL, cfg.ServiceName)
	if err != nil {
		return Resources{}, err
	}

	stockStorage, err := pgstocks.New(ctx, cfg.StocksRepo)
	if err != nil {
		return Resources{}, err
	}

	orderStorage, err := pgorders.New(ctx, cfg.OrdersRepo)
	if err != nil {
		return Resources{}, err
	}

	trace, err := tracing.New(ctx, cfg.Telemetry)
	if err != nil {
		return Resources{}, err
	}

	return Resources{
		Log: log,
		UseCase: usecase.New(
			log,
			orderStorage,
			stockStorage,
		),
		stopResources: []func() error{
			orderStorage.Close,
			stockStorage.Close,
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
