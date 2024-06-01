package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Address         string `env:"CART_ADDR" env-default:"localhost:8082"`
	LoggerLVL       string `env:"CART_LOG_LVL" env-default:"info"`
	ProductProvider ProductProviderCfg
}

type ProductProviderCfg struct {
	Address          string `env:"PRODUCT_SERVICE_ADDR" env-default:"http://route256.pavl.uk:8080"`
	AccessToken      string `env:"PRODUCT_SERVICE_TOKEN" env-default:"testtoken"`
	Retries          int    `env:"PRODUCT_SERVICE_RETRIES" env-default:"3"`
	MaxDelayForRetry int    `env:"PRODUCT_SERVICE_RETRY_DELAY" env-default:"3"`
	Timeout          int    `env:"PRODUCT_SERVICE_TIMEOUT" env-default:"15"`
}

func New() (Config, error) {
	var res Config
	if err := cleanenv.ReadEnv(&res); err != nil {
		return Config{}, err
	}
	return res, nil
}
