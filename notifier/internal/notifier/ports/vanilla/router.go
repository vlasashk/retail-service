package vanilla

import (
	"net/http"

	"route256/notifier/config"
	"route256/notifier/internal/notifier/ports/vanilla/healthz"
	"route256/notifier/internal/notifier/ports/vanilla/muxer"
)

func NewServer(cfg config.Config) *http.Server {
	mux := muxer.NewMyMux()

	mux.HandleFunc("GET /healthz", healthz.HealthCheck)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: mux.Chain(),
	}
}
