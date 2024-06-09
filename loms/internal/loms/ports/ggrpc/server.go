package ggrpc

import (
	"route256/loms/config"
	"route256/loms/internal/loms/ports/ggrpc/impl"
	"route256/loms/internal/loms/ports/ggrpc/interceptors"
	"route256/loms/internal/loms/ports/ggrpc/resources"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(cfg config.Config) (*grpc.Server, error) {
	res, err := resources.New(cfg)
	if err != nil {
		return nil, err
	}

	options, err := buildOptions(res.Log)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(options...)
	reflection.Register(server)

	serverAdapter := impl.New(res.Log, res.UseCase)

	lomsservicev1.RegisterLOMSServer(server, serverAdapter)

	return server, err
}

func buildOptions(log zerolog.Logger) ([]grpc.ServerOption, error) {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor(log),
			interceptors.RecoverPanic,
		),
	}, nil
}
