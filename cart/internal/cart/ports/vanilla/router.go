package vanilla

import (
	"net/http"
	"net/http/pprof"

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

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(cfg config.Config, res resources.Resources) *http.Server {
	mux := muxer.NewMyMux()

	addRoutes(mux, res)
	registerPprofHandlers(mux)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: mux.Chain(),
	}
}

func addRoutes(mux *muxer.MyMux, resources resources.Resources) {
	mux.Use(middleware.TraceMiddleware)
	mux.Use(middleware.LoggingMiddleware(resources.Log))
	mux.Use(middleware.Metrics)
	mux.Use(middleware.Recover)

	mux.Handle("POST /user/{user_id}/cart/{sku_id}", additem.New(resources.UseCase))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", removeitem.New(resources.UseCase))
	mux.Handle("DELETE /user/{user_id}/cart", removecart.New(resources.UseCase))
	mux.Handle("GET /user/{user_id}/cart/list", getcart.New(resources.UseCase))
	mux.Handle("GET /user/{user_id}/cart/checkout", checkout.New(resources.UseCase))

	mux.Handle("GET /metrics", promhttp.Handler())

	mux.HandleFunc("GET /healthz", healthz.HealthCheck)
}

func registerPprofHandlers(mux *muxer.MyMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
}
