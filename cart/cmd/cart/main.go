package main

import (
	"context"

	"route256/cart/config"
	"route256/cart/internal/cart"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err = cart.Run(ctx, cfg); err != nil {
		log.Fatal().Err(err).Send()
	}
}
