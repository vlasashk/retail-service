package httpgate

import (
	"context"
	"net/http"

	"route256/loms/config"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

	return &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}, nil
}
