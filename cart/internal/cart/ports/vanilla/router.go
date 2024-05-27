package vanilla

import (
	"net/http"

	"route256/cart/config"
	"route256/cart/internal/cart/ports/vanilla/handlers/additem"
	"route256/cart/internal/cart/ports/vanilla/handlers/getcart"
	"route256/cart/internal/cart/ports/vanilla/handlers/removecart"
	"route256/cart/internal/cart/ports/vanilla/handlers/removeitem"
	"route256/cart/internal/cart/ports/vanilla/middleware"
	"route256/cart/internal/cart/ports/vanilla/muxer"
	"route256/cart/internal/cart/ports/vanilla/resources"
)

func NewServer(cfg config.CartConfig) (*http.Server, error) {
	mux := muxer.NewMyMux()

	res, err := resources.NewResources(cfg)
	if err != nil {
		return nil, err
	}

	addRoutes(mux, res)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: mux.Chain(),
	}, nil
}

func addRoutes(mux *muxer.MyMux, resources resources.Resources) {
	mux.Use(middleware.LoggingMiddleware(resources.Log))

	mux.Handle("POST /user/{user_id}/cart/{sku_id}", additem.New(resources.Log, resources.Adder, resources.Provider))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", removeitem.New(resources.Log, resources.ItemRemover))
	mux.Handle("DELETE /user/{user_id}/cart", removecart.New(resources.Log, resources.CartRemover))
	mux.Handle("GET /user/{user_id}/cart", getcart.New(resources.Log, resources.Retriever, resources.Provider))
}
