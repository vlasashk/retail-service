package main

import (
	"context"

	"route256/loms/config"
	"route256/loms/internal/loms"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err = loms.Run(ctx, cfg); err != nil {
		log.Fatal().Err(err).Send()
	}
}
