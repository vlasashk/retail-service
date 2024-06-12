package loms

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"route256/loms/config"
	"route256/loms/internal/loms/ports/ggrpc"
	"route256/loms/internal/loms/ports/ggrpc/resources"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		return err
	}
	res, err := resources.New(cfg)
	if err != nil {
		return err
	}
	grpcServer := ggrpc.NewServer(cfg, res)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Info().Msg(fmt.Sprintf("starting server: %s", cfg.Port))
		return grpcServer.Serve(grpcListener)
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Info().Msg("Got interruption signal")
		grpcServer.GracefulStop()
		return nil
	})

	if err = g.Wait(); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	log.Info().Msg("server was gracefully shut down")
	return nil
}
