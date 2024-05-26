package resources

import (
	"route256/cart/config"
	"route256/cart/internal/cart/adapters/inmem"
	"route256/cart/internal/cart/adapters/prodservice"
	"route256/cart/internal/cart/ports/vanilla/handlers/additem"
	"route256/cart/internal/cart/ports/vanilla/handlers/common"
	"route256/cart/internal/cart/ports/vanilla/handlers/getcart"
	"route256/cart/internal/cart/ports/vanilla/handlers/removeitem"
	"route256/cart/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log       zerolog.Logger
	Adder     additem.CartAdder
	Deleter   removeitem.CartDeleter
	Retriever getcart.CartRetriever
	Provider  common.ProductProvider
}

func NewResources(cfg config.CartConfig) (Resources, error) {
	resources := Resources{}

	log, err := logger.New(cfg.LoggerLVL)
	if err != nil {
		return Resources{}, err
	}
	resources.Log = log

	inMemStorage := inmem.NewStorage()

	resources.Adder = inMemStorage
	resources.Deleter = inMemStorage
	resources.Retriever = inMemStorage

	resources.Provider = prodservice.New(cfg.ProductProvider, log)

	return resources, nil
}
