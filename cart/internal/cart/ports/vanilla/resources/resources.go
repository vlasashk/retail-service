package resources

import (
	"route256/cart/config"
	"route256/cart/internal/cart/usecase"
	"route256/cart/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log     zerolog.Logger
	UseCase *usecase.UseCase
}

func NewResources(cfg config.CartConfig) (Resources, error) {
	log, err := logger.New(cfg.LoggerLVL)
	if err != nil {
		return Resources{}, err
	}

	return Resources{
		Log:     log,
		UseCase: usecase.New(cfg.ProductProvider, log),
	}, nil
}
