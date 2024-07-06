package httpgate

import (
	"context"
	"net/http"
	"net/http/pprof"

	"route256/loms/config"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(ctx context.Context, cfg config.HTTPGateCfg) (*http.Server, error) {
	rmux := runtime.NewServeMux(runtime.WithErrorHandler(errHandler))

	conn, err := grpc.NewClient(cfg.LOMSAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := lomsservicev1.NewLOMSClient(conn)

	if err = lomsservicev1.RegisterLOMSHandlerClient(ctx, rmux, client); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.SwaggerFilePath)
	})

	// mount the Swagger UI that uses the OpenAPI specification path above
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir(cfg.SwaggerDirPath))))

	registerPprofHandlers(mux)

	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}, nil
}

func registerPprofHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register other pprof handlers
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
}
