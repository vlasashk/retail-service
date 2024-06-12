package resources

import (
	"route256/loms/config"
	"route256/loms/internal/loms/adapters/ordersmem"
	"route256/loms/internal/loms/adapters/stocksmem"
	"route256/loms/internal/loms/usecase"

	"route256/loms/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log           zerolog.Logger
	UseCase       *usecase.UseCase
	stopResources []func() error
}

func New(cfg config.Config) (Resources, error) {
	log, err := logger.New(cfg.LoggerLVL)
	if err != nil {
		return Resources{}, err
	}

	stockStorage, err := stocksmem.New()
	if err != nil {
		return Resources{}, err
	}

	orderStorage := ordersmem.New()

	return Resources{
		Log: log,
		UseCase: usecase.New(
			log,
			orderStorage,
			stockStorage,
		),
		stopResources: []func() error{},
	}, nil
}

func (r Resources) Stop() {
	for _, stop := range r.stopResources {
		if err := stop(); err != nil {
			r.Log.Error().Err(err).Send()
		}
	}
}
