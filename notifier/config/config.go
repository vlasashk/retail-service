package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address     string `env:"NOTIFIER_ADDR" env-default:"localhost:8099"`
	LoggerLVL   string `env:"NOTIFIER_LOG_LVL" env-default:"debug"`
	ServiceName string `env:"SERVICE_NAME" env-default:"notifier"`
	Consumer    KafkaCfg
}

type KafkaCfg struct {
	Brokers []string `env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	Topic   []string `env:"KAFKA_TOPICS" env-separator:"," env-default:"loms.order-events"`
	GroupID string   `env:"KAFKA_GROUP_ID" env-default:"order_events"`
}

// New - parses environment variables
func New(opts ...Option) (Config, error) {
	o := defaultOptions()

	for _, opt := range opts {
		opt.apply(&o)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(o.configPath, &cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
		log.Warn().Err(err).Msg("yaml config file not found")

		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
