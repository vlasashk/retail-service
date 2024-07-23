package main

import (
	"context"

	"route256/notifier/config"
	"route256/notifier/internal/notifier"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err = notifier.Run(ctx, cfg); err != nil {
		log.Fatal().Err(err).Send()
	}
}
