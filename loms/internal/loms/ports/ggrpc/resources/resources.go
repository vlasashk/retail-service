package resources

import (
	"route256/loms/config"
	"route256/loms/internal/loms/adapters/ordersmem"
	"route256/loms/internal/loms/adapters/stocksmem"
	"route256/loms/internal/loms/usecase"

	"route256/cart/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log     zerolog.Logger
	UseCase *usecase.UseCase
}

func New(cfg config.Config) (Resources, error) {
	log, err := logger.New(cfg.LoggerLVL)
	if err != nil {
		return Resources{}, err
	}

	orderStorage := ordersmem.New()
	stockStorage := stocksmem.New()

	return Resources{
		Log: log,
		UseCase: usecase.New(
			log,
			orderStorage,
			stockStorage,
		),
	}, nil
}
