package resources

import (
	"route256/cart/config"
	"route256/cart/internal/cart/adapters/inmem"
	"route256/cart/internal/cart/adapters/prodservice"
	"route256/cart/internal/cart/usecase"
	"route256/cart/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log     zerolog.Logger
	UseCase *usecase.UseCase
}

func NewResources(cfg config.Config) (Resources, error) {
	log, err := logger.New(cfg.LoggerLVL)
	if err != nil {
		return Resources{}, err
	}

	inMemStorage := inmem.NewStorage()
	productProvider := prodservice.New(cfg.ProductProvider, log)

	return Resources{
		Log: log,
		UseCase: usecase.New(
			inMemStorage,
			inMemStorage,
			inMemStorage,
			inMemStorage,
			productProvider,
		),
	}, nil
}
