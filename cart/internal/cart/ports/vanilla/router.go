package vanilla

import (
	"net/http"

	"route256/cart/config"
	"route256/cart/internal/cart/ports/vanilla/handlers/additem"
	"route256/cart/internal/cart/ports/vanilla/handlers/checkout"
	"route256/cart/internal/cart/ports/vanilla/handlers/getcart"
	"route256/cart/internal/cart/ports/vanilla/handlers/healthz"
	"route256/cart/internal/cart/ports/vanilla/handlers/removecart"
	"route256/cart/internal/cart/ports/vanilla/handlers/removeitem"
	"route256/cart/internal/cart/ports/vanilla/middleware"
	"route256/cart/internal/cart/ports/vanilla/muxer"
	"route256/cart/internal/cart/resources"
)

func NewServer(cfg config.Config, res resources.Resources) *http.Server {
	mux := muxer.NewMyMux()

	addRoutes(mux, res)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: mux.Chain(),
	}
}

func addRoutes(mux *muxer.MyMux, resources resources.Resources) {
	mux.Use(middleware.LoggingMiddleware(resources.Log))
	mux.Use(middleware.Recover)

	mux.Handle("POST /user/{user_id}/cart/{sku_id}", additem.New(resources.Log, resources.UseCase))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", removeitem.New(resources.Log, resources.UseCase))
	mux.Handle("DELETE /user/{user_id}/cart", removecart.New(resources.Log, resources.UseCase))
	mux.Handle("GET /user/{user_id}/cart/list", getcart.New(resources.Log, resources.UseCase))
	mux.Handle("GET /user/{user_id}/cart/checkout", checkout.New(resources.Log, resources.UseCase))
	mux.HandleFunc("GET /healthz", healthz.HealthCheck)
}
