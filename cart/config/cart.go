package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address         string `env:"CART_ADDR" env-default:"localhost:8082"`
	LoggerLVL       string `env:"CART_LOG_LVL" env-default:"debug"`
	ProductProvider ProductProviderCfg
	StocksProvider  StocksProviderCfg
}

type ProductProviderCfg struct {
	Address          string        `env:"PRODUCT_SERVICE_ADDR" env-default:"http://route256.pavl.uk:8080"`
	AccessToken      string        `env:"PRODUCT_SERVICE_TOKEN" env-default:"testtoken"`
	Retries          int           `env:"PRODUCT_SERVICE_RETRIES" env-default:"3"`
	MaxDelayForRetry int           `env:"PRODUCT_SERVICE_RETRY_DELAY" env-default:"3"`
	Timeout          time.Duration `env:"PRODUCT_SERVICE_TIMEOUT" env-default:"15s"`
}

type StocksProviderCfg struct {
	Address        string        `env:"LOMS_SERVICE_ADDR" env-default:"localhost:50000"`
	MaxConnTimeout time.Duration `env:"LOMS_SERVICE_MAX_TIMEOUT" env-default:"5s"`
	BaseDelay      time.Duration `env:"LOMS_SERVICE_BASE_DELAY" env-default:"1s"`
	Multiplier     float64       `env:"LOMS_SERVICE_MULTIPLIER" env-default:"1.6"`
	Jitter         float64       `env:"LOMS_SERVICE_JITTER" env-default:"0.2"`
	MaxDelay       time.Duration `env:"LOMS_SERVICE_MAX_DELAY" env-default:"5s"`
}

func New() (Config, error) {
	var res Config
	if err := cleanenv.ReadEnv(&res); err != nil {
		return Config{}, err
	}
	return res, nil
}
