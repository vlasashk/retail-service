package config

import (
	"errors"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type Config struct {
	Address         string `env:"CART_ADDR" env-default:"localhost:8082"`
	LoggerLVL       string `env:"CART_LOG_LVL" env-default:"debug"`
	ServiceName     string `env:"SERVICE_NAME" env-default:"cart"`
	ProductProvider ProductProviderCfg
	StocksProvider  StocksProviderCfg
	Telemetry       TelemetryCfg
}

type TelemetryCfg struct {
	ServiceName string `env:"SERVICE_NAME" env-default:"cart"`
	Address     string `env:"TELEMETRY_ADDR" env-default:"localhost:4317"`
}

type ProductProviderCfg struct {
	Address          string        `env:"PRODUCT_SERVICE_ADDR" env-default:"http://route256.pavl.uk:8080"`
	AccessToken      string        `env:"PRODUCT_SERVICE_TOKEN" env-default:"testtoken"`
	Retries          int           `env:"PRODUCT_SERVICE_RETRIES" env-default:"3"`
	MaxDelayForRetry int           `env:"PRODUCT_SERVICE_RETRY_DELAY" env-default:"3"`
	Timeout          time.Duration `env:"PRODUCT_SERVICE_TIMEOUT" env-default:"15s"`
	RateLimit        rate.Limit    `env:"PRODUCT_SERVICE_RATE_LIMIT" env-default:"10"`
	BurstLimit       int           `env:"PRODUCT_SERVICE_BURST_LIMIT" env-default:"10"`
}

type StocksProviderCfg struct {
	Address        string        `env:"LOMS_SERVICE_ADDR" env-default:"localhost:50000"`
	MaxConnTimeout time.Duration `env:"LOMS_SERVICE_MAX_TIMEOUT" env-default:"5s"`
	BaseDelay      time.Duration `env:"LOMS_SERVICE_BASE_DELAY" env-default:"1s"`
	Multiplier     float64       `env:"LOMS_SERVICE_MULTIPLIER" env-default:"1.6"`
	Jitter         float64       `env:"LOMS_SERVICE_JITTER" env-default:"0.2"`
	MaxDelay       time.Duration `env:"LOMS_SERVICE_MAX_DELAY" env-default:"5s"`
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
