package loms

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"route256/loms/config"
	"route256/loms/internal/loms/ports/ggrpc"
	"route256/loms/internal/loms/ports/ggrpc/resources"
	"route256/loms/internal/loms/ports/httpgate"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const gracefulTimeout = 10 * time.Second

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
	defer res.Stop()

	g, gCtx := errgroup.WithContext(ctx)

	grpcServer := ggrpc.NewServer(cfg, res)
	gatewayServer, err := httpgate.New(gCtx, cfg.Gateway)
	if err != nil {
		return err
	}

	g.Go(func() error {
		log.Info().Msg(fmt.Sprintf("starting server: %s", cfg.Port))
		return grpcServer.Serve(grpcListener)
	})
	g.Go(func() error {
		log.Info().Msg(fmt.Sprintf("starting gateway: %s", cfg.Gateway.Address))
		if err = gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Info().Msg("Got interruption signal")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()
		grpcServer.GracefulStop()
		return gatewayServer.Shutdown(shutDownCtx)
	})

	if err = g.Wait(); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	log.Info().Msg("server was gracefully shut down")
	return nil
}
