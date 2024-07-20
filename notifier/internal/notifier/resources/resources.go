package resources

import (
	"context"

	"route256/notifier/config"
	"route256/notifier/internal/notifier/resources/consumer"
	"route256/notifier/pkg/logger"

	"github.com/rs/zerolog"
)

type Resources struct {
	Log           zerolog.Logger
	Consumer      *consumer.Consumer
	stopResources []func() error
}

func NewResources(_ context.Context, cfg config.Config) (Resources, error) {
	log, err := logger.New(cfg.LoggerLVL, cfg.ServiceName)
	if err != nil {
		return Resources{}, err
	}

	kfkConsumer, err := consumer.New(cfg.Consumer, log)
	if err != nil {
		return Resources{}, err
	}

	return Resources{
		Log:           log,
		Consumer:      kfkConsumer,
		stopResources: []func() error{},
	}, nil
}

func (r Resources) Stop() {
	for _, stop := range r.stopResources {
		if err := stop(); err != nil {
			r.Log.Error().Err(err).Send()
		}
	}
}
