package ggrpc

import (
	"route256/loms/config"
	"route256/loms/internal/loms/ports/ggrpc/impl"
	"route256/loms/internal/loms/ports/ggrpc/interceptors"
	"route256/loms/internal/loms/ports/ggrpc/resources"
	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewServer(_ config.Config, res resources.Resources) *grpc.Server {
	options := buildOptions(res.Log)

	server := grpc.NewServer(options...)
	reflection.Register(server)

	serverAdapter := impl.New(res.Log, res.UseCase)

	lomsservicev1.RegisterLOMSServer(server, serverAdapter)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	return server
}

func buildOptions(log zerolog.Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor(log),
			interceptors.RecoverPanic,
			interceptors.Validate,
		),
	}
}
