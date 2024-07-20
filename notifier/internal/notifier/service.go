package notifier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"route256/notifier/config"
	"route256/notifier/internal/notifier/ports/vanilla"
	"route256/notifier/internal/notifier/resources"

	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog/log"
)

const gracefulTimeout = 10 * time.Second

func Run(ctx context.Context, cfg config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	res, err := resources.NewResources(ctx, cfg)
	if err != nil {
		return err
	}
	defer res.Stop()

	srv := vanilla.NewServer(cfg)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Info().Msg(fmt.Sprintf("starting server: %s", cfg.Address))
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		log.Info().Msg(fmt.Sprintf("starting consumer: %s", cfg.Consumer.GroupID))
		return res.Consumer.Run(gCtx)
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
