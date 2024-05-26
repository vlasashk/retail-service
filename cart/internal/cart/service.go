package cart

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"route256/cart/config"
	"route256/cart/internal/cart/ports/vanilla"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const gracefulTimeout = 10 * time.Second

func Run(ctx context.Context, cfg config.CartConfig) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv, err := vanilla.NewServer(cfg)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Info().Msg(fmt.Sprintf("starting server: %s", cfg.Address))
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Info().Msg("Got interruption signal")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()
		return srv.Shutdown(shutDownCtx)
	})

	if err = g.Wait(); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	log.Info().Msg("server was gracefully shut down")
	return nil
}
